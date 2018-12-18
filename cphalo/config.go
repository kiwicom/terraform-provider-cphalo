package cphalo

import (
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
)

type Config struct {
	ApplicationKey    string
	ApplicationSecret string
}

func (c *Config) Client() *api.Client {

	client := api.NewClient(c.ApplicationKey, c.ApplicationSecret)

	logInfo("CP Client configured.")

	return client
}
