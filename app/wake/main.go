package main

import (
	"log"
	"os"
	"github.com/WPMedia/coral/wake/listener"
)

func init() {
	s := []string{"AMQP_URL", "AMQP_EXCHANGE", "ES_URL"}
	for _, one := range s {
		if os.Getenv(one) == "" {
			log.Fatalf("Error initializing Wake: %s not found.\n", one)
		}
	}
}

func main() {
	listener.Listen()
}

