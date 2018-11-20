package cphalo

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"log"
)

func resourceCPHaloServerGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceCPHaloServerGroupCreate,
		Read:   resourceCPHaloServerGroupRead,
		Update: resourceCPHaloServerGroupUpdate,
		Delete: resourceCPHaloServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCPHaloServerGroupCreate(d *schema.ResourceData, i interface{}) error {
	group := api.ServerGroup{
		Name:     d.Get("name").(string),
		ParentID: d.Get("parent_id").(string),
		Tag:      d.Get("tag").(string),
	}

	client := i.(*api.Client)

	resp, err := client.CreateServerGroup(group)
	if err != nil {
		return fmt.Errorf("cannot create server group: %v", err)
	}

	d.SetId(resp.Group.ID)

	return resourceCPHaloServerGroupRead(d, i)
}

func resourceCPHaloServerGroupRead(d *schema.ResourceData, i interface{}) error {
	client := i.(*api.Client)

	resp, err := client.GetServerGroup(d.Id())

	if err != nil {
		//d.SetId("") // TODO: check if needed
		return fmt.Errorf("cannot read server group %s: %v", d.Id(), err)
	}

	group := resp.Group

	d.Set("name", group.Name)
	d.Set("parent_id", group.ParentID)
	d.Set("tag", group.Tag)

	return nil
}

func resourceCPHaloServerGroupUpdate(d *schema.ResourceData, i interface{}) error {
	client := i.(*api.Client)
	_, err := client.GetServerGroup(d.Id())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	d.Partial(true)
	if d.HasChange("name") {
		if err := client.UpdateServerGroup(api.ServerGroup{ID: d.Id(), Name: d.Get("name").(string)}); err != nil {
			return fmt.Errorf("updating name of %s failed: %v", d.Id(), err)
		}
		log.Println("updated name")
	}

	if d.HasChange("tag") {
		if err := client.UpdateServerGroup(api.ServerGroup{ID: d.Id(), Tag: d.Get("tag").(string)}); err != nil {
			return fmt.Errorf("updating tag of %s failed: %v", d.Id(), err)
		}
		log.Println("updated tag")
	}

	if d.HasChange("parent_id") {
		if err := client.UpdateServerGroup(api.ServerGroup{ID: d.Id(), ParentID: d.Get("parent_id").(string)}); err != nil {
			return fmt.Errorf("updating parent_id of %s failed: %v", d.Id(), err)
		}
		log.Println("updated parent_id")
	}
	d.Partial(false)

	return resourceCPHaloServerGroupRead(d, i)
}

func resourceCPHaloServerGroupDelete(d *schema.ResourceData, i interface{}) error {
	client := i.(*api.Client)

	if err := client.DeleteServerGroup(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	log.Printf("server %s deleted\n", d.Id())

	return nil
}
