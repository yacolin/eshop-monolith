package models

import (
	"time"

	domain "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/utils"
	"gorm.io/gorm"
)

// ── AttributePO ──────────────────────────────────────────────────────

type AttributePO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"type:varchar(100);not null;uniqueIndex"`
	SortOrder int            `gorm:"default:0"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (AttributePO) TableName() string { return "attribute_attributes" }

func (po *AttributePO) ToDomain() *domain.Attribute {
	return &domain.Attribute{
		ID: po.ID, Name: po.Name, SortOrder: po.SortOrder,
		CreatedAt: utils.Timestamp(po.CreatedAt), UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
}

func AttributeFromDomain(a *domain.Attribute) *AttributePO {
	return &AttributePO{
		ID: a.ID, Name: a.Name, SortOrder: a.SortOrder,
		CreatedAt: time.Time(a.CreatedAt), UpdatedAt: time.Time(a.UpdatedAt),
	}
}

// ── AttributeValuePO ────────────────────────────────────────────────

type AttributeValuePO struct {
	ID          int64          `gorm:"primaryKey;autoIncrement"`
	AttributeID int64          `gorm:"index;not null"`
	Value       string         `gorm:"type:varchar(100);not null"`
	SortOrder   int            `gorm:"default:0"`
	CreatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt   time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (AttributeValuePO) TableName() string { return "attribute_values" }

func (po *AttributeValuePO) ToDomain() *domain.AttributeValue {
	return &domain.AttributeValue{
		ID: po.ID, AttributeID: po.AttributeID, Value: po.Value, SortOrder: po.SortOrder,
		CreatedAt: utils.Timestamp(po.CreatedAt), UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
}

func AttributeValueFromDomain(v *domain.AttributeValue) *AttributeValuePO {
	return &AttributeValuePO{
		ID: v.ID, AttributeID: v.AttributeID, Value: v.Value, SortOrder: v.SortOrder,
		CreatedAt: time.Time(v.CreatedAt), UpdatedAt: time.Time(v.UpdatedAt),
	}
}
