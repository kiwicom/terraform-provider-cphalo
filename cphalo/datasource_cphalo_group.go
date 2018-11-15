package cphalo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCpHaloGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCpHaloGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceCpHaloGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client)

	gs, err := client.List()
	if err != nil {
		return err
	}
	n := d.Get("name").(string)
	log.Println(gs)

	gId := ""
	for i := range gs {
		if strings.TrimSpace(gs[i].Name) == n {
			gId = gs[i].Id
		}
	}

	if gId == "" {
		return fmt.Errorf("Resouce does not exists")
	}

	g, err := client.Read(gId)
	if err != nil {
		return err
	}

	log.Println(g.Group)

	d.SetId(g.Group.Id)
	d.Set("name", g.Group.Name)

	return nil
}
