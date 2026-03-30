package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseIntParam 解析路径参数为int64
func ParseIntParam(c *gin.Context, param string) (int64, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
