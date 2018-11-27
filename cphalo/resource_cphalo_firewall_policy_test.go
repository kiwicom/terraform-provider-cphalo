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
						"uncomplicated firewall",
						"",
						true,
						false,
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						"complicated firewall",
						"awesome",
						true,
						false,
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, 3),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						"complicated firewall",
						"awesome",
						true,
						true,
					)
				}),
			},
		},
	})
}

func testFirewallPolicyAttributes(name, description string, shared, ignoreFwdRules bool) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListFirewallPolicies()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall policy: %v", err)
	}

	if resp.Count != 1 {
		return fmt.Errorf("expected exactly 1 firewall policy, got %d", resp.Count)
	}

	var found api.FirewallPolicy
	var policies []string
	for _, g := range resp.Policies {
		policies = append(policies, g.Name)
		if g.Name == name {
			found = g
		}
	}

	if found.Name == "" {
		return fmt.Errorf("could not find firewall policy %s; found only: %s", name, strings.Join(policies, ","))
	}

	if found.Name != name {
		return fmt.Errorf("expected firewall policy %s; found only: %s", name, strings.Join(policies, ","))
	}

	if found.Description != description {
		return fmt.Errorf("expected firewall policy %s to have description %s; got: %s", name, description, found.Description)
	}

	if found.Shared != shared {
		return fmt.Errorf("expected firewall policy %s to have shared %t; got: %t", name, shared, found.Shared)
	}

	if found.IgnoreForwardingRules != ignoreFwdRules {
		return fmt.Errorf("expected firewall policy %s to have ignore_forwarding_rules %t; got: %t", name, ignoreFwdRules, found.IgnoreForwardingRules)
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
