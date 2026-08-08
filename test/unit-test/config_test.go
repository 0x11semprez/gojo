package config

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg, err := NewConfig()
	if err != nil {
		t.Fatal(err)
	}

	cfgtest, err := NewConfigTest()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("thank you for testing config, this is crucial to launch our application\n")
	fmt.Printf("your listening on port -> %q\n", cfg.Port)
	fmt.Printf("your database url is -> %q\n", cfg.DatabaseURLTest)
	fmt.Printf("your database url for migrations is -> %q\n", cfg.DatabaseURLMigration)
	fmt.Printf("your database url for tests is -> %q\n", cfgtest.DatabaseURL)
	fmt.Printf("your database url for migrations tests is -> %q\n", cfgtest.DatabaseURLMigration)
}
