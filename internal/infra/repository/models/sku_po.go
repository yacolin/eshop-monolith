package models

import (
	"time"

	domain "eshop-monolith/internal/inventory/domain/models"
	"eshop-monolith/pkg/utils"

	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

// ── SkuPO ────────────────────────────────────────────────────────────

type SkuPO struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	ProductID int64          `gorm:"index;not null"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Price     int64          `gorm:"type:bigint;not null"`
	SKUCode   string         `gorm:"type:varchar(100);uniqueIndex;not null"`
	Image     string         `gorm:"type:varchar(500);default:''"`
	Spec      string         `gorm:"type:json"`
	CreatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP()"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP();onUpdate:CURRENT_TIMESTAMP()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SkuPO) TableName() string { return "skus" }

func (po *SkuPO) ToDomain() *domain.Sku {
	var specMap map[string]string
	if po.Spec != "" {
		sonic.Unmarshal([]byte(po.Spec), &specMap)
	}
	return &domain.Sku{
		ID:        po.ID,
		ProductID: po.ProductID,
		Name:      po.Name,
		Price:     po.Price,
		SKUCode:   po.SKUCode,
		Image:     po.Image,
		Spec:      specMap,
		CreatedAt: utils.Timestamp(po.CreatedAt),
		UpdatedAt: utils.Timestamp(po.UpdatedAt),
	}
}

func SkuFromDomain(s *domain.Sku) *SkuPO {
	var specStr string
	if len(s.Spec) > 0 {
		data, _ := sonic.Marshal(s.Spec)
		specStr = string(data)
	}
	return &SkuPO{
		ID:        s.ID,
		ProductID: s.ProductID,
		Name:      s.Name,
		Price:     s.Price,
		SKUCode:   s.SKUCode,
		Image:     s.Image,
		Spec:      specStr,
		CreatedAt: time.Time(s.CreatedAt),
		UpdatedAt: time.Time(s.UpdatedAt),
	}
}

// ── SkuAttributePO ──────────────────────────────────────────────────

type SkuAttributePO struct {
	SkuID           int64 `gorm:"primaryKey"`
	AttributeID     int64 `gorm:"primaryKey"`
	AttributeValueID int64 `gorm:"not null"`
}

func (SkuAttributePO) TableName() string { return "sku_attributes" }

func (po *SkuAttributePO) ToDomain() *domain.SkuAttribute {
	return &domain.SkuAttribute{
		SkuID: po.SkuID, AttributeID: po.AttributeID, AttributeValueID: po.AttributeValueID,
	}
}

func SkuAttributeFromDomain(sa *domain.SkuAttribute) *SkuAttributePO {
	return &SkuAttributePO{
		SkuID: sa.SkuID, AttributeID: sa.AttributeID, AttributeValueID: sa.AttributeValueID,
	}
}
