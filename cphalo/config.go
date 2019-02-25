package cphalo

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/logging"
	"gitlab.com/kiwicom/cphalo-go"
)

type config struct {
	applicationKey    string
	applicationSecret string
}

func (c *config) client() *cphalo.Client {

	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	transport := logging.NewTransport("CPHalo", t)

	httpClient := &http.Client{Transport: transport}

	client := cphalo.NewClient(c.applicationKey, c.applicationSecret, httpClient)

	logInfo("CP Client configured.")

	return client
}
