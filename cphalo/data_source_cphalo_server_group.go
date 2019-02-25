package cphalo

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
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
				Computed: true,
			},
		},
	}
}

func dataSourceCPHaloServerGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cphalo.Client)

	gs, err := client.ListServerGroups()
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	var selectedGroup cphalo.ServerGroup
	for _, g := range gs.Groups {
		if strings.TrimSpace(g.Name) == name {
			selectedGroup = g
			break
		}
	}

	if selectedGroup.Name == "" {
		return fmt.Errorf("resouce %s does not exists", name)
	}

	d.SetId(selectedGroup.ID)
	_ = d.Set("name", selectedGroup.Name)
	_ = d.Set("parent_id", selectedGroup.ParentID)

	return nil
}
