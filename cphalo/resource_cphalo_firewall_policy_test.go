package cphalo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.com/kiwicom/cphalo-go"
)

type expectedFirewallPolicy struct {
	name           string
	description    string
	shared         cphalo.StringableBool
	ignoreFwdRules bool
	rules          []expectedFirewallRule
}

type expectedFirewallRule struct {
	chain       string
	action      string
	states      string
	position    int
	log         bool
	logPrefix   string
	comment     string
	fwInterface expectedFirewallInterface
	fwService   expectedFirewallService
	fwSource    expectedFirewallRuleSourceTarget
	fwTarget    expectedFirewallRuleSourceTarget
}

type expectedFirewallRuleSourceTarget struct {
	name       string
	ipAddress  string
	kind       string
	dataSource bool
}

type expectedFirewallService struct {
	name       string
	protocol   string
	port       string
	dataSource bool
}

type expectedFirewallInterface struct {
	name       string
	dataSource bool
}

func TestAccFirewallPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFirewallPolicyCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallPolicyConfig(t, "basic", 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_fw_policy",
							description:    "",
							shared:         false,
							ignoreFwdRules: false,
							rules: []expectedFirewallRule{
								{chain: "INPUT", action: "DROP", states: "NEW", position: 1},
								{chain: "INPUT", action: "DROP", states: "ESTABLISHED", position: 2},
								{chain: "INPUT", action: "DROP", states: "NEW, ESTABLISHED", position: 3},
								{chain: "OUTPUT", action: "DROP", states: "NEW", position: 1},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, "basic", 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_fw_policy_changed",
							description:    "awesome",
							shared:         true,
							ignoreFwdRules: false,
							rules: []expectedFirewallRule{
								{chain: "INPUT", action: "DROP", states: "NEW", position: 1},
								{chain: "INPUT", action: "DROP", states: "ESTABLISHED", position: 2},
								{chain: "INPUT", action: "DROP", states: "NEW, ESTABLISHED", position: 3},
								{chain: "OUTPUT", action: "DROP", states: "NEW", position: 1},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, "basic", 3),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_fw_policy_changed",
							description:    "awesome",
							shared:         true,
							ignoreFwdRules: true,
							rules: []expectedFirewallRule{
								{chain: "INPUT", action: "DROP", states: "NEW", position: 1},
								{chain: "OUTPUT", action: "DROP", states: "NEW", position: 1},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, "integration", 1),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_fw_policy",
							description:    "awesome",
							shared:         true,
							ignoreFwdRules: true,
							rules: []expectedFirewallRule{
								{
									chain:    "INPUT",
									action:   "ACCEPT",
									states:   "NEW, ESTABLISHED",
									position: 1,
									fwInterface: expectedFirewallInterface{
										name: "eth42",
									},
									fwService: expectedFirewallService{
										name:     "tf_acc_fw_svc",
										protocol: "TCP",
										port:     "2222",
									},
									fwSource: expectedFirewallRuleSourceTarget{
										name:      "tf_acc_fw_in_zone",
										ipAddress: "1.1.1.1",
										kind:      firewallRuleSourceTargetKindFirewallZone,
									},
								},
								{
									chain:    "OUTPUT",
									action:   "ACCEPT",
									states:   "NEW, ESTABLISHED",
									position: 1,
									fwInterface: expectedFirewallInterface{
										name: "eth42",
									},
									fwService: expectedFirewallService{
										name:     "tf_acc_fw_svc",
										protocol: "TCP",
										port:     "2222",
									},
									fwTarget: expectedFirewallRuleSourceTarget{
										name:       "All active servers",
										kind:       firewallRuleSourceTargetKindGroup,
										dataSource: true,
									},
								},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, "data_source", 1),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_data_source_policy",
							description:    "",
							shared:         true,
							ignoreFwdRules: false,
							rules: []expectedFirewallRule{
								{
									chain:    "INPUT",
									action:   "ACCEPT",
									states:   "NEW, ESTABLISHED",
									position: 1,
									fwInterface: expectedFirewallInterface{
										name:       "eth0",
										dataSource: true,
									},
									fwService: expectedFirewallService{
										name:       "http",
										protocol:   "TCP",
										port:       "80",
										dataSource: true,
									},
									fwSource: expectedFirewallRuleSourceTarget{
										name:       "any",
										ipAddress:  "0.0.0.0/0",
										kind:       firewallRuleSourceTargetKindFirewallZone,
										dataSource: true,
									},
								},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, "any_connection_states", 1),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:           "tf_acc_any_conn_states_fw_policy",
							description:    "",
							shared:         true,
							ignoreFwdRules: false,
							rules: []expectedFirewallRule{
								{
									chain:    "INPUT",
									action:   "ACCEPT",
									states:   "ANY",
									position: 1,
								},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, "logs_and_comments", 1),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:   "tf_acc_fw_logging_policy",
							shared: true,
							rules: []expectedFirewallRule{
								{
									chain:     "INPUT",
									action:    "ACCEPT",
									states:    "ANY",
									position:  1,
									log:       true,
									logPrefix: "tf_acc_test_",
									comment:   "tf_acc",
								},
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallPolicyConfig(t, "logs_and_comments", 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallPolicyAttributes(
						expectedFirewallPolicy{
							name:   "tf_acc_fw_logging_policy",
							shared: true,
							rules: []expectedFirewallRule{
								{
									chain:    "INPUT",
									action:   "ACCEPT",
									states:   "ANY",
									position: 1,
									log:      false,
									comment:  "tf_acc_v2",
								},
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
		client = testAccProvider.Meta().(*cphalo.Client)
		policy cphalo.FirewallPolicy
		resp   cphalo.ListFirewallRulesResponse
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

func testHelperFindFirewallPolicyByName(name string) (policy cphalo.FirewallPolicy, err error) {
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListFirewallPolicies()

	name = testID + name

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

func testHelperCompareFirewallPolicyAttributes(policy cphalo.FirewallPolicy, expectedPolicy expectedFirewallPolicy) error {
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

func testHelperCompareFirewallPolicyRuleAttributes(client *cphalo.Client, rules []cphalo.FirewallRule, policy cphalo.FirewallPolicy, expectedPolicy expectedFirewallPolicy) (err error) {
	var (
		rule         cphalo.FirewallRule
		ruleResp     cphalo.GetFirewallRuleResponse
		fetchedRules = make(map[string]cphalo.FirewallRule, len(rules))
	)

	// TODO: cleanup this mess
	for _, expectedRule := range expectedPolicy.rules {
		var found cphalo.FirewallRule

		for _, r := range rules {
			var ok bool
			if rule, ok = fetchedRules[r.ID]; !ok {
				if ruleResp, err = client.GetFirewallRule(policy.ID, r.ID); err != nil {
					return fmt.Errorf("cannot fetch details for firewall rule (%s): %v", r.ID, err)
				}

				fetchedRules[r.ID] = ruleResp.Rule
				rule = ruleResp.Rule
			}

			matches := []bool{
				rule.Chain == expectedRule.chain,
				rule.Action == expectedRule.action,
				rule.ConnectionStates == expectedRule.states,
				rule.Position == expectedRule.position,
				rule.Log == expectedRule.log,
				rule.LogPrefix == expectedRule.logPrefix,
				rule.Comment == expectedRule.comment,
			}

			if expectedRule.fwInterface.name != "" {
				if rule.FirewallInterface == nil {
					matches = append(matches, false)
				} else {
					var expectedName string
					if expectedRule.fwInterface.dataSource {
						expectedName = expectedRule.fwInterface.name
					} else {
						expectedName = testID + expectedRule.fwInterface.name
					}
					matches = append(matches, rule.FirewallInterface.Name == expectedName)
				}
			}

			if expectedRule.fwService.name != "" {
				if rule.FirewallService == nil {
					matches = append(matches, false)
				} else {
					var expectedName string
					if expectedRule.fwService.dataSource {
						expectedName = expectedRule.fwService.name
					} else {
						expectedName = testID + expectedRule.fwService.name
					}

					matches = append(matches, rule.FirewallService.Name == expectedName)
					matches = append(matches, rule.FirewallService.Protocol == expectedRule.fwService.protocol)
					matches = append(matches, rule.FirewallService.Port == expectedRule.fwService.port)
				}
			}

			if expectedRule.fwSource.name != "" {
				if rule.FirewallSource == nil {
					matches = append(matches, false)
				} else {
					var expectedName string
					if expectedRule.fwSource.dataSource {
						expectedName = expectedRule.fwSource.name
					} else {
						expectedName = testID + expectedRule.fwSource.name
					}

					matches = append(matches, rule.FirewallSource.Name == expectedName)
					matches = append(matches, rule.FirewallSource.IPAddress == expectedRule.fwSource.ipAddress)
					matches = append(matches, rule.FirewallSource.Kind == expectedRule.fwSource.kind)
				}
			}

			if expectedRule.fwTarget.name != "" {
				if rule.FirewallTarget == nil {
					matches = append(matches, false)
				} else {
					var expectedName string
					if expectedRule.fwTarget.dataSource {
						expectedName = expectedRule.fwTarget.name
					} else {
						expectedName = testID + expectedRule.fwTarget.name
					}

					matches = append(matches, rule.FirewallTarget.Name == expectedName)
					matches = append(matches, rule.FirewallTarget.IPAddress == expectedRule.fwTarget.ipAddress)
					matches = append(matches, rule.FirewallTarget.Kind == expectedRule.fwTarget.kind)
				}
			}

			allMatch := true
			for _, match := range matches {
				if !match {
					allMatch = false
				}
			}

			if allMatch {
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
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListFirewallPolicies()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall policies on destroy: %v", err)
	}

	var policies []string
	for _, p := range resp.Policies {
		if strings.HasPrefix(p.Name, testID) {
			policies = append(policies, p.Name)
		}
	}

	if len(policies) > 0 {
		return fmt.Errorf("found %d firewall policies after destroy: %s", resp.Count, strings.Join(policies, ","))
	}

	return nil
}

func testAccFirewallPolicyConfig(t *testing.T, prefix string, step int) string {
	path := fmt.Sprintf("firewall_policies/%s_%.2d.tf", prefix, step)

	data, err := readTestTemplateData(path, testID)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
