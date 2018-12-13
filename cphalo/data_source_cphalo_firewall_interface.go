package cphalo

import (
	"fmt"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCPHaloFirewallInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFirewallInterfaceRead,

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

func dataSourceFirewallInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	var (
		err               error
		client            = meta.(*api.Client)
		name              = d.Get("name").(string)
		interfaces        api.ListFirewallInterfacesResponse
		selectedInterface api.FirewallInterface
	)

	if interfaces, err = client.ListFirewallInterfaces(); err != nil {
		return err
	}

	for _, ifc := range interfaces.Interfaces {
		if strings.TrimSpace(ifc.Name) == name {
			selectedInterface = ifc
			break
		}
	}

	if selectedInterface.Name == "" {
		return fmt.Errorf("resource %s does not exists", name)
	}

	d.SetId(selectedInterface.ID)
	d.Set("name", selectedInterface.Name)

	return nil
}
