package netapp

import (
	"fmt"
	"os"
)

type App interface {
	Run() error
}

func getPortEnv() (uint16, error) {
	if portString := os.Getenv("PORT"); portString != "" {
		var port uint16
		_, err := fmt.Sscan(portString, &port)

		if err != nil {
			return 0, fmt.Errorf("failed to parse port env variable")
		}

		return port, nil
	}

	return 0, fmt.Errorf("no port env variable found")
}
