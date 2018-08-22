package cphalo

import (
	"log"
)

type Config struct {
	ApplicationKey    string
	ApplicationSecret string
}

func (c *Config) Client() *Client {

	client := newClient(c.ApplicationKey, c.ApplicationSecret)

	log.Printf("[INFO] CP Client configured.")

	return client
}
