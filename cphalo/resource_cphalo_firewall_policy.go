package cphalo

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
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

	//if d.HasChange("shared") {
	//	if err := client.UpdateFirewallPolicy(api.FirewallPolicy{ID: d.Id(), Shared: d.Get("shared").(bool)}); err != nil {
	//		return fmt.Errorf("updating shared of %s failed: %v", d.Id(), err)
	//	}
	//	d.SetPartial("shared")
	//	log.Println("updated shared")
	//}

	//if d.HasChange("ignore_forwarding_rules") {
	//	if err := client.UpdateFirewallPolicy(api.FirewallPolicy{ID: d.Id(), IgnoreForwardingRules: d.Get("ignore_forwarding_rules").(bool)}); err != nil {
	//		return fmt.Errorf("updating ignore_forwarding_rules of %s failed: %v", d.Id(), err)
	//	}
	//	d.SetPartial("ignore_forwarding_rules")
	//	log.Println("updated ignore_forwarding_rules")
	//}

	d.Partial(false)

	return resourceFirewallPolicyRead(d, i)
}

func resourceFirewallPolicyDelete(d *schema.ResourceData, i interface{}) (err error) {
	client := i.(*api.Client)

	if err := client.DeleteFirewallPolicy(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"waiting"},
		Target:     []string{"deleted"},
		MinTimeout: time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Refresh: func() (result interface{}, state string, err error) {
			resp, err := client.GetFirewallPolicy(d.Id())

			if err == nil {
				return resp, "waiting", nil
			}

			if _, ok := err.(*api.ResponseError404); ok {
				return resp, "deleted", nil
			}

			return resp, "", err
		},
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for firewall policy %s to be deleted: %v", d.Id(), err)
	}

	log.Printf("firewall policy %s deleted\n", d.Id())

	return nil
}
