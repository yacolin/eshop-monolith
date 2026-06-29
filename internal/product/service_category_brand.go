package product

import (
	"context"
	"errors"

	"eshop-monolith/pkg/errcode"

	"gorm.io/gorm"
)

type CategoryBrandService struct {
	repo         IcategoryBrandRepository
	categoryRepo IcategoryRepository
	brandRepo    IbrandRepository
	db           *gorm.DB
}

func NewCategoryBrandService(
	repo IcategoryBrandRepository,
	categoryRepo IcategoryRepository,
	brandRepo IbrandRepository,
	db *gorm.DB,
) *CategoryBrandService {
	return &CategoryBrandService{repo: repo, categoryRepo: categoryRepo, brandRepo: brandRepo, db: db}
}

// SetCategoryBrands 全量替换类目下的品牌关联
func (s *CategoryBrandService) SetCategoryBrands(ctx context.Context, categoryID int64, req *SetCategoryBrandsReq) error {
	// 校验类目存在
	if _, err := s.categoryRepo.FindByID(ctx, categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrCategoryNotFound
		}
		return err
	}
	// 校验品牌存在
	seen := map[int64]bool{}
	items := make([]CategoryBrand, 0, len(req.BrandIDs))
	for _, brandID := range req.BrandIDs {
		if seen[brandID] {
			continue
		}
		seen[brandID] = true
		if _, err := s.brandRepo.FindByID(ctx, brandID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errcode.ErrBrandNotFound
			}
			return err
		}
		items = append(items, CategoryBrand{
			CategoryID: categoryID,
			BrandID:    brandID,
			SortOrder:  req.SortOrder,
		})
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.DeleteByCategory(tx, categoryID); err != nil {
			return err
		}
		return s.repo.AddBatch(tx, items)
	})
}

// ListCategoryBrands 查类目下的品牌列表
func (s *CategoryBrandService) ListCategoryBrands(ctx context.Context, categoryID int64) ([]CategoryBrand, error) {
	if _, err := s.categoryRepo.FindByID(ctx, categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrCategoryNotFound
		}
		return nil, err
	}
	return s.repo.FindBrandsByCategory(ctx, categoryID)
}
