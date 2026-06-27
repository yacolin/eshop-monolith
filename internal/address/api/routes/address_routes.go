package routes

import (
	"eshop-monolith/internal/address/api/handlers"
	"eshop-monolith/internal/address/domain/repositories"
	"eshop-monolith/internal/address/service"
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAddressRoutes(v1 *gin.RouterGroup, db *gorm.DB, rabbit *rabbitmq.Client) {
	repo := repositories.NewAddressRepository(db)
	svc := service.NewAddressService(repo, db, rabbit)
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
