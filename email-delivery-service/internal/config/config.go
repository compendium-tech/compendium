package config

import (
	"fmt"
	"log"
	"os"
)

type AppConfig struct {
	SmtpHost     string
	SmtpPort     uint16
	SmtpUsername string
	SmtpPassword string
	SmtpFrom     string
	KafkaBroker  string
	KafkaTopic   string
	KafkaGroupID string
}

func LoadAppConfig() AppConfig {
	appConfig := AppConfig{
		KafkaBroker:  os.Getenv("KAFKA_BROKER"),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
		SmtpHost:     os.Getenv("SMTP_HOST"),
		SmtpUsername: os.Getenv("SMTP_USERNAME"),
		SmtpPassword: os.Getenv("SMTP_PASSWORD"),
		SmtpFrom:     os.Getenv("SMTP_FROM"),
	}

	if port := os.Getenv("SMTP_PORT"); port != "" {
		var smtpPort uint16
		_, err := fmt.Sscan(port, &smtpPort)

		if err == nil {
			appConfig.SmtpPort = smtpPort
		} else {
			log.Printf("Failed to parse smtp port: %s", port)
		}
	}

	return appConfig
}
