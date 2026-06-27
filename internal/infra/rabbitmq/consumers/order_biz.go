package consumers

import (
	"context"
	"encoding/json"

	"eshop-monolith/internal/infra/rabbitmq"
	flashSvc "eshop-monolith/internal/flashsale/service"
	orderSvc "eshop-monolith/internal/order/service"
	paymentEvents "eshop-monolith/internal/payment/events"
)

func StartOrderBizConsumers(ctx context.Context, client *rabbitmq.Client, oSvc *orderSvc.OrderService, fSvc *flashSvc.FlashService) error {
	// order-svc: 消费 payment.success
	orderConsumer := rabbitmq.NewConsumer(client, rabbitmq.ConsumerConfig{
		Queue:    "eshop.order-svc",
		Bindings: []string{"payment.success"},
	})
	orderConsumer.HandleFunc("payment.success", func(msg rabbitmq.Message) error {
		var e paymentEvents.PaymentSuccessEvent
		json.Unmarshal(msg.Payload, &e)
		if e.OrderType != "flash" {
			return oSvc.HandlePaidSuccess(context.Background(), e.OrderID)
		}
		return nil
	})

	// flash-svc: 消费 flash-order.paid
	flashConsumer := rabbitmq.NewConsumer(client, rabbitmq.ConsumerConfig{
		Queue:    "eshop.flash-svc",
		Bindings: []string{"flash-order.paid"},
	})
	flashConsumer.HandleFunc("flash-order.paid", func(msg rabbitmq.Message) error {
		var e paymentEvents.PaymentSuccessEvent
		json.Unmarshal(msg.Payload, &e)
		return fSvc.HandlePaidSuccess(context.Background(), e.OrderID)
	})

	if err := orderConsumer.Start(ctx); err != nil {
		return err
	}
	return flashConsumer.Start(ctx)
}
