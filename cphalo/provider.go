package cphalo

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"application_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CP_APPLICATION_KEY", nil),
				Description: descriptions["application_key"],
			},
			"application_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CP_APPLICATION_SECRET", nil),
				Description: descriptions["application_secret"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cphalo_group": dataSourceCpHaloGroup(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"application_key":    "The CP API application key",
		"application_secret": "The CP API application secret",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		ApplicationKey:    d.Get("application_key").(string),
		ApplicationSecret: d.Get("application_secret").(string),
	}
	log.Println("[INFO] Initializing cphalo client")

	client := config.Client()

	ok, err := client.Validate()

	if err != nil {
		return client, err
	}

	if ok == false {
		return client, errors.New("No valid credential sources found for CPHalo Provider.")
	}

	return client, nil
}
