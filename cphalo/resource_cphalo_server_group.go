package cphalo

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"log"
	"time"
)

func resourceCPHaloServerGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"linux_firewall_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_profile_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Create: resourceCPHaloServerGroupCreate,
		Read:   resourceCPHaloServerGroupRead,
		Update: resourceCPHaloServerGroupUpdate,
		Delete: resourceCPHaloServerGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func resourceCPHaloServerGroupCreate(d *schema.ResourceData, i interface{}) error {
	policyID := d.Get("linux_firewall_policy_id").(string)

	group := api.ServerGroup{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		ParentID:              d.Get("parent_id").(string),
		Tag:                   d.Get("tag").(string),
		LinuxFirewallPolicyID: api.NullableString(policyID),
		AlertProfileIDs:       expandStringList(d.Get("alert_profile_ids")),
	}

	client := i.(*api.Client)

	resp, err := client.CreateServerGroup(group)
	if err != nil {
		return fmt.Errorf("cannot create server group: %v", err)
	}

	d.SetId(resp.Group.ID)

	err = createStateChangeDefault(d, func() (interface{}, error) {
		return client.GetServerGroup(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for server group %s to be created: %v", d.Id(), err)
	}

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
	d.Set("description", group.Description)
	d.Set("parent_id", group.ParentID)
	d.Set("tag", group.Tag)
	d.Set("linux_firewall_policy_id", group.LinuxFirewallPolicyID)
	d.Set("alert_profile_ids", group.AlertProfileIDs)

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
		d.SetPartial("name")
		log.Println("updated name")
	}

	if d.HasChange("description") {
		if err := client.UpdateServerGroup(api.ServerGroup{ID: d.Id(), Description: d.Get("description").(string)}); err != nil {
			return fmt.Errorf("updating description of %s failed: %v", d.Id(), err)
		}
		d.SetPartial("description")
		log.Println("updated description")
	}

	if d.HasChange("tag") {
		if err := client.UpdateServerGroup(api.ServerGroup{ID: d.Id(), Tag: d.Get("tag").(string)}); err != nil {
			return fmt.Errorf("updating tag of %s failed: %v", d.Id(), err)
		}
		d.SetPartial("tag")
		log.Println("updated tag")
	}

	if d.HasChange("parent_id") {
		if err := client.UpdateServerGroup(api.ServerGroup{ID: d.Id(), ParentID: d.Get("parent_id").(string)}); err != nil {
			return fmt.Errorf("updating parent_id of %s failed: %v", d.Id(), err)
		}
		d.SetPartial("parent_id")
		log.Println("updated parent_id")
	}

	if d.HasChange("linux_firewall_policy_id") {
		policyID := d.Get("linux_firewall_policy_id").(string)

		log.Println("POLICY ID: ", policyID)
		if err := client.UpdateServerGroup(api.ServerGroup{ID: d.Id(), LinuxFirewallPolicyID: api.NullableString(policyID)}); err != nil {
			return fmt.Errorf("updating linux_firewall_policy_id of %s failed: %v", d.Id(), err)
		}
		d.SetPartial("linux_firewall_policy_id")
		log.Println("updated linux_firewall_policy_id")
	}

	if d.HasChange("alert_profile_ids") {
		ids := expandStringList(d.Get("alert_profile_ids"))

		log.Println("Profile IDs: ", ids)
		if err := client.UpdateServerGroup(api.ServerGroup{ID: d.Id(), AlertProfileIDs: ids}); err != nil {
			return fmt.Errorf("updating alert_profile_ids of %s failed: %v", d.Id(), err)
		}
		d.SetPartial("alert_profile_ids")
		log.Println("updated alert_profile_ids")
	}
	d.Partial(false)

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetServerGroup(d.Id())

		if err != nil {
			return resp, "", err
		}

		matches := []bool{
			resp.Group.Name == d.Get("name").(string),
			resp.Group.Tag == d.Get("tag").(string),
			resp.Group.ParentID == d.Get("parent_id").(string),
			resp.Group.LinuxFirewallPolicyID == api.NullableString(d.Get("linux_firewall_policy_id").(string)),
			assertStringSlice(resp.Group.AlertProfileIDs, expandStringList(d.Get("alert_profile_ids"))),
		}

		for _, match := range matches {
			if !match {
				return resp, StateChangeWaiting, err
			}
		}

		return resp, StateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for server group %s to be updated: %v", d.Id(), err)
	}

	return resourceCPHaloServerGroupRead(d, i)
}

func resourceCPHaloServerGroupDelete(d *schema.ResourceData, i interface{}) (err error) {
	client := i.(*api.Client)

	if err := client.DeleteServerGroup(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	err = deleteStateChangeDefault(d, func() (interface{}, error) {
		return client.GetServerGroup(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for server group %s to be deleted: %v", d.Id(), err)
	}

	log.Printf("server %s deleted\n", d.Id())

	return nil
}