package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Ch *amqp.Channel

const OrderStatusQueue = "order_status_updates"

// OrderStatusEvent is the message we publish when an order status changes
type OrderStatusEvent struct {
	OrderID  uint   `json:"order_id"`
	UserID   uint   `json:"user_id"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

// ConnectRabbitMQ establishes connection to RabbitMQ
func ConnectRabbitMQ(url string) error {
	var err error

	Conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	Ch, err = Conn.Channel()
	if err != nil {
		return err
	}

	// Declare the queue (creates it if it doesn't exist)
	_, err = Ch.QueueDeclare(
		OrderStatusQueue, // queue name
		true,             // durable (survives broker restart)
		false,            // auto-delete
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	log.Println("RabbitMQ connected and queue declared successfully!")
	return nil
}

// PublishOrderStatusEvent publishes an order status update to the queue
func PublishOrderStatusEvent(ctx context.Context, event OrderStatusEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return Ch.PublishWithContext(
		ctx,
		"",               // exchange (default)
		OrderStatusQueue, // routing key (queue name)
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Message survives broker restart
		},
	)
}

// Close cleans up connections
func Close() {
	if Ch != nil {
		Ch.Close()
	}
	if Conn != nil {
		Conn.Close()
	}
}
