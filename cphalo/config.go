package cphalo

import (
	"log"
	"os"
)

type Config struct {
	ApplicationKey    string
	ApplicationSecret string
}

func (c *Config) Client() (*Client, error) {

	if v := os.Getenv("CP_APPLICATION_KEY"); v != "" {
		c.ApplicationKey = v
	}

	if v := os.Getenv("CP_APPLICATION_SECRET"); v != "" {
		c.ApplicationSecret = v
	}

	client := newClient(c.ApplicationKey, c.ApplicationSecret)

	log.Printf("[INFO] CP Client configured.")

	return client, nil
}
