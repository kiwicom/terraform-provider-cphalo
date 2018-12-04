package cphalo

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"application_key":    "The CP API application key",
		"application_secret": "The CP API application secret",
	}
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"application_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CP_APPLICATION_KEY", nil),
				Description: descriptions["application_key"],
			},
			"application_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CP_APPLICATION_SECRET", nil),
				Description: descriptions["application_secret"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cphalo_server_group":    dataSourceCPHaloServerGroup(),
			"cphalo_firewall_policy": dataSourceCPHaloFirewallPolicy(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"cphalo_server_group":       resourceCPHaloServerGroup(),
			"cphalo_firewall_policy":    resourceCPHaloFirewallPolicy(),
			"cphalo_firewall_interface": resourceCPHaloFirewallInterface(),
			"cphalo_firewall_service":   resourceCPHaloFirewallService(),
			"cphalo_csp_account":        resourceCPHaloCSPAccount(),
		},

		ConfigureFunc: ConfigureProvider,
	}
}

func ConfigureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ApplicationKey:    d.Get("application_key").(string),
		ApplicationSecret: d.Get("application_secret").(string),
	}
	log.Println("[INFO] Initializing CPHalo client")

	client := config.Client()

	return client, nil
}
