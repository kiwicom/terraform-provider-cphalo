package cphalo

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CP_ENDPOINT", "https://api.cloudpassage.com/v1"),
				Description: descriptions["endpoint"],
			},
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
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"endpoint":           "The CP API endpoint url, default to https://api.cloudpassage.com/v1",
		"application_key":    "The CP API application key",
		"application_secret": "The CP API application secret",
	}
}
