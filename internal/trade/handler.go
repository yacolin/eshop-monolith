package trade

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/internal/infra/rabbitmq"
)

func RegisterTradeRoutes(v1 *gin.RouterGroup, db *gorm.DB, mqClient *rabbitmq.Client) {
	RegisterCartRoutes(v1, db)
	RegisterPaymentRoutes(v1, db)
	RegisterOrderRoutes(v1, db, mqClient)
}
