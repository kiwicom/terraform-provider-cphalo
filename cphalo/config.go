package cphalo

import (
	"gitlab.com/kiwicom/cphalo-go"
)

type Config struct {
	ApplicationKey    string
	ApplicationSecret string
}

func (c *Config) Client() *cphalo.Client {

	client := cphalo.NewClient(c.ApplicationKey, c.ApplicationSecret)

	logInfo("CP Client configured.")

	return client
}
