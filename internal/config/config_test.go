package config

import (
	"fmt"
	"log"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("thank you for testing config, this is crucial to launch our application\n")
	fmt.Printf("your listening on port -> %q\n", config.Port)
	fmt.Printf("your database url is -> %q\n", config.DatabaseURL)
}
