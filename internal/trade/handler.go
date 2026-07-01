package trade

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"eshop-monolith/internal/infra/rabbitmq"
)

func RegisterTradeRoutes(v1 *gin.RouterGroup, db *gorm.DB, mqClient *rabbitmq.Client, rdb *redis.Client, invalidateCache ...func()) {
	RegisterCartRoutes(v1, db, rdb)
	RegisterPaymentRoutes(v1, db)
	RegisterOrderRoutes(v1, db, mqClient, invalidateCache...)
}
