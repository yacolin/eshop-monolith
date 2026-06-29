package trade

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterTradeRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	RegisterCartRoutes(v1, db)
	RegisterPaymentRoutes(v1, db)
	RegisterOrderRoutes(v1, db)
}
