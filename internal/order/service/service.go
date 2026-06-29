package service

// OrderService 桩（旧模块引用兼容）
type OrderService struct{}

func (s *OrderService) HandlePaidSuccess(ctx interface{}, orderID int64) error {
	return nil
}
