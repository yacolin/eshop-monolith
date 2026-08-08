package staff

import "github.com/gin-gonic/gin"

// currentStaffID 从 JWTAuth 写入的 context 取当前员工 ID(int64)
func currentStaffID(c *gin.Context) int64 {
	v, ok := c.Get("user_id")
	if !ok {
		return 0
	}
	id, ok := v.(int64)
	if !ok {
		return 0
	}
	return id
}
