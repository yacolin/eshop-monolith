package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/internal/cart/api/handlers"
	"eshop-monolith/internal/cart/service"
	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/repository"
	invService "eshop-monolith/internal/inventory/service"
)

// RegisterCartRoutes 注册购物车相关路由
func RegisterCartRoutes(router *gin.RouterGroup, repos *repository.Repositories, db *gorm.DB, rabbit *rabbitmq.Client) {
	// 创建库存服务实例
	inventoryService := invService.NewInventoryService(repos.Inventory, repos.Sku, repos.Product, rabbit)

	// 创建产品服务实例
	productService := invService.NewProductService(repos.Product, repos.Inventory, repos.Sku, rabbit, nil, repos.Redis)

	// 创建 SKU 服务实例
	skuService := invService.NewSkuService(repos.Sku, repos.Product, rabbit, db)

	// 创建购物车服务实例
	cartService := service.NewCartService(repos.Cart, inventoryService, productService, skuService, rabbit)

	// 创建购物车处理器实例
	cartHandler := handlers.NewCartHandler(cartService)

	carts := router.Group("/carts")
	{
		// 获取购物车
		carts.GET("", cartHandler.GetCart)
		// 清空购物车
		carts.DELETE("", cartHandler.ClearCart)

		// 购物车项管理
		items := carts.Group("/items")
		{
			// 添加商品到购物车
			items.POST("", cartHandler.AddToCart)
			// 更新购物车项
			items.PUT("/:item_id", cartHandler.UpdateCartItem)
			// 删除购物车项
			items.DELETE("/:item_id", cartHandler.RemoveCartItem)
		}
	}
}
