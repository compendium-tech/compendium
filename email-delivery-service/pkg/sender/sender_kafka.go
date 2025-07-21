package sender

import (
	"context"
	"encoding/json"

	"github.com/compendium-tech/compendium/common/pkg/log"

	"github.com/compendium-tech/compendium/email-delivery-service/pkg/domain"
	"github.com/segmentio/kafka-go"
	"github.com/ztrue/tracerr"
)

type kafkaEmailMessageProducer struct {
	writer *kafka.Writer
}

func NewKafkaEmailMessageProducer(broker, topic string) EmailSender {
	return &kafkaEmailMessageProducer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{broker},
			Topic:   topic,
		}),
	}
}

func (kp *kafkaEmailMessageProducer) SendMessage(ctx context.Context, msg domain.EmailMessage) error {
	messageBytes, err := json.Marshal(msg)

	if err != nil {
		return tracerr.Errorf("failed to marshal email message to JSON: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.To),
		Value: messageBytes,
	}

	log.L(ctx).Printf("Attempting to produce message for recipient: %s, Subject: %s", msg.To, msg.Subject)

	err = kp.writer.WriteMessages(context.Background(), kafkaMsg)
	if err != nil {
		return tracerr.Errorf("failed to write message to Kafka: %w", err)
	}

	log.L(ctx).Printf("Produced message for recipient: %s, Subject: %s", msg.To, msg.Subject)

	return nil
}
