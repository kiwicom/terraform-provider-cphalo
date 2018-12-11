package cphalo

import (
	"fmt"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCPHaloServerGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCPHaloServerGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceCPHaloServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	gs, err := client.ListServerGroups()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	var selectedGroup api.ServerGroup
	for _, g := range gs.Groups {
		if strings.TrimSpace(g.Name) == name {
			selectedGroup = g
			break
		}
	}

	if selectedGroup.Name == "" {
		return fmt.Errorf("resouce %s does not exists", name)
	}

	log.Println(selectedGroup)

	d.SetId(selectedGroup.ID)
	d.Set("name", selectedGroup.Name)
	d.Set("parent_id", selectedGroup.ParentID)

	return nil
}
