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
	allowedFirewallPolicyPlatforms    = []string{"linux", "windows"}
	allowedFirewallRuleChains         = []string{"INPUT", "OUTPUT"}
	allowedFirewallRuleActions        = []string{"ACCEPT", "DROP", "REJECT"}
	allowedFirewallRuleConnStates     = []string{"NEW", "RELATED", "ESTABLISHED"}
	allowedFirewallRuleSourcesTargets = []string{"User", "UserGroup", "Group", "FirewallZone"}
)

const (
	firewallRuleSourceTargetKindUser         = "User"
	firewallRuleSourceTargetKindGroup        = "Group"
	firewallRuleSourceTargetKindUserGroup    = "UserGroup"
	firewallRuleSourceTargetKindFirewallZone = "FirewallZone"
)

func resourceCPHaloFirewallPolicy() *schema.Resource {
	newSourceTargetResource := func() *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Required: true,
				},
				"kind": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(allowedFirewallRuleSourcesTargets, false),
				},
			},
		}
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "linux",
				ValidateFunc: validation.StringInSlice(allowedFirewallPolicyPlatforms, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true, // TODO: this value must be sent as string and is returned as boolean
				ForceNew: true,
			},
			"ignore_forwarding_rules": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"rule": {
				Description: "Firewall rule",
				Type:        schema.TypeSet,
				Optional:    true,
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
							Description: "Firewall source",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        newSourceTargetResource(),
							Set:         resourceCPHaloFirewallPolicyRuleSourceTargetHash,
						},
						"firewall_target": {
							Description: "Firewall target",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        newSourceTargetResource(),
							Set:         resourceCPHaloFirewallPolicyRuleSourceTargetHash,
						},
					},
				},
				Set: resourceCPHaloFirewallPolicyRuleHash,
			},
		},
		Create:        resourceFirewallPolicyCreate,
		Read:          resourceFirewallPolicyRead,
		Update:        resourceFirewallPolicyUpdate,
		Delete:        resourceFirewallPolicyDelete,
		CustomizeDiff: firewallRuleChecker,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func firewallRuleChecker(d *schema.ResourceDiff, _ interface{}) error {
	var (
		inputRules  = make(map[int]bool)
		outputRules = make(map[int]bool)
		inputCount  int
		outputCount int
	)

	rules := d.Get("rule").(*schema.Set)
	for _, rawData := range rules.List() {
		data := rawData.(map[string]interface{})

		if v, ok := data["firewall_source"]; ok {
			s := v.(*schema.Set)
			if s.Len() > 1 {
				return fmt.Errorf("only one unique firewall_source is allowed per firewall rule")
			}
		}

		if v, ok := data["firewall_target"]; ok {
			s := v.(*schema.Set)
			if s.Len() > 1 {
				return fmt.Errorf("only one unique firewall_target is allowed per firewall rule")
			}
		}

		chain := data["chain"].(string)
		position := data["position"].(int)

		if chain == "OUTPUT" {
			outputRules[position] = true
			outputCount++
		} else {
			inputRules[position] = true
			inputCount++
		}
	}

	if len(inputRules) == 0 && len(outputRules) == 0 {
		return fmt.Errorf("firewall policy %s does not contain any rules", d.Get("name").(string))
	}

	if len(inputRules) != inputCount {
		return fmt.Errorf("INPUT rules have duplicate positions")
	}

	if len(outputRules) != outputCount {
		return fmt.Errorf("OUTPUT rules have duplicate positions")
	}

	for i := 1; i <= inputCount; i++ {
		if _, ok := inputRules[i]; !ok {
			return fmt.Errorf("INPUT rule positions are not in order")
		}
	}

	for i := 1; i <= outputCount; i++ {
		if _, ok := outputRules[i]; !ok {
			return fmt.Errorf("OUTPUT rule positions are not in order")
		}
	}

	return nil
}

