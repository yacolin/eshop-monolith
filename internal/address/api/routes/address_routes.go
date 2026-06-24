package routes

import (
	"eshop-monolith/internal/address/api/handlers"
	"eshop-monolith/internal/address/domain/repositories"
	"eshop-monolith/internal/address/service"
	"eshop-monolith/internal/infra/eventbus"
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAddressRoutes(v1 *gin.RouterGroup, db *gorm.DB, bus *eventbus.Bus) {
	repo := repositories.NewAddressRepository(db)
	svc := service.NewAddressService(repo, db, bus)
	handler := handlers.NewAddressHandler(svc)

	addresses := v1.Group("/addresses")
	addresses.Use(middleware.JWTAuth())
	{
		addresses.POST("", handler.Create)
		addresses.GET("", handler.List)
		addresses.GET("/default", handler.GetDefault)
		addresses.GET("/:id", handler.Get)
		addresses.PUT("/:id", handler.Update)
		addresses.DELETE("/:id", handler.Delete)
	}
}
