package email

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type kafkaEmailMessageProducer struct {
	writer *kafka.Writer
}

func NewKafkaEmailMessageProducer(broker, topic string) Sender {
	return &kafkaEmailMessageProducer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(broker),
			Topic: topic,
		},
	}
}

func (kp *kafkaEmailMessageProducer) SendMessage(msg Message) {
	messageBytes, err := json.Marshal(msg)

	if err != nil {
		panic(fmt.Errorf("failed to marshal email message to JSON: %w", err))
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(msg.To),
		Value: messageBytes,
	}

	err = kp.writer.WriteMessages(context.Background(), kafkaMsg)
	if err != nil {
		panic(fmt.Errorf("failed to write message to Kafka: %w", err))
	}
}
