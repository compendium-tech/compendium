package email

import (
	"context"
	"encoding/json"

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

func (kp *kafkaEmailMessageProducer) SendMessage(msg EmailMessage) error {
	messageBytes, err := json.Marshal(msg)

	if err != nil {
		return tracerr.Errorf("failed to marshal email message to JSON: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.To),
		Value: messageBytes,
	}

	err = kp.writer.WriteMessages(context.Background(), kafkaMsg)
	if err != nil {
		return tracerr.Errorf("failed to write message to Kafka: %w", err)
	}

	return nil
}
