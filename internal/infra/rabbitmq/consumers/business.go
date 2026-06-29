package consumers

import (
	"context"
	"encoding/json"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/trade"

	"eshop-monolith/internal/base"
)

func StartBusinessConsumer(ctx context.Context, client *rabbitmq.Client, nSvc *base.NotificationService) error {
	consumer := rabbitmq.NewConsumer(client, rabbitmq.ConsumerConfig{
		Queue: "eshop.business",
		Bindings: []string{
			"order.paid", "order.shipped", "order.delivered", "order.cancelled",
			"payment.success", "payment.failed",
			"payment.refund.created", "payment.refund.status-updated", "payment.refund.failed",
			"flash-order.created", "flash-order.paid", "flash-order.cancelled",
			"inventory.low",
		},
	})

	consumer.HandleFunc("*", func(msg rabbitmq.Message) error {
		switch msg.RoutingKey {
		case "payment.success":
			var e trade.PaymentSuccessEvent
			if err := json.Unmarshal(msg.Payload, &e); err != nil {
				return err
			}
			if e.OrderType != "flash" {
			}
		case "flash-order.paid":
			var e trade.PaymentSuccessEvent
			if err := json.Unmarshal(msg.Payload, &e); err != nil {
				return err
			}
		}
		return nSvc.HandleMessage(msg)
	})

	return consumer.Start(ctx)
}
