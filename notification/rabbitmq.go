package notification

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/lordbasex/parsewatchdog/config"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQNotifier struct {
	config *config.Config
}

// NewRabbitMQNotifier initializes the RabbitMQ notifier
func NewRabbitMQNotifier(cfg *config.Config) *RabbitMQNotifier {
	return &RabbitMQNotifier{config: cfg}
}

// Send sends a JSON message to a RabbitMQ queue

func (n *RabbitMQNotifier) Send(subject, message string) error {
	// Corregir la cadena de conexión usando %d para el puerto si es un entero
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/",
		n.config.RabbitMQ.Type,
		n.config.RabbitMQ.User,
		n.config.RabbitMQ.Password,
		n.config.RabbitMQ.IP,
		n.config.RabbitMQ.Port,
	)
	//log.Printf("Connecting to RabbitMQ with: %s", connStr) // Log de conexión

	// Attempt to connect to RabbitMQ
	conn, err := amqp091.Dial(connStr)
	if err != nil {
		log.Printf("Connection error: %v", err) // Log detailed error
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	// Attempt to open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Channel error: %v", err) // Log detailed error
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Declare a queue (idempotent, only creates if not exists)
	q, err := ch.QueueDeclare(
		n.config.RabbitMQ.Queue, // queue name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		log.Printf("Queue declaration error: %v", err) // Log detailed error
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Create the JSON message payload
	payload := map[string]string{
		"subject": subject,
		"message": message,
	}
	jsonMessage, err := json.Marshal(payload)
	if err != nil {
		log.Printf("JSON encoding error: %v", err) // Log detailed error
		return fmt.Errorf("failed to encode message to JSON: %w", err)
	}

	// Attempt to publish the message to the queue
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonMessage,
		},
	)
	if err != nil {
		log.Printf("Publishing error: %v", err) // Log detailed error
		return fmt.Errorf("failed to publish message: %w", err)
	}

	//log.Printf("Message sent to RabbitMQ queue %s: %s", q.Name, string(jsonMessage))
	return nil
}
