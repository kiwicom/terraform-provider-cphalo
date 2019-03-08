package cphalo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.com/kiwicom/cphalo-go"
)

func TestAccServerGroupFirewallPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCPHaloCheckFirewallPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCPHaloServerGroupFirewallPolicyConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testServerGroupFirewallPolicyAttributes("root group", "tf_acc_fw_policy")
				}),
			},
			{
				Config: testAccCPHaloServerGroupFirewallPolicyConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					var (
						client    = testAccProvider.Meta().(*cphalo.Client)
						groupName = testID + "root group"
					)

					group, err := findServerGroup(client, groupName)
					if err != nil {
						return err
					}

					groupPolicy, err := client.GetServerGroupFirewallPolicy(group.ID)
					if err != nil {
						return err
					}

					if groupPolicy.Group.LinuxFirewallPolicyID != cphalo.NullableString("") {
						return fmt.Errorf("expected server group %s to not have a linux firewall policy; found: %s", groupName, groupPolicy.Group.LinuxFirewallPolicyID)
					}

					return nil
				}),
			},
		},
	})
}

func findServerGroup(client *cphalo.Client, name string) (cphalo.ServerGroup, error) {
	var result cphalo.ServerGroup

	resp, err := client.ListServerGroups()
	if err != nil {
		return result, fmt.Errorf("could not list server groups: %v", err)
	}

	var servers []string
	for _, g := range resp.Groups {
		servers = append(servers, g.Name)
		if g.Name == name {
			result = g
		}
	}

	if result.Name == "" {
		return result, fmt.Errorf("could not find server group %s; found only: %s", name, strings.Join(servers, ","))
	}

	return result, nil
}

func findFirewallPolicy(client *cphalo.Client, name string) (cphalo.FirewallPolicy, error) {
	var result cphalo.FirewallPolicy

	resp, err := client.ListFirewallPolicies()
	if err != nil {
		return result, fmt.Errorf("could not list firewall policies: %v", err)
	}

	var policies []string
	for _, p := range resp.Policies {
		policies = append(policies, p.Name)
		if p.Name == name {
			result = p
		}
	}

	if result.Name == "" {
		return result, fmt.Errorf("could not find firewall policy %s; found only: %s", name, strings.Join(policies, ","))
	}

	return result, nil
}

func testServerGroupFirewallPolicyAttributes(groupName, policyName string) error {
	client := testAccProvider.Meta().(*cphalo.Client)

	serverGroup, err := findServerGroup(client, testID+groupName)
	if err != nil {
		return err
	}

	firewallPolicy, err := findFirewallPolicy(client, testID+policyName)
	if err != nil {
		return err
	}

	groupPolicy, err := client.GetServerGroupFirewallPolicy(serverGroup.ID)
	if err != nil {
		return fmt.Errorf("could not get server group firewall policies: %v", err)
	}

	if groupPolicy.Group.LinuxFirewallPolicyID != cphalo.NullableString(firewallPolicy.ID) {
		return fmt.Errorf("expected server group %s to have linux_firwall_policy_id '%s'; got: '%s'", groupName, firewallPolicy.ID, groupPolicy.Group.LinuxFirewallPolicyID)
	}

	return nil
}

func testAccCPHaloCheckFirewallPolicyDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListServerGroups()

	if err != nil {
		return fmt.Errorf("cannot fetch server groups on destroy: %v", err)
	}

	var servers []string
	for _, g := range resp.Groups {
		if strings.HasPrefix(g.Name, testID) {
			servers = append(servers, g.Name)
		}
	}

	if len(servers) > 0 {
		return fmt.Errorf("found %d server groups after destroy: %s", resp.Count, strings.Join(servers, ","))
	}

	return nil
}

func testAccCPHaloServerGroupFirewallPolicyConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("server_group_firewall_policies/basic_%.2d.tf", step)

	data, err := readTestTemplateData(path, testID)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
