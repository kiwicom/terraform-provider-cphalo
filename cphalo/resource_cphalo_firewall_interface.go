package cphalo

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
	"log"
	"time"
)

func resourceCPHaloFirewallInterface() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceFirewallInterfaceCreate,
		Read:   resourceFirewallInterfaceRead,
		Update: resourceFirewallInterfaceUpdate,
		Delete: resourceFirewallInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func resourceFirewallInterfaceCreate(d *schema.ResourceData, i interface{}) error {
	policy := cphalo.FirewallInterface{
		Name: d.Get("name").(string),
	}

	client := i.(*cphalo.Client)

	resp, err := client.CreateFirewallInterface(policy)
	if err != nil {
		return fmt.Errorf("cannot create firewall interface: %v", err)
	}

	d.SetId(resp.Interface.ID)

	err = createStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallInterface(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall interface %s to be created: %v", d.Id(), err)
	}

	return resourceFirewallInterfaceRead(d, i)
}

func resourceFirewallInterfaceRead(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)

	resp, err := client.GetFirewallInterface(d.Id())

	if err != nil {
		return fmt.Errorf("cannot read firewall interface %s: %v", d.Id(), err)
	}

	fwInterface := resp.Interface

	d.Set("name", fwInterface.Name)

	return nil
}

func resourceFirewallInterfaceUpdate(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)
	_, err := client.GetFirewallInterface(d.Id())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	d.Partial(true)

	if d.HasChange("name") {
		err := client.UpdateFirewallInterface(cphalo.FirewallInterface{
			ID:   d.Id(),
			Name: d.Get("name").(string),
		})

		if err != nil {
			return fmt.Errorf("updating name of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("name")
		log.Println("updated name")
	}

	d.Partial(false)

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetFirewallInterface(d.Id())

		if err != nil {
			return resp, "", err
		}

		matches := []bool{
			resp.Interface.Name == d.Get("name").(string),
		}

		for _, match := range matches {
			if !match {
				return resp, StateChangeWaiting, err
			}
		}

		return resp, StateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall interface %s to be updated: %v", d.Id(), err)
	}

	return resourceFirewallInterfaceRead(d, i)
}

func resourceFirewallInterfaceDelete(d *schema.ResourceData, i interface{}) (err error) {
	client := i.(*cphalo.Client)

	if err := client.DeleteFirewallInterface(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	err = deleteStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallInterface(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall interface %s to be deleted: %v", d.Id(), err)
	}

	log.Printf("firewall interface %s deleted\n", d.Id())

	return nil
}
