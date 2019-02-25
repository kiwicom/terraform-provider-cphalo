package cphalo

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
)

func dataSourceCPHaloAlertProfile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlertProfileRead,

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

func dataSourceAlertProfileRead(d *schema.ResourceData, meta interface{}) error {
	var (
		err      error
		client   = meta.(*cphalo.Client)
		name     = d.Get("name").(string)
		resp     cphalo.ListAlertProfilesResponse
		selected cphalo.AlertProfile
	)

	if resp, err = client.ListAlertProfiles(); err != nil {
		return err
	}

	for _, ap := range resp.AlertProfiles {
		if ap.Name == name {
			selected = ap
			break
		}
	}

	if selected.Name == "" {
		return fmt.Errorf("resource %s does not exists", name)
	}

	d.SetId(selected.ID)
	_ = d.Set("name", selected.Name)

	return nil
}
