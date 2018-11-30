package cphalo

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"io/ioutil"
	"strings"
	"testing"
)

type expectedFirewallRule struct {
	chain    string
	action   string
	states   string
	position int
}

func TestAccFirewallRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFirewallRuleCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallRuleConfig(t, "basic", 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallRuleAttributes(
						"awesome firewall policy",
						[]expectedFirewallRule{
							{"INPUT", "DROP", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "ACCEPT", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "DROP", "RELATED", 2},
						},
					)
				}),
			},
			{
				Config: testAccFirewallRuleConfig(t, "basic", 2),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallRuleAttributes(
						"awesome firewall policy",
						[]expectedFirewallRule{
							{"INPUT", "DROP", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "ACCEPT", "ESTABLISHED", 1},
							{"OUTPUT", "DROP", "NEW", 2},
						},
					)
				}),
			},
			{
				Config: testAccFirewallRuleConfig(t, "basic", 3),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallRuleAttributes(
						"awesome firewall policy",
						[]expectedFirewallRule{
							{"INPUT", "DROP", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "DROP", "NEW", 1},
						},
					)
				}),
			},
			{
				Config: testAccFirewallRuleConfig(t, "position", 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallRuleAttributes(
						"awesome firewall policy",
						[]expectedFirewallRule{
							{"INPUT", "DROP", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "ACCEPT", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "DROP", "RELATED", 2},
						},
					)
				}),
			},
			{
				Config: testAccFirewallRuleConfig(t, "position", 2),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallRuleAttributes(
						"awesome firewall policy",
						[]expectedFirewallRule{
							{"INPUT", "DROP", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "ACCEPT", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "DROP", "RELATED", 3},
							{"OUTPUT", "DROP", "NEW", 2},
						},
					)
				}),
			},
			{
				Config: testAccFirewallRuleConfig(t, "position", 3),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallRuleAttributes(
						"awesome firewall policy",
						[]expectedFirewallRule{
							{"INPUT", "DROP", "NEW, ESTABLISHED", 1},
							{"OUTPUT", "DROP", "NEW", 1},
							{"OUTPUT", "ACCEPT", "NEW, ESTABLISHED", 2},
							{"OUTPUT", "DROP", "RELATED", 3},
						},
					)
				}),
			},
		},
	})
}

func testHelperFindFirewallPolicyByName(name string) (policy api.FirewallPolicy, err error) {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListFirewallPolicies()

	if err != nil {
		return policy, fmt.Errorf("cannot fetch firewall policy: %v", err)
	}

	for _, p := range resp.Policies {
		if p.Name == name {
			policy = p
		}
	}

	return policy, nil
}

func testFirewallRuleAttributes(policyName string, expectedRules []expectedFirewallRule) (err error) {
	var (
		client   = testAccProvider.Meta().(*api.Client)
		policy   api.FirewallPolicy
		rule     api.FirewallRule
		resp     api.ListFirewallRulesResponse
		ruleResp api.GetFirewallRuleResponse
	)

	if policy, err = testHelperFindFirewallPolicyByName(policyName); err != nil {
		return err
	}

	if resp, err = client.ListFirewallRules(policy.ID); err != nil {
		return fmt.Errorf("cannot fetch firewall rule: %v", err)
	}

	fetchedRules := make(map[string]api.FirewallRule, len(resp.Rules))

	for _, expectedRule := range expectedRules {
		var found api.FirewallRule

		for _, r := range resp.Rules {
			var ok bool
			if rule, ok = fetchedRules[r.ID]; !ok {
				if ruleResp, err = client.GetFirewallRule(policy.ID, r.ID); err != nil {
					return fmt.Errorf("cannot fetch details for firewall rule (%s): %v", r.ID, err)
				}

				fetchedRules[r.ID] = ruleResp.Rule
				rule = ruleResp.Rule
			}

			sameChain := rule.Chain == expectedRule.chain
			sameAction := rule.Action == expectedRule.action
			sameStates := rule.ConnectionStates == expectedRule.states
			samePosition := rule.Position == expectedRule.position

			if sameChain && sameAction && sameStates && samePosition {
				found = r
			}
		}

		if found.Chain == "" {
			return fmt.Errorf("could not find correct firewall rule on position %d", expectedRule.position)
		}
	}

	return nil
}

func testAccFirewallRuleCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListFirewallPolicies()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall rules on destroy: %v", err)
	}

	if resp.Count != 0 {
		var policies []string
		for _, g := range resp.Policies {
			policies = append(policies, g.Name)
		}

		return fmt.Errorf("found %d firewall rules after destroy: %s", resp.Count, strings.Join(policies, ","))
	}

	return nil
}

func testAccFirewallRuleConfig(t *testing.T, category string, step int) string {
	path := fmt.Sprintf("testdata/firewall_rules/%s_%.2d.tf", category, step)
	b, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("cannot read file %s: %v", path, err)
	}

	return string(b)
}
