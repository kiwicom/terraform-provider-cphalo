package cphalo

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"log"
	"math"
	"strings"
	"time"
)

var (
	allowedFirewallRuleChains     = []string{"INPUT", "OUTPUT"}
	allowedFirewallRuleActions    = []string{"ACCEPT", "DROP", "REJECT"}
	allowedFirewallRuleConnStates = []string{"NEW", "RELATED", "ESTABLISHED"}
)

func resourceCPHaloFirewallRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"chain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(allowedFirewallRuleChains, false),
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(allowedFirewallRuleActions, false),
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"connection_states": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateFirewallRuleConnectionStates(),
			},
			"position": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, math.MaxInt64), // too much?
			},
		},
		Create: resourceFirewallRuleCreate,
		Read:   resourceFirewallRuleRead,
		Update: resourceFirewallRuleUpdate,
		Delete: resourceFirewallRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func validateFirewallRuleConnectionStates() schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if !ok {
			return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
		}

		var found bool
		for _, checkedValue := range strings.Split(v, ",") {
			found = false
			checkedValue = strings.TrimSpace(checkedValue)

			for _, allowedValue := range allowedFirewallRuleConnStates {
				if checkedValue == allowedValue {
					found = true
				}
			}

			if !found {
				return nil, []error{fmt.Errorf("invalid value for %s (%s)", k, checkedValue)}
			}
		}

		return nil, nil
	}
}

func resourceFirewallRuleCreate(d *schema.ResourceData, i interface{}) error {
	var (
		parentID = d.Get("parent_id").(string)
		client   = i.(*api.Client)
	)

	rule := api.FirewallRule{
		Chain:            d.Get("chain").(string),
		Action:           d.Get("action").(string),
		Active:           d.Get("active").(bool),
		ConnectionStates: d.Get("connection_states").(string),
		Position:         d.Get("position").(int),
	}

	resp, err := client.CreateFirewallRule(parentID, rule)
	if err != nil {
		return fmt.Errorf("cannot create firewall rule: %v", err)
	}

	d.SetId(resp.Rule.ID)

	err = createStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallRule(parentID, d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall rule %s to be created: %v", d.Id(), err)
	}

	return resourceFirewallRuleRead(d, i)
}

func resourceFirewallRuleRead(d *schema.ResourceData, i interface{}) error {
	var (
		parentID = d.Get("parent_id").(string)
		client   = i.(*api.Client)
	)

	resp, err := client.GetFirewallRule(parentID, d.Id())

	if err != nil {
		return fmt.Errorf("cannot read firewall rule %s: %v", d.Id(), err)
	}

	rule := resp.Rule

	d.Set("chain", rule.Chain)
	d.Set("action", rule.Action)
	d.Set("active", rule.Active)
	d.Set("connection_states", rule.ConnectionStates)
	d.Set("position", rule.Position)

	return nil
}

func resourceFirewallRuleUpdate(d *schema.ResourceData, i interface{}) error {
	var (
		parentID = d.Get("parent_id").(string)
		client   = i.(*api.Client)
	)

	_, err := client.GetFirewallRule(parentID, d.Id())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	d.Partial(true)
	if d.HasChange("chain") {
		err := client.UpdateFirewallRule(parentID, api.FirewallRule{
			ID:    d.Id(),
			Chain: d.Get("chain").(string),
		})

		if err != nil {
			return fmt.Errorf("updating chain of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("chain")
		log.Println("updated chain")
	}

	if d.HasChange("action") {
		err := client.UpdateFirewallRule(parentID, api.FirewallRule{
			ID:     d.Id(),
			Action: d.Get("action").(string),
		})

		if err != nil {
			return fmt.Errorf("updating action of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("action")
		log.Println("updated action")
	}

	if d.HasChange("active") {
		err := client.UpdateFirewallRule(parentID, api.FirewallRule{
			ID:     d.Id(),
			Active: d.Get("active").(bool),
		})

		if err != nil {
			return fmt.Errorf("updating active of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("active")
		log.Println("updated active")
	}

	if d.HasChange("connection_states") {
		err := client.UpdateFirewallRule(parentID, api.FirewallRule{
			ID:               d.Id(),
			ConnectionStates: d.Get("connection_states").(string),
		})

		if err != nil {
			return fmt.Errorf("updating connection states of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("connection_states")
		log.Println("updated connection_states")
	}

	if d.HasChange("position") {
		err := client.UpdateFirewallRule(parentID, api.FirewallRule{
			ID:       d.Id(),
			Position: d.Get("position").(int),
		})

		if err != nil {
			return fmt.Errorf("updating position of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("position")
		log.Println("updated position")
	}

	d.Partial(false)

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetFirewallRule(parentID, d.Id())

		if err != nil {
			return resp, "", err
		}

		matches := []bool{
			resp.Rule.Chain == d.Get("chain").(string),
			resp.Rule.Action == d.Get("action").(string),
			resp.Rule.Active == d.Get("active").(bool),
			resp.Rule.ConnectionStates == d.Get("connection_states").(string),
			resp.Rule.Position == d.Get("position").(int),
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

	return resourceFirewallRuleRead(d, i)
}

func resourceFirewallRuleDelete(d *schema.ResourceData, i interface{}) (err error) {
	var (
		parentID = d.Get("parent_id").(string)
		client   = i.(*api.Client)
	)

	if err := client.DeleteFirewallRule(parentID, d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	err = deleteStateChangeDefault(d, func() (interface{}, error) {
		return client.GetFirewallRule(parentID, d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for firewall rule %s to be deleted: %v", d.Id(), err)
	}

	log.Printf("firewall rule %s deleted\n", d.Id())

	return nil
}
