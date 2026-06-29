package inventory

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eshop-monolith/pkg/middleware"
	"eshop-monolith/pkg/response"
)

type InventoryHandler struct {
	svc *InventoryService
}

func NewInventoryHandler(svc *InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

// Lock 下单预占库存
// @Summary 下单预占库存
// @Tags inventories
// @Accept json
// @Produce json
// @Param req body LockStockReq true "预占请求"
// @Success 200 {object} response.Response{data=Inventory}
// @Router /api/v1/inventories/lock [post]
func (h *InventoryHandler) Lock(c *gin.Context) {
	var req LockStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Lock(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Unlock 取消释放库存
// @Summary 取消释放库存
// @Tags inventories
// @Accept json
// @Produce json
// @Param req body UnlockStockReq true "释放请求"
// @Success 200 {object} response.Response{data=Inventory}
// @Router /api/v1/inventories/unlock [post]
func (h *InventoryHandler) Unlock(c *gin.Context) {
	var req UnlockStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Unlock(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Deduct 支付扣减库存
// @Summary 支付扣减库存
// @Tags inventories
// @Accept json
// @Produce json
// @Param req body DeductStockReq true "扣减请求"
// @Success 200 {object} response.Response{data=Inventory}
// @Router /api/v1/inventories/deduct [post]
func (h *InventoryHandler) Deduct(c *gin.Context) {
	var req DeductStockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Deduct(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// Restock 入库/补货
// @Summary 入库/补货
// @Tags inventories
// @Accept json
// @Produce json
// @Param req body RestockReq true "入库请求"
// @Success 200 {object} response.Response{data=Inventory}
// @Router /api/v1/inventories/restock [post]
func (h *InventoryHandler) Restock(c *gin.Context) {
	var req RestockReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.Restock(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// GetStock 查询库存
// @Summary 查询库存
// @Tags inventories
// @Produce json
// @Param sku_id query int true "SKU ID"
// @Success 200 {object} response.Response{data=Inventory}
// @Router /api/v1/inventories/stock [get]
func (h *InventoryHandler) GetStock(c *gin.Context) {
	var req struct {
		SkuID int64 `form:"sku_id" binding:"required"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.GetStock(c, req.SkuID)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListLogs 查询库存变更流水
// @Summary 查询库存变更流水
// @Tags inventories
// @Produce json
// @Param sku_id query int true "SKU ID"
// @Param change_type query string false "变更类型"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页条数" default(20)
// @Success 200 {object} response.Response{data=InventoryLogListResult}
// @Router /api/v1/inventories/logs [get]
func (h *InventoryHandler) ListLogs(c *gin.Context) {
	var req InventoryLogQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		c.Error(err)
		return
	}
	result, err := h.svc.ListLogs(c, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, result)
}

// ── Routes ────────────────────────────────────────

func RegisterInventoryRoutes(v1 *gin.RouterGroup, db *gorm.DB) {
	repo := NewInventoryRepository(db)
	svc := NewInventoryService(repo, db)
	h := NewInventoryHandler(svc)

	inv := v1.Group("/inventories")
	{
		inv.GET("/stock", h.GetStock)
		inv.GET("/logs", h.ListLogs)
	}
	auth := v1.Group("/inventories")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/lock", h.Lock)
		auth.POST("/unlock", h.Unlock)
		auth.POST("/deduct", h.Deduct)
		auth.POST("/restock", h.Restock)
	}
}
