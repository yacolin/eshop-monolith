package routes

import (
	"eshop-monolith/internal/flashsale/api/handlers"
	"eshop-monolith/internal/flashsale/domain/repositories"
	"eshop-monolith/internal/flashsale/service"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterFlashRoutes(v1 *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB) {
	flashRepo := repositories.NewFlashRepository(db)
	if err := flashRepo.AutoMigrate(); err != nil {
		panic("failed to auto migrate flash tables: " + err.Error())
	}

	flashService := service.NewFlashService(db, repos.Redis, flashRepo)
	flashHandler := handlers.NewFlashHandler(flashService)

	flash := v1.Group("/flash")
	{
		flash.POST("/activities", flashHandler.CreateActivity)
		flash.POST("/activities/:id/load-stock", flashHandler.LoadStock)
		flash.POST("/buy", flashHandler.FlashBuy)
		flash.GET("/activities", flashHandler.ListActivities)
		flash.GET("/activities/:id", flashHandler.GetActivity)
		flash.GET("/orders/:id", flashHandler.GetOrder)
		flash.GET("/users/:user_id/orders", flashHandler.GetUserOrders)
	}
}