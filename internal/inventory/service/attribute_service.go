package service

import (
	"context"

	"eshop-monolith/internal/inventory/api/dto"
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/internal/inventory/domain/repositories"
)

type AttributeService struct {
	attrRepo repositories.IattributeRepository
}

func NewAttributeService(attrRepo repositories.IattributeRepository) *AttributeService {
	return &AttributeService{attrRepo: attrRepo}
}

// ── Attribute ──────────────────────────────────────────────────

func (s *AttributeService) CreateAttribute(ctx context.Context, req *dto.CreateAttributeDTO) (*models.Attribute, error) {
	attr := &models.Attribute{Name: req.Name, SortOrder: req.SortOrder}
	if err := s.attrRepo.Create(ctx, attr); err != nil {
		return nil, err
	}
	return attr, nil
}

func (s *AttributeService) GetAttribute(ctx context.Context, id int64) (*models.Attribute, error) {
	return s.attrRepo.FindByID(ctx, id)
}

func (s *AttributeService) ListAttributes(ctx context.Context, q dto.AttributeListQuery) (*dto.AttributeListResult, error) {
	offset := (q.Page - 1) * q.Size
	list, err := s.attrRepo.FindAll(ctx, offset, q.Size)
	if err != nil {
		return nil, err
	}
	total, err := s.attrRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]dto.AttributeResponse, len(list))
	for i, a := range list {
		items[i] = dto.AttributeToResponse(&a)
	}
	return &dto.AttributeListResult{List: items, Total: total}, nil
}

func (s *AttributeService) UpdateAttribute(ctx context.Context, id int64, req *dto.UpdateAttributeDTO) (*models.Attribute, error) {
	attr, err := s.attrRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		attr.Name = *req.Name
	}
	if req.SortOrder != nil {
		attr.SortOrder = *req.SortOrder
	}
	if err := s.attrRepo.Update(ctx, attr); err != nil {
		return nil, err
	}
	return attr, nil
}

func (s *AttributeService) DeleteAttribute(ctx context.Context, id int64) error {
	return s.attrRepo.Delete(ctx, id)
}

// ── AttributeValue ─────────────────────────────────────────────

func (s *AttributeService) CreateAttributeValue(ctx context.Context, req *dto.CreateAttributeValueDTO) (*models.AttributeValue, error) {
	val := &models.AttributeValue{
		AttributeID: req.AttributeID,
		Value:       req.Value,
		SortOrder:   req.SortOrder,
	}
	if err := s.attrRepo.CreateValue(ctx, val); err != nil {
		return nil, err
	}
	return val, nil
}

func (s *AttributeService) GetAttributeValue(ctx context.Context, id int64) (*models.AttributeValue, error) {
	return s.attrRepo.FindValueByID(ctx, id)
}

func (s *AttributeService) ListAttributeValues(ctx context.Context, attributeID int64) ([]dto.AttributeValueResponse, error) {
	vals, err := s.attrRepo.FindValuesByAttributeID(ctx, attributeID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.AttributeValueResponse, len(vals))
	for i, v := range vals {
		items[i] = dto.AttributeValueToResponse(&v)
	}
	return items, nil
}

func (s *AttributeService) UpdateAttributeValue(ctx context.Context, id int64, req *dto.UpdateAttributeValueDTO) (*models.AttributeValue, error) {
	val, err := s.attrRepo.FindValueByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Value != nil {
		val.Value = *req.Value
	}
	if req.SortOrder != nil {
		val.SortOrder = *req.SortOrder
	}
	if err := s.attrRepo.UpdateValue(ctx, val); err != nil {
		return nil, err
	}
	return val, nil
}

func (s *AttributeService) DeleteAttributeValue(ctx context.Context, id int64) error {
	return s.attrRepo.DeleteValue(ctx, id)
}
