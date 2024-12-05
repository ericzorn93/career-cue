package api

import (
	"context"
	"log"
	"time"

	commonv1 "packages/proto-gen/go/common/v1"
	pb "packages/proto-gen/go/webhooks/inboundwebhooksapi/v1"

	amqp "github.com/rabbitmq/amqp091-go"
)

// UserRegistered handles incoming Webhooks from Auth0 and will attach the message
// to an exchange within the message broker
func (s *InboundWebhooksAPIServer) UserRegistered(
	_ context.Context,
	req *pb.UserRegisteredRequest,
) (*commonv1.Empty, error) {
	log.Println("hit grpc endpoint")
	log.Println(req)

	connection, _ := amqp.Dial("amqp://guest:guest@lavinmq:5672")
	defer connection.Close()

	channel, _ := connection.Channel()
	log.Print("[✅] Connection over channel established")
	defer channel.Close()

	queue, _ := channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World"
	channel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf("[📥] Message sent to queue:  %s\n", body)

	return &commonv1.Empty{}, nil
}
