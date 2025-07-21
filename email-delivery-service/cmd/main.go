package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/compendium-tech/compendium/email-delivery-service/internal/config"
	"github.com/compendium-tech/compendium/email-delivery-service/pkg/domain"
	"github.com/compendium-tech/compendium/email-delivery-service/pkg/sender"
	"github.com/joho/godotenv"
	kafka "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file, using environmental variables instead: %v\n", err)
	}

	cfg := config.LoadAppConfig()
	smtpSender := sender.NewSmtpEmailSender(cfg.SmtpHost, cfg.SmtpPort, cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpFrom)
	consumeAndSendEmails(ctx, cfg, smtpSender)
}

func consumeAndSendEmails(
	ctx context.Context,
	cfg config.AppConfig, sender sender.EmailSender) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBroker},
		Topic:    cfg.KafkaTopic,
		GroupID:  cfg.KafkaGroupId,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})
	defer reader.Close()

	for {
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			logrus.Printf("Error fetching message from Kafka: %v", err)
			time.Sleep(5 * time.Second)

			continue
		}

		logrus.Printf("Received message from Kafka - Topic: %s, Partition: %d, Offset: %d, Key: %s",
			m.Topic, m.Partition, m.Offset, string(m.Key))

		var emailMsg domain.EmailMessage

		if err := json.Unmarshal(m.Value, &emailMsg); err != nil {
			logrus.Printf("Error unmarshaling email message: %v, Message value: %s", err, string(m.Value))

			continue
		}

		if err := sender.SendMessage(ctx, emailMsg); err != nil {
			logrus.Errorf("Error sending email to %s: %v", emailMsg.To, err)
		} else {
			logrus.Printf("Email sent successfully to %s, Subject: %s", emailMsg.To, emailMsg.Subject)
		}

		if err := reader.CommitMessages(ctx, m); err != nil {
			logrus.Printf("Error committing message: %v", err)
		}
	}
}
