package dto

import (
	"eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/query"
	"eshop-monolith/pkg/utils"
)

// ── Attribute ──────────────────────────────────────────────────

type CreateAttributeDTO struct {
	Name      string `json:"name" binding:"required,max=100"`
	SortOrder int    `json:"sort_order"`
}

type UpdateAttributeDTO struct {
	Name      *string `json:"name" binding:"omitempty,max=100"`
	SortOrder *int    `json:"sort_order"`
}

type AttributeResponse struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	SortOrder int             `json:"sort_order"`
	CreatedAt utils.Timestamp `json:"created_at"`
	UpdatedAt utils.Timestamp `json:"updated_at"`
}

type AttributeListResult struct {
	List  []AttributeResponse `json:"list"`
	Total int64               `json:"total"`
}

type AttributeListQuery struct {
	query.Pagination
}

// ── AttributeValue ─────────────────────────────────────────────

type CreateAttributeValueDTO struct {
	AttributeID int64  `json:"attribute_id" binding:"required"`
	Value       string `json:"value" binding:"required,max=100"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateAttributeValueDTO struct {
	Value     *string `json:"value" binding:"omitempty,max=100"`
	SortOrder *int    `json:"sort_order"`
}

type AttributeValueResponse struct {
	ID          int64           `json:"id"`
	AttributeID int64           `json:"attribute_id"`
	Value       string          `json:"value"`
	SortOrder   int             `json:"sort_order"`
	CreatedAt   utils.Timestamp `json:"created_at"`
	UpdatedAt   utils.Timestamp `json:"updated_at"`
}

// ── Converters ─────────────────────────────────────────────────

func AttributeToResponse(a *models.Attribute) AttributeResponse {
	return AttributeResponse{
		ID: a.ID, Name: a.Name, SortOrder: a.SortOrder,
		CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
	}
}

func AttributeValueToResponse(v *models.AttributeValue) AttributeValueResponse {
	return AttributeValueResponse{
		ID: v.ID, AttributeID: v.AttributeID, Value: v.Value, SortOrder: v.SortOrder,
		CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt,
	}
}
