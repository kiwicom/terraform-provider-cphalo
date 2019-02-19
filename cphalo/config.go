package cphalo

import (
	"gitlab.com/kiwicom/cphalo-go"
)

type config struct {
	applicationKey    string
	applicationSecret string
}

func (c *config) client() *cphalo.Client {

	client := cphalo.NewClient(c.applicationKey, c.applicationSecret)

	logInfo("CP Client configured.")

	return client
}
