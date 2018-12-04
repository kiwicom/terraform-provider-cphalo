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

type expectedFirewallPolicy struct {
	name           string
	description    string
	shared         bool
	ignoreFwdRules bool
	rules          []expectedFirewallRule
}

type expectedFirewallRule struct {
	chain    string
	action   string
	states   string
	position int
}

func TestAccFirewallPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFirewallPolicyCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallPolicyConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_fw_policy",
							description:    "",
							shared:         true,
							ignoreFwdRules: false,
							rules: []expectedFirewallRule{
								{"INPUT", "DROP", "NEW", 1},
								{"INPUT", "DROP", "ESTABLISHED", 2},
								{"INPUT", "DROP", "NEW, ESTABLISHED", 3},
								{"OUTPUT", "DROP", "NEW", 1},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_fw_policy_changed",
							description:    "awesome",
							shared:         true,
							ignoreFwdRules: false,
							rules: []expectedFirewallRule{
								{"INPUT", "DROP", "NEW", 1},
								{"INPUT", "DROP", "ESTABLISHED", 2},
								{"INPUT", "DROP", "NEW, ESTABLISHED", 3},
								{"OUTPUT", "DROP", "NEW", 1},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, 3),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_fw_policy_changed",
							description:    "awesome",
							shared:         true,
							ignoreFwdRules: true,
							rules: []expectedFirewallRule{
								{"INPUT", "DROP", "NEW", 1},
								{"OUTPUT", "DROP", "NEW", 1},
							},
						},
					)
				}),
			},
		},
	})
}

func testFirewallPolicyAttributes(expectedPolicy expectedFirewallPolicy) (err error) {
	var (
		client = testAccProvider.Meta().(*api.Client)
		policy api.FirewallPolicy
		resp   api.ListFirewallRulesResponse
	)

	if policy, err = testHelperFindFirewallPolicyByName(expectedPolicy.name); err != nil {
		return err
	}

	if resp, err = client.ListFirewallRules(policy.ID); err != nil {
		return fmt.Errorf("cannot fetch firewall rule: %v", err)
	}

	if err := testHelperCompareFirewallPolicyAttributes(policy, expectedPolicy); err != nil {
		return err
	}

	if err := testHelperCompareFirewallPolicyRuleAttributes(client, resp.Rules, policy, expectedPolicy); err != nil {
		return err
	}

	return nil
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

func testHelperCompareFirewallPolicyAttributes(policy api.FirewallPolicy, expectedPolicy expectedFirewallPolicy) error {
	if policy.Name == "" {
		return fmt.Errorf("could not find firewall policy %s", expectedPolicy.name)
	}

	if policy.Description != expectedPolicy.description {
		return fmt.Errorf("expected firewall policy %s to have description %s; got: %s", expectedPolicy.name, expectedPolicy.description, policy.Description)
	}

	if policy.Shared != expectedPolicy.shared {
		return fmt.Errorf("expected firewall policy %s to have shared %t; got: %t", expectedPolicy.name, expectedPolicy.shared, policy.Shared)
	}

	if policy.IgnoreForwardingRules != expectedPolicy.ignoreFwdRules {
		return fmt.Errorf("expected firewall policy %s to have ignore_forwarding_rules %t; got: %t", expectedPolicy.name, expectedPolicy.ignoreFwdRules, policy.IgnoreForwardingRules)
	}

	return nil
}

func testHelperCompareFirewallPolicyRuleAttributes(client *api.Client, rules []api.FirewallRule, policy api.FirewallPolicy, expectedPolicy expectedFirewallPolicy) (err error) {
	var (
		rule         api.FirewallRule
		ruleResp     api.GetFirewallRuleResponse
		fetchedRules = make(map[string]api.FirewallRule, len(rules))
	)

	for _, expectedRule := range expectedPolicy.rules {
		var found api.FirewallRule

		for _, r := range rules {
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

func testAccFirewallPolicyCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListFirewallPolicies()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall policies on destroy: %v", err)
	}

	if resp.Count != 0 {
		var policies []string
		for _, g := range resp.Policies {
			policies = append(policies, g.Name)
		}

		return fmt.Errorf("found %d firewall policies after destroy: %s", resp.Count, strings.Join(policies, ","))
	}

	return nil
}

func testAccFirewallPolicyConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("testdata/firewall_policies/basic_%.2d.tf", step)
	b, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("cannot read file %s: %v", path, err)
	}

	return string(b)
}
