package service

type OrderService struct{}

func (s *OrderService) HandlePaidSuccess(ctx interface{}, orderID int64) error {
	return nil
}
