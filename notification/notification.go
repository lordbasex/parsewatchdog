package notification

import (
	"log"

	"github.com/lordbasex/parsewatchdog/config"
)

func NotifyAll(cfg *config.Config, subject, message string) {
	if cfg.SMTP.Enabled {
		if err := NewEmailNotifier(cfg).Send(subject, message); err != nil {
			log.Println("Error sending email:", err)
		}
	}
	if cfg.Telegram.Enabled {
		if err := NewTelegramNotifier(cfg).Send(subject, message); err != nil {
			log.Println("Error sending Telegram message:", err)
		}
	}
	if cfg.API.Enabled {
		if err := NewAPINotifier(cfg).Send(subject, message); err != nil {
			log.Println("Error sending API notification:", err)
		}
	}
	if cfg.RabbitMQ.Enabled {
		if err := NewRabbitMQNotifier(cfg).Send(subject, message); err != nil {
			log.Println("Error sending RabbitMQ notification:", err)
		}
	}
}
