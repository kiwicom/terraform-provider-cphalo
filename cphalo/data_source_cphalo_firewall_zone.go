package cphalo

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
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
			"ip_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceFirewallZoneRead(d *schema.ResourceData, meta interface{}) error {
	var (
		err          error
		client       = meta.(*cphalo.Client)
		name         = d.Get("name").(string)
		zones        cphalo.ListFirewallZonesResponse
		selectedZone cphalo.FirewallZone
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
	_ = d.Set("name", selectedZone.Name)
	_ = d.Set("ip_address", selectedZone.IPAddress)
	_ = d.Set("description", selectedZone.Description)

	return nil
}
