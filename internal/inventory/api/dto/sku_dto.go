package dto

import (
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/utils"
)

type CreateSkuDTO struct {
	ProductID int64             `json:"product_id" binding:"required"`
	Name      string            `json:"name" binding:"required,max=255"`
	Price     int64             `json:"price" binding:"required,gt=0"`
	SKUCode   string            `json:"sku_code" binding:"required,max=100"`
	Image     string            `json:"image" binding:"max=500"`
	Spec      map[string]string `json:"spec"`
}

type UpdateSkuDTO struct {
	Name     *string            `json:"name" binding:"omitempty,max=255"`
	Price    *int64             `json:"price" binding:"omitempty,gt=0"`
	SKUCode  *string            `json:"sku_code" binding:"omitempty,max=100"`
	Image    *string            `json:"image" binding:"omitempty,max=500"`
	Spec     map[string]string  `json:"spec"`
}

type SkuResponse struct {
	ID        int64             `json:"id"`
	ProductID int64             `json:"product_id"`
	Name      string            `json:"name"`
	Price     int64             `json:"price"`
	SKUCode   string            `json:"sku_code"`
	Image     string            `json:"image,omitempty"`
	Spec      map[string]string `json:"spec,omitempty"`
	CreatedAt utils.Timestamp   `json:"created_at"`
	UpdatedAt utils.Timestamp   `json:"updated_at"`
}

type SkuListResult struct {
	List  []SkuResponse `json:"list"`
	Total int64         `json:"total"`
}

func SkuToResponse(s *models.Sku) SkuResponse {
	return SkuResponse{
		ID: s.ID, ProductID: s.ProductID, Name: s.Name, Price: s.Price,
		SKUCode: s.SKUCode, Image: s.Image, Spec: s.Spec,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}
