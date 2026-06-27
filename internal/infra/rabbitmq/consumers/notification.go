package consumers

import (
	"context"

	"eshop-monolith/internal/infra/rabbitmq"
	"eshop-monolith/internal/notification/service"
)

func StartNotificationConsumer(ctx context.Context, client *rabbitmq.Client, notifSvc *service.NotificationService) error {
	consumer := rabbitmq.NewConsumer(client, rabbitmq.ConsumerConfig{
		Queue:    "eshop.notification",
		Bindings: []string{
			"order.paid", "order.shipped", "order.delivered", "order.cancelled",
			"flash-order.created", "flash-order.paid", "flash-order.cancelled",
			"payment.refund.created", "payment.refund.failed",
			"inventory.low",
		},
	})
	consumer.HandleFunc("*", func(msg rabbitmq.Message) error {
		return notifSvc.HandleMessage(msg)
	})
	return consumer.Start(ctx)
}
