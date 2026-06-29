package utils

import (
	"fmt"
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

// ParseQueryIntParam 解析查询参数为int64
func ParseQueryIntParam(c *gin.Context, param string) (int64, error) {
	idStr := c.Query(param)
	if idStr == "" {
		return 0, fmt.Errorf("missing query parameter: %s", param)
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
