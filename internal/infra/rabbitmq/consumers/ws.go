package consumers

import (
	"context"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/infra/ws"
)

func StartWSConsumer(ctx context.Context, client *rabbitmq.Client, hub *ws.Hub) error {
	consumer := rabbitmq.NewConsumer(client, rabbitmq.ConsumerConfig{
		Queue:    "eshop.ws-push",
		Bindings: []string{
			"order.paid", "order.shipped", "order.delivered", "order.cancelled",
			"payment.success", "payment.failed",
			"payment.refund.created", "payment.refund.status-updated", "payment.refund.failed",
			"flash-order.created", "flash-order.paid", "flash-order.cancelled",
			"inventory.low",
		},
	})
	consumer.HandleFunc("*", func(msg rabbitmq.Message) error {
		hub.HandleMessage(msg)
		return nil
	})
	return consumer.Start(ctx)
}
