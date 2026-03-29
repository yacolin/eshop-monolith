package routes

import (
	"errors"
	"eshop-monolith/internal/pkg/config"
	"eshop-monolith/internal/pkg/response"
	"eshop-monolith/internal/repository"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, repos *repository.Repositories) *gin.Engine {
	router := gin.Default()

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		// c.JSON(200, gin.H{
		// 	"status": "ok",
		// })
		response.Success(c, gin.H{
			"status": "ok",
		})
	})

	// API v1 路由组
	v1 := router.Group("/api/v1")
	{
		// 健康检查
		v1.GET("/health", func(c *gin.Context) {
			response.Success(c, gin.H{
				"status":  "ok",
				"version": "v1",
			})
		})

		// 错误返回测试
		v1.GET("/error", func(c *gin.Context) {
			response.SysError(c, errors.New("test error"))
		})

		// 这里将添加具体的路由
		// 例如：订单、产品、用户等路由
	}

	return router
}
