package product

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"eshop-monolith/pkg/errcode"
)

type SKUService struct {
	repo IspuRepository
}

func NewSKUService(repo IspuRepository) *SKUService {
	return &SKUService{repo: repo}
}

func (s *SKUService) GetByID(ctx context.Context, id int64) (*SKU, error) {
	sku, err := s.repo.FindSKUByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return sku, nil
}

func (s *SKUService) GetByCode(ctx context.Context, code string) (*SKU, error) {
	sku, err := s.repo.FindSKUByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}
	return sku, nil
}

func (s *SKUService) Update(ctx context.Context, id int64, req *UpdateSKUReq) (*SKU, error) {
	sku, err := s.repo.FindSKUByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrNotFound
		}
		return nil, err
	}

	if req.Price != nil {
		sku.Price = *req.Price
	}
	if req.MarketPrice != nil {
		sku.MarketPrice = *req.MarketPrice
	}
	if req.CostPrice != nil {
		sku.CostPrice = *req.CostPrice
	}
	if req.Status != nil {
		sku.Status = *req.Status
	}
	if req.Image != nil {
		sku.Image = *req.Image
	}
	if req.Barcode != nil {
		sku.Barcode = *req.Barcode
	}
	if req.Weight != nil {
		sku.Weight = *req.Weight
	}
	if req.Volume != nil {
		sku.Volume = *req.Volume
	}
	if req.Length != nil {
		sku.Length = *req.Length
	}
	if req.Width != nil {
		sku.Width = *req.Width
	}
	if req.Height != nil {
		sku.Height = *req.Height
	}

	if err := s.repo.UpdateSKU(ctx, sku); err != nil {
		return nil, err
	}
	return sku, nil
}

func (s *SKUService) Delete(ctx context.Context, id int64) error {
	_, err := s.repo.FindSKUByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errcode.ErrNotFound
		}
		return err
	}
	return s.repo.DeleteSKU(ctx, id)
}

func (s *SKUService) ListByProduct(ctx context.Context, productID int64) ([]SKU, error) {
	return s.repo.FindSKUsByProductID(ctx, productID)
}
