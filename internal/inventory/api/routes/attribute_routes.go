package routes

import (
	"eshop-monolith/internal/infra/repository"
	"eshop-monolith/internal/inventory/api/handlers"
	"eshop-monolith/internal/inventory/domain/repositories"
	"eshop-monolith/internal/inventory/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAttributeRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	attrRepo := repositories.NewAttributeRepository(db)
	svc := service.NewAttributeService(attrRepo)
	h := handlers.NewAttributeHandler(svc)

	// 属性维度管理 /api/v1/attributes
	attrs := v1.Group("/attributes")
	{
		attrs.GET("", h.ListAttributes)
		attrs.POST("", h.CreateAttribute)
		attrs.GET("/:id", h.GetAttribute)
		attrs.PUT("/:id", h.UpdateAttribute)
		attrs.DELETE("/:id", h.DeleteAttribute)

		// 属性值列表作为子资源：/api/v1/attributes/:id/values
		attrs.GET("/:id/values", h.ListAttributeValues)
	}

	// 属性值增删改独立路由 /api/v1/attribute-values
	vals := v1.Group("/attribute-values")
	{
		vals.POST("", h.CreateAttributeValue)
		vals.PUT("/:id", h.UpdateAttributeValue)
		vals.DELETE("/:id", h.DeleteAttributeValue)
	}
}
