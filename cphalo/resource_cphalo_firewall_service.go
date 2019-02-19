package cphalo

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"gitlab.com/kiwicom/cphalo-go"
)

var (
	allowedFirewallServiceProtocols = []string{"TCP", "UDP", "ICMP"}
)

func resourceCPHaloFirewallService() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(allowedFirewallServiceProtocols, false),
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceFirewallServiceCreate,
		Read:   resourceFirewallServiceRead,
		Update: resourceFirewallServiceUpdate,
		Delete: resourceFirewallServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func resourceFirewallServiceCreate(d *schema.ResourceData, i interface{}) error {
	policy := cphalo.FirewallService{
		Name:     d.Get("name").(string),
		Protocol: d.Get("protocol").(string),
		Port:     d.Get("port").(string),
	}

	client := i.(*cphalo.Client)

	resp, err := client.CreateFirewallService(policy)
	if err != nil {
		return fmt.Errorf("cannot create firewall service: %v", err)
	}

	d.SetId(resp.Service.ID)

	err = createStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallService(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall service %s to be created: %v", d.Id(), err)
	}

	return resourceFirewallServiceRead(d, i)
}

func resourceFirewallServiceRead(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)

	resp, err := client.GetFirewallService(d.Id())

	if err != nil {
		return fmt.Errorf("cannot read firewall service %s: %v", d.Id(), err)
	}

	service := resp.Service

	_ = d.Set("name", service.Name)
	_ = d.Set("protocol", service.Protocol)
	_ = d.Set("port", service.Port)

	return nil
}

func resourceFirewallServiceUpdate(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)
	_, err := client.GetFirewallService(d.Id())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	d.Partial(true)

	if d.HasChange("name") {
		err = client.UpdateFirewallService(cphalo.FirewallService{
			ID:   d.Id(),
			Name: d.Get("name").(string),
		})

		if err != nil {
			return fmt.Errorf("updating name of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("name")
		log.Println("updated name")
	}

	if d.HasChange("protocol") {
		err = client.UpdateFirewallService(cphalo.FirewallService{
			ID:       d.Id(),
			Protocol: d.Get("protocol").(string),
		})

		if err != nil {
			return fmt.Errorf("updating protocol of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("protocol")
		log.Println("updated protocol")
	}

	if d.HasChange("port") {
		err = client.UpdateFirewallService(cphalo.FirewallService{
			ID:   d.Id(),
			Port: d.Get("port").(string),
		})

		if err != nil {
			return fmt.Errorf("updating port of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("port")
		log.Println("updated port")
	}

	d.Partial(false)

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetFirewallService(d.Id())

		if err != nil {
			return resp, "", err
		}

		matches := []bool{
			resp.Service.Name == d.Get("name").(string),
			resp.Service.Protocol == d.Get("protocol").(string),
			resp.Service.Port == d.Get("port").(string),
		}

		for _, match := range matches {
			if !match {
				return resp, stateChangeWaiting, err
			}
		}

		return resp, stateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall service %s to be updated: %v", d.Id(), err)
	}

	return resourceFirewallServiceRead(d, i)
}

func resourceFirewallServiceDelete(d *schema.ResourceData, i interface{}) (err error) {
	client := i.(*cphalo.Client)

	if err = client.DeleteFirewallService(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	err = deleteStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallService(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall service %s to be deleted: %v", d.Id(), err)
	}

	log.Printf("firewall service %s deleted\n", d.Id())

	return nil
}
