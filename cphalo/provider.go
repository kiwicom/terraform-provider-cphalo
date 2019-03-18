package cphalo

import (
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

// Provider returns CPHalo Resource provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"application_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPHALO_APPLICATION_KEY", nil),
				Description: descriptions["application_key"],
			},
			"application_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CPHALO_APPLICATION_SECRET", nil),
				Description: descriptions["application_secret"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cphalo_server_group":       dataSourceCPHaloServerGroup(),
			"cphalo_firewall_policy":    dataSourceCPHaloFirewallPolicy(),
			"cphalo_firewall_interface": dataSourceCPHaloFirewallInterface(),
			"cphalo_firewall_service":   dataSourceCPHaloFirewallService(),
			"cphalo_firewall_zone":      dataSourceCPHaloFirewallZone(),
			"cphalo_alert_profile":      dataSourceCPHaloAlertProfile(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"cphalo_server_group":                 resourceCPHaloServerGroup(),
			"cphalo_server_group_firewall_policy": resourceCPHaloServerGroupFirewallPolicy(),
			"cphalo_firewall_policy":              resourceCPHaloFirewallPolicy(),
			"cphalo_firewall_interface":           resourceCPHaloFirewallInterface(),
			"cphalo_firewall_service":             resourceCPHaloFirewallService(),
			"cphalo_firewall_zone":                resourceCPHaloFirewallZone(),
			"cphalo_csp_aws_account":              resourceCPHaloCSPAWSAccount(),
		},

		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	c := config{
		applicationKey:    d.Get("application_key").(string),
		applicationSecret: d.Get("application_secret").(string),
	}
	logInfo("Initializing CPHalo client")

	client := c.client()

	return client, nil
}
