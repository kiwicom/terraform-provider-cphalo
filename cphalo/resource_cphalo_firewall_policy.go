package cphalo

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"log"
	"time"
)

func resourceCPHaloFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "linux", // FIXME: acc create test breaks if this is not set
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true, // FIXME: is this ok?
				ForceNew: true,
			},
			"ignore_forwarding_rules": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
		Create: resourceFirewallPolicyCreate,
		Read:   resourceFirewallPolicyRead,
		Update: resourceFirewallPolicyUpdate,
		Delete: resourceFirewallPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func resourceFirewallPolicyCreate(d *schema.ResourceData, i interface{}) error {
	policy := api.FirewallPolicy{
		Name:                  d.Get("name").(string),
		Platform:              d.Get("platform").(string),
		Description:           d.Get("description").(string),
		Shared:                d.Get("shared").(bool),
		IgnoreForwardingRules: d.Get("ignore_forwarding_rules").(bool),
	}

	client := i.(*api.Client)

	resp, err := client.CreateFirewallPolicy(policy)
	if err != nil {
		return fmt.Errorf("cannot create firewall policy: %v", err)
	}

	d.SetId(resp.Policy.ID)

	err = createStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallPolicy(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall policy %s to be created: %v", d.Id(), err)
	}

	return resourceFirewallPolicyRead(d, i)
}

func resourceFirewallPolicyRead(d *schema.ResourceData, i interface{}) error {
	client := i.(*api.Client)

	resp, err := client.GetFirewallPolicy(d.Id())

	if err != nil {
		return fmt.Errorf("cannot read server policy %s: %v", d.Id(), err)
	}

	policy := resp.Policy

	d.Set("name", policy.Name)
	d.Set("platform", policy.Platform)
	d.Set("description", policy.Description)
	d.Set("shared", policy.Shared)
	d.Set("ignore_forwarding_rules", policy.IgnoreForwardingRules)

	return nil
}

func resourceFirewallPolicyUpdate(d *schema.ResourceData, i interface{}) error {
	client := i.(*api.Client)
	_, err := client.GetFirewallPolicy(d.Id())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	d.Partial(true)
	if d.HasChange("name") {
		err := client.UpdateFirewallPolicy(api.FirewallPolicy{
			ID:   d.Id(),
			Name: d.Get("name").(string),
		})

		if err != nil {
			return fmt.Errorf("updating name of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("name")
		log.Println("updated name")
	}

	if d.HasChange("platform") {
		err := client.UpdateFirewallPolicy(api.FirewallPolicy{
			ID:       d.Id(),
			Platform: d.Get("platform").(string),
		})

		if err != nil {
			return fmt.Errorf("updating platform of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("platform")
		log.Println("updated platform")
	}

	if d.HasChange("description") {
		err := client.UpdateFirewallPolicy(api.FirewallPolicy{
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

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetFirewallPolicy(d.Id())

		if err != nil {
			return resp, "", err
		}

		matches := []bool{
			resp.Policy.Name == d.Get("name").(string),
			resp.Policy.Description == d.Get("description").(string),
			resp.Policy.Platform == d.Get("platform").(string),
		}

		for _, match := range matches {
			if !match {
				return resp, StateChangeWaiting, err
			}
		}

		return resp, StateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall rule %s to be updated: %v", d.Id(), err)
	}

	return resourceFirewallPolicyRead(d, i)
}

func resourceFirewallPolicyDelete(d *schema.ResourceData, i interface{}) (err error) {
	client := i.(*api.Client)

	if err := client.DeleteFirewallPolicy(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	err = deleteStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallPolicy(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall policy %s to be deleted: %v", d.Id(), err)
	}

	log.Printf("firewall policy %s deleted\n", d.Id())

	return nil
}
