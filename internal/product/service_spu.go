package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type SpuService struct {
	repo         IspuRepository
	categoryRepo IcategoryRepository
	brandRepo    IbrandRepository
	attrRepo     IattributeRepository
	db           *gorm.DB
}

func NewSpuService(
	repo IspuRepository,
	categoryRepo IcategoryRepository,
	brandRepo IbrandRepository,
	attrRepo IattributeRepository,
	db *gorm.DB,
) *SpuService {
	return &SpuService{
		repo:         repo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
		attrRepo:     attrRepo,
		db:           db,
	}
}

type SPUListResult struct {
	Total int64  `json:"total"`
	List  []*SPU `json:"list"`
}

// Create 创建商品（事务内写 SPU + SKU + Description + ProductAttribute）
func (s *SpuService) Create(ctx context.Context, req *CreateSPUReq) (*SPU, error) {
	// 校验类目
	if _, err := s.categoryRepo.FindByID(ctx, req.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrCategoryNotFound
		}
		return nil, err
	}
	// 校验品牌
	if req.BrandID > 0 {
		if _, err := s.brandRepo.FindByID(ctx, req.BrandID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errcode.ErrBrandNotFound
			}
			return nil, err
		}
	}
	// 校验 SKU 编码重复（请求内）
	seen := map[string]bool{}
	for _, sku := range req.SKUs {
		if seen[sku.SkuCode] {
			return nil, fmt.Errorf("%w: %s", errcode.ErrSKUCodeExists, sku.SkuCode)
		}
		seen[sku.SkuCode] = true
	}
	// 校验 SKU 编码重复（数据库）
	for _, sku := range req.SKUs {
		existing, err := s.repo.FindSKUByCode(ctx, sku.SkuCode)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, fmt.Errorf("%w: %s", errcode.ErrSKUCodeExists, sku.SkuCode)
		}
	}
	// 校验属性
	if len(req.Attributes) > 0 {
		attrIDs := map[int64]bool{}
		for _, pa := range req.Attributes {
			if attrIDs[pa.AttributeID] {
				return nil, errcode.ErrProductAttrDuplicate
			}
			attrIDs[pa.AttributeID] = true
			if _, err := s.attrRepo.FindByID(ctx, pa.AttributeID); err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errcode.ErrAttributeNotFound
				}
				return nil, err
			}
		}
	}

	// 序列化 images
	var imagesJSON string
	if len(req.Images) > 0 {
		b, _ := json.Marshal(req.Images)
		imagesJSON = string(b)
	}

	// 事务创建
	var spu *SPU
	err := s.db.Transaction(func(tx *gorm.DB) error {
		spu = &SPU{
			Name:        req.Name,
			Subtitle:    req.Subtitle,
			CategoryID:  req.CategoryID,
			BrandID:     req.BrandID,
			Unit:        req.Unit,
			MainImage:   req.MainImage,
			Images:      imagesJSON,
			VideoURL:    req.VideoURL,
			SortOrder:   req.SortOrder,
			CreatedBy:   req.CreatedBy,
			Status:      0, // draft
		}
		if req.Description != "" || req.MobileDesc != "" {
			spu.HasDescription = 1
		}
		if err := s.repo.CreateSPUWithTx(tx, spu); err != nil {
			return err
		}

		// 批量创建 SKU
		for _, sku := range req.SKUs {
			specJSON := "{}"
			if sku.Spec != nil {
				b, _ := json.Marshal(sku.Spec)
				specJSON = string(b)
			}
			skuModel := &SKU{
				ProductID:     spu.ID,
				SkuCode:       sku.SkuCode,
				Barcode:       sku.Barcode,
				Spec:          specJSON,
				Price:         sku.Price,
				MarketPrice:   sku.MarketPrice,
				CostPrice:     sku.CostPrice,
				Weight:        sku.Weight,
				Volume:        sku.Volume,
				Length:        sku.Length,
				Width:         sku.Width,
				Height:        sku.Height,
				MinPurchaseQty: sku.MinPurchase,
				MaxPurchaseQty: sku.MaxPurchase,
				Image:         sku.Image,
				Status:        1,
			}
			if skuModel.MinPurchaseQty <= 0 {
				skuModel.MinPurchaseQty = 1
			}
			if err := s.repo.CreateSKUWithTx(tx, skuModel); err != nil {
				return err
			}
		}

		// 创建图文详情
		if req.Description != "" || req.MobileDesc != "" {
			desc := &Description{
				ProductID:         spu.ID,
				Description:       req.Description,
				MobileDescription: req.MobileDesc,
			}
			if err := s.repo.CreateDescriptionWithTx(tx, desc); err != nil {
				return err
			}
		}

		// 创建属性值
		for _, pa := range req.Attributes {
			paModel := &ProductAttribute{
				ProductID:   spu.ID,
				AttributeID: pa.AttributeID,
				Value:       pa.Value,
			}
			if err := s.repo.CreateProductAttrWithTx(tx, paModel); err != nil {
				return err
			}
		}

		// 更新聚合价格
		if err := s.repo.UpdateSPUPriceWithTx(tx, spu.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return spu, nil
}

// GetByID 获取商品详情（含 SKU）
func (s *SpuService) GetByID(ctx context.Context, id int64) (*SPU, error) {
	spu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrSPUNotFound
		}
		return nil, err
	}
	return spu, nil
}

// List 商品列表
func (s *SpuService) List(ctx context.Context, req *SPUListReq) (*SPUListResult, error) {
	req.Normalize()
	list, total, err := s.repo.List(ctx, req.Name, req.CategoryID, req.BrandID, req.Status, req.PriceMin, req.PriceMax, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*SPU, len(list))
	for i := range list {
		items[i] = &list[i]
	}
	return &SPUListResult{Total: total, List: items}, nil
}

// Update 更新商品基本信息
func (s *SpuService) Update(ctx context.Context, id int64, req *UpdateSPUReq) (*SPU, error) {
	spu, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrSPUNotFound
		}
		return nil, err
	}
	if req.Name != nil {
		spu.Name = *req.Name
	}
	if req.Subtitle != nil {
		spu.Subtitle = *req.Subtitle
	}
	if req.Unit != nil {
		spu.Unit = *req.Unit
	}
	if req.MainImage != nil {
		spu.MainImage = *req.MainImage
	}
	if req.Images != nil {
		b, _ := json.Marshal(*req.Images)
		spu.Images = string(b)
	}
	if req.VideoURL != nil {
		spu.VideoURL = *req.VideoURL
	}
	if req.SortOrder != nil {
		spu.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		spu.Status = *req.Status
	}
	if req.UpdatedBy != nil {
		spu.UpdatedBy = *req.UpdatedBy
	}
	if err := s.repo.Update(ctx, spu); err != nil {
		return nil, err
	}
	return spu, nil
}

// Delete 删除商品（软删除）
func (s *SpuService) Delete(ctx context.Context, id int64) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrSPUNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
