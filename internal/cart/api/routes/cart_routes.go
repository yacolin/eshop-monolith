package routes

import (
	"eshop-monolith/internal/cart/api/handlers"
	"eshop-monolith/internal/cart/service"
	"eshop-monolith/internal/infra/eventbus"
	invService "eshop-monolith/internal/inventory/service"
	"eshop-monolith/internal/infra/repository"

	"github.com/gin-gonic/gin"
)

// RegisterCartRoutes 注册购物车相关路由
func RegisterCartRoutes(router *gin.RouterGroup, repos *repository.Repositories) {
	// 创建事件总线实例
	bus := eventbus.NewBus()

	// 创建库存服务实例
	inventoryService := invService.NewInventoryService(repos.Inventory, bus)

	// 创建产品服务实例
	productService := invService.NewProductService(repos.Product, repos.Inventory, bus, nil, repos.Redis)

	// 创建购物车服务实例
	cartService := service.NewCartService(repos.Cart, inventoryService, productService, bus)

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
