package cphalo

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
)

func resourceCPHaloFirewallZone() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip_address": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Create: resourceFirewallZoneCreate,
		Read:   resourceFirewallZoneRead,
		Update: resourceFirewallZoneUpdate,
		Delete: resourceFirewallZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func resourceFirewallZoneCreate(d *schema.ResourceData, i interface{}) error {
	policy := cphalo.FirewallZone{
		Name:        d.Get("name").(string),
		IPAddress:   flattenIpAddress(d.Get("ip_address").([]interface{})),
		Description: d.Get("description").(string),
	}

	client := i.(*cphalo.Client)

	resp, err := client.CreateFirewallZone(policy)
	if err != nil {
		return fmt.Errorf("cannot create firewall zone: %v", err)
	}

	d.SetId(resp.Zone.ID)

	err = createStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallZone(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall zone %s to be created: %v", d.Id(), err)
	}

	return resourceFirewallZoneRead(d, i)
}

func resourceFirewallZoneRead(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)

	resp, err := client.GetFirewallZone(d.Id())

	if err != nil {
		return fmt.Errorf("cannot read firewall zone %s: %v", d.Id(), err)
	}

	zone := resp.Zone

	_ = d.Set("name", zone.Name)
	_ = d.Set("ip_address", expandIpAddress(zone.IPAddress))
	_ = d.Set("description", zone.Description)

	return nil
}

func resourceFirewallZoneUpdate(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)
	_, err := client.GetFirewallZone(d.Id())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	d.Partial(true)

	if d.HasChange("name") {
		err = client.UpdateFirewallZone(cphalo.FirewallZone{
			ID:   d.Id(),
			Name: d.Get("name").(string),
		})

		if err != nil {
			return fmt.Errorf("updating name of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("name")
		log.Println("updated name")
	}

	if d.HasChange("ip_address") {
		err = client.UpdateFirewallZone(cphalo.FirewallZone{
			ID:        d.Id(),
			IPAddress: flattenIpAddress(d.Get("ip_address").([]interface{})),
		})

		if err != nil {
			return fmt.Errorf("updating ip_address of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("ip_address")
		log.Println("updated ip_address")
	}

	if d.HasChange("description") {
		err = client.UpdateFirewallZone(cphalo.FirewallZone{
			ID:          d.Id(),
			Description: d.Get("description").(string),
		})

		if err != nil {
			return fmt.Errorf("updating description of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("description")
		log.Println("updated description")
	}

	d.Partial(false)

	flatIPs := flattenIpAddress(d.Get("ip_address").([]interface{}))

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetFirewallZone(d.Id())

		if err != nil {
			return resp, "", err
		}

		matches := []bool{
			resp.Zone.Name == d.Get("name").(string),
			resp.Zone.IPAddress == flatIPs,
		}

		for _, match := range matches {
			if !match {
				return resp, stateChangeWaiting, err
			}
		}

		return resp, stateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall zone %s to be updated: %v", d.Id(), err)
	}

	return resourceFirewallZoneRead(d, i)
}

func resourceFirewallZoneDelete(d *schema.ResourceData, i interface{}) (err error) {
	client := i.(*cphalo.Client)

	if err = client.DeleteFirewallZone(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	err = deleteStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallZone(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall zone %s to be deleted: %v", d.Id(), err)
	}

	log.Printf("firewall zone %s deleted\n", d.Id())

	return nil
}

func flattenIpAddress(ips []interface{}) string {
	var result []string
	for _, ip := range ips {
		result = append(result, ip.(string))
	}
	return strings.Join(result, ",")
}

func expandIpAddress(ips string) []string {
	return strings.Split(ips, ",")
}
