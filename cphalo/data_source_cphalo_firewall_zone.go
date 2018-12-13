package cphalo

import (
	"fmt"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCPHaloFirewallZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFirewallZoneRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceFirewallZoneRead(d *schema.ResourceData, meta interface{}) error {
	var (
		err          error
		client       = meta.(*api.Client)
		name         = d.Get("name").(string)
		zones        api.ListFirewallZonesResponse
		selectedZone api.FirewallZone
	)

	if zones, err = client.ListFirewallZones(); err != nil {
		return err
	}

	for _, zone := range zones.Zones {
		if strings.TrimSpace(zone.Name) == name {
			selectedZone = zone
			break
		}
	}

	if selectedZone.Name == "" {
		return fmt.Errorf("resource %s does not exists", name)
	}

	d.SetId(selectedZone.ID)
	d.Set("name", selectedZone.Name)

	return nil
}
