package cphalo

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"log"
	"math"
	"strings"
	"time"
)

var (
	allowedFirewallPolicyPlatforms = []string{"linux", "windows"}
	allowedFirewallRuleChains      = []string{"INPUT", "OUTPUT"}
	allowedFirewallRuleActions     = []string{"ACCEPT", "DROP", "REJECT"}
	allowedFirewallRuleConnStates  = []string{"NEW", "RELATED", "ESTABLISHED"}
)

func resourceCPHaloFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "linux", // FIXME: acc create test breaks if this is not set
				ValidateFunc: validation.StringInSlice(allowedFirewallPolicyPlatforms, false),
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
			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							ValidateFunc: validation.IntBetween(1, math.MaxInt32),
						},
						"firewall_interface": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"firewall_service": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"firewall_source": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"firewall_target": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
				Set: resourceCPHaloFirewallPolicyRuleHash,
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

func resourceCPHaloFirewallPolicyRuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	strElements := []string{
		"chain",
		"action",
		"connection_states",
		"firewall_interface",
		"firewall_service",
		"firewall_source",
		"firewall_target",
	}

	boolElements := []string{
		"active",
	}

	intElements := []string{
		"position",
	}

	for _, element := range strElements {
		if v, ok := m[element]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	for _, element := range boolElements {
		if v, ok := m[element]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
	}

	for _, element := range intElements {
		if v, ok := m[element]; ok {
			buf.WriteString(fmt.Sprintf("%d-", v.(int)))
		}
	}

	return hashcode.String(buf.String())
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

func resourceFirewallPolicyCreate(d *schema.ResourceData, i interface{}) error {
	policy := api.FirewallPolicy{
		Name:                  d.Get("name").(string),
		Platform:              d.Get("platform").(string),
		Description:           d.Get("description").(string),
		Shared:                d.Get("shared").(bool),
		IgnoreForwardingRules: d.Get("ignore_forwarding_rules").(bool),
	}

	client := i.(*api.Client)

	// parse firewall rules before creating a policy, in case there are some issues
	inputRules, outputRules, err := parseFirewallPolicyRuleSet(d.Get("rule"))
	if err != nil {
		return fmt.Errorf("could not parse firewall rules: %v", err)
	}

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

	// create firewall rules
	if err = applyFirewallPolicyRules(d, client, d.Id(), inputRules, outputRules); err != nil {
		return fmt.Errorf("updating rules of %s failed: %v", d.Id(), err)
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

	if d.HasChange("rule") {
		inputRules, outputRules, err := parseFirewallPolicyRuleSet(d.Get("rule"))
		if err != nil {
			return fmt.Errorf("could not parse firewall rules: %v", err)
		}

		err = applyFirewallPolicyRules(d, client, d.Id(), inputRules, outputRules)
		if err != nil {
			return fmt.Errorf("updating rules of %s failed: %v", d.Id(), err)
		}

		d.SetPartial("rule")
		log.Println("updated rule")
	}

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

func applyFirewallPolicyRules(d *schema.ResourceData, client *api.Client, policyID string, inputRules, outputRules map[int]api.FirewallRule) error {
	existingRules, err := client.ListFirewallRules(policyID)
	if err != nil {
		return fmt.Errorf("could not fetch existing firewall rules for policy %s: %v", policyID, err)
	}

	for _, rule := range existingRules.Rules {
		if err = client.DeleteFirewallRule(policyID, rule.ID); err != nil {
			return fmt.Errorf("failed to delete firewall rule %s: %v", rule.ID, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{StateChangeWaiting},
		Target:     []string{StateChangeChanged},
		MinTimeout: time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Refresh: func() (result interface{}, state string, err error) {
			resp, err := client.ListFirewallRules(policyID)

			if err != nil {
				return resp, "", err
			}

			if len(resp.Rules) > 0 {
				return resp, StateChangeWaiting, nil
			}

			return resp, StateChangeChanged, nil
		},
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error waiting for firewall rules for policy %s to be deleted: %v", d.Id(), err)
	}

	log.Printf("existing rules for policy %s have been deleted\n", d.Id())

	createRule := func(rule api.FirewallRule) error {
		_, err := client.CreateFirewallRule(policyID, rule)
		if err != nil {
			return fmt.Errorf("cannot create firewall rule (%+v): %v", rule, err)
		}

		err = createStateChangeDefault(d, func() (interface{}, error) {
			return client.GetFirewallRule(policyID, rule.ID)
		})

		if err != nil {
			return fmt.Errorf("error waiting for firewall rule %s to be created: %v", d.Id(), err)
		}

		return nil
	}

	for i := 1; i <= len(inputRules); i++ {
		if err := createRule(inputRules[i]); err != nil {
			return err
		}

		log.Println("firewall rules for INPUT chain have been created")
	}

	for i := 1; i <= len(outputRules); i++ {
		if err := createRule(outputRules[i]); err != nil {
			return err
		}

		log.Println("firewall rules for OUTPUT chain have been created")
	}

	log.Println("all firewall rules have been created")

	return nil
}

func parseFirewallPolicyRuleSet(rules interface{}) (map[int]api.FirewallRule, map[int]api.FirewallRule, error) {
	if rules == nil {
		rules = new(schema.Set)
	}

	var (
		inputRules  = make(map[int]api.FirewallRule)
		outputRules = make(map[int]api.FirewallRule)
		inputCount  int
		outputCount int
	)

	for _, rawData := range rules.(*schema.Set).List() {
		data := rawData.(map[string]interface{})

		rule := api.FirewallRule{
			FirewallInterface: nil,
			FirewallService:   nil,
			FirewallSource:    nil,
			FirewallTarget:    nil,
		}

		if v, ok := data["chain"]; ok {
			rule.Chain = v.(string)
		}

		if v, ok := data["action"]; ok {
			rule.Action = v.(string)
		}

		if v, ok := data["connection_states"]; ok {
			rule.ConnectionStates = v.(string)
		}

		if v, ok := data["position"]; ok {
			rule.Position = v.(int)
		}

		if v, ok := data["active"]; ok {
			rule.Active = v.(bool)
		}

		if v, ok := data["firewall_interface"]; ok {
			if id := v.(string); id != "" {
				rule.FirewallInterface = &api.FirewallInterface{
					ID: id,
				}
			}
		}

		if v, ok := data["firewall_service"]; ok {
			if id := v.(string); id != "" {
				rule.FirewallService = &api.FirewallService{
					ID: id,
				}
			}
		}

		if v, ok := data["firewall_source"]; ok {
			if id := v.(string); id != "" {
				rule.FirewallSource = &api.FirewallRuleInlineSourceTarget{
					ID: id,
				}
			}
		}

		if v, ok := data["firewall_target"]; ok {
			if id := v.(string); id != "" {
				rule.FirewallTarget = &api.FirewallRuleInlineSourceTarget{
					ID: v.(string),
				}
			}
		}

		if rule.Chain == "OUTPUT" {
			outputRules[rule.Position] = rule
			outputCount++
		} else {
			inputRules[rule.Position] = rule
			inputCount++
		}
	}

	if len(inputRules) != inputCount {
		return inputRules, outputRules, fmt.Errorf("INPUT rules have duplicate positions")
	}

	if len(outputRules) != outputCount {
		return inputRules, outputRules, fmt.Errorf("OUTPUT rules have duplicate positions")
	}

	for i := 1; i <= inputCount; i++ {
		if _, ok := inputRules[i]; !ok {
			return inputRules, outputRules, fmt.Errorf("INPUT rule positions are not in order")
		}
	}

	for i := 1; i <= outputCount; i++ {
		if _, ok := outputRules[i]; !ok {
			return inputRules, outputRules, fmt.Errorf("OUTPUT rule positions are not in order")
		}
	}

	return inputRules, outputRules, nil
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
