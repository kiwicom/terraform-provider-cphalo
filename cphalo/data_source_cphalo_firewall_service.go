package cphalo

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
)

func dataSourceCPHaloFirewallService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFirewallServiceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFirewallServiceRead(d *schema.ResourceData, meta interface{}) error {
	var (
		err             error
		client          = meta.(*cphalo.Client)
		name            = d.Get("name").(string)
		services        cphalo.ListFirewallServicesResponse
		selectedService cphalo.FirewallService
	)

	if services, err = client.ListFirewallServices(); err != nil {
		return err
	}

	for _, service := range services.Services {
		if strings.TrimSpace(service.Name) == name {
			selectedService = service
			break
		}
	}

	if selectedService.Name == "" {
		return fmt.Errorf("resource %s does not exists", name)
	}

	d.SetId(selectedService.ID)
	_ = d.Set("name", selectedService.Name)
	_ = d.Set("protocol", selectedService.Protocol)
	_ = d.Set("port", selectedService.Port)

	return nil
}
