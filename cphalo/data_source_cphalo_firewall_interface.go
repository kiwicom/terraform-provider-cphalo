package cphalo

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
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
		client            = meta.(*cphalo.Client)
		name              = d.Get("name").(string)
		interfaces        cphalo.ListFirewallInterfacesResponse
		selectedInterface cphalo.FirewallInterface
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
	_ = d.Set("name", selectedInterface.Name)

	return nil
}
