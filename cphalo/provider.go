package cphalo

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
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
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"application_key":    "The CP API application key",
		"application_secret": "The CP API application secret",
	}
}