func resourceCPHaloFirewallPolicyRuleSourceTargetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	strElements := []string{
		"id",
		"kind",
	}

	for _, element := range strElements {
		if v, ok := m[element]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}
	}

	return hashcode.String(buf.String())
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
	}

	boolElements := []string{
		"active",
	}

	intElements := []string{
		"position",
	}

	sourceTargetElements := []string{
		"firewall_source",
		"firewall_target",
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

	for _, element := range sourceTargetElements {
		if v, ok := m[element]; ok {
			s := v.(*schema.Set)
			if s.Len() > 0 {
				sourceTarget := s.List()[0].(map[string]interface{})

				sourceTargetID := sourceTarget["id"].(string)
				sourceTargetType := sourceTarget["kind"].(string)

				buf.WriteString(fmt.Sprintf("%s-%s-", sourceTargetID, sourceTargetType))
			}
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
	inputRules, outputRules := parseFirewallPolicyRuleSet(d.Get("rule"))

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

	apiRules, err := client.ListFirewallRules(d.Id())
	if err != nil {
		return fmt.Errorf("cannot read server policy (%s) rules: %v", d.Id(), err)
	}

	rules := &schema.Set{
		F: resourceCPHaloFirewallPolicyRuleHash,
	}

	for _, r := range apiRules.Rules {
		details, err := client.GetFirewallRule(d.Id(), r.ID)
		if err != nil {
			return fmt.Errorf("cannot read server policy (%s) rule details: %v", r.ID, err)
		}

		var (
			fwInterfaceID string
			fwServiceID   string
			fwSource      = &schema.Set{F: resourceCPHaloFirewallPolicyRuleSourceTargetHash}
			fwTarget      = &schema.Set{F: resourceCPHaloFirewallPolicyRuleSourceTargetHash}
		)

		if details.Rule.FirewallInterface != nil {
			fwInterfaceID = details.Rule.FirewallInterface.ID
		}

		if details.Rule.FirewallService != nil {
			fwServiceID = details.Rule.FirewallService.ID
		}

		applySourceTarget := func(output *schema.Set, input *api.FirewallRuleSourceTarget) {
			// normally we take the ID
			id := input.ID

			// special cases (all servers, all users, ...) require to use `name` property instead of ID
			if input.ID == "" {
				id = input.Name
			}

			// fix mismatch between CloudPassage requests/responses
			if strings.ToLower(id) == "all active servers" {
				id = "All Active Servers"
			} else if strings.ToLower(id) == "all ghostports users" {
				id = "All GhostPorts users"
			}

			output.Add(map[string]interface{}{
				"id":   id,
				"kind": input.Kind,
			})
		}

		if details.Rule.FirewallSource != nil {
			applySourceTarget(fwSource, details.Rule.FirewallSource)
		}

		if details.Rule.FirewallTarget != nil {
			applySourceTarget(fwTarget, details.Rule.FirewallTarget)
		}

		rules.Add(map[string]interface{}{
			"chain":              details.Rule.Chain,
			"action":             details.Rule.Action,
			"active":             details.Rule.Active,
			"connection_states":  details.Rule.ConnectionStates,
			"position":           details.Rule.Position,
			"firewall_interface": fwInterfaceID,
			"firewall_service":   fwServiceID,
			"firewall_source":    fwSource,
			"firewall_target":    fwTarget,
		})
	}

	if err := d.Set("rule", rules); err != nil {
		return fmt.Errorf("cannot set rules for firewall policy (%s): %v", d.Id(), err)
	}

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
		inputRules, outputRules := parseFirewallPolicyRuleSet(d.Get("rule"))

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

func parseFirewallPolicyRuleSet(rules interface{}) (map[int]api.FirewallRule, map[int]api.FirewallRule) {
	if rules == nil {
		rules = new(schema.Set)
	}

	var (
		inputRules  = make(map[int]api.FirewallRule)
		outputRules = make(map[int]api.FirewallRule)
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
			s := v.(*schema.Set)
			if s.Len() == 1 {
				source := s.List()[0].(map[string]interface{})

				rule.FirewallSource = &api.FirewallRuleSourceTarget{
					ID:   source["id"].(string),
					Kind: source["kind"].(string),
				}
			}
		}

		if v, ok := data["firewall_target"]; ok {
			s := v.(*schema.Set)
			if s.Len() == 1 {
				target := s.List()[0].(map[string]interface{})

				rule.FirewallTarget = &api.FirewallRuleSourceTarget{
					ID:   target["id"].(string),
					Kind: target["kind"].(string),
				}
			}
		}

		if rule.Chain == "OUTPUT" {
			outputRules[rule.Position] = rule
		} else {
			inputRules[rule.Position] = rule
		}
	}

	return inputRules, outputRules
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
