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

func TestAccFirewallInterface_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFirewallInterfaceCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallInterfaceConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallInterfaceAttributes("eth42")
				}),
			},
			{
				Config: testAccFirewallInterfaceConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallInterfaceAttributes("eth9001")
				}),
			},
		},
	})
}

func testFirewallInterfaceAttributes(name string) (err error) {
	var (
		client = testAccProvider.Meta().(*api.Client)
		resp   api.ListFirewallInterfacesResponse
	)

	if resp, err = client.ListFirewallInterfaces(); err != nil {
		return fmt.Errorf("cannot fetch firewall interfaces: %v", err)
	}

	var found api.FirewallInterface
	var interfaces []string
	for _, i := range resp.Interfaces {
		interfaces = append(interfaces, i.Name)
		if i.Name == name {
			found = i
		}
	}

	if found.Name == "" {
		return fmt.Errorf("could not find firewall interface %s; found only: %s", name, strings.Join(interfaces, ","))
	}

	if found.Name != name {
		return fmt.Errorf("expected firewall interface %s; found only: %s", name, strings.Join(interfaces, ","))
	}

	return nil
}

func testAccFirewallInterfaceCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListFirewallInterfaces()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall interfaces on destroy: %v", err)
	}

	// FIXME: magic number...
	if resp.Count != 7 {
		var interfaces []string
		for _, i := range resp.Interfaces {
			interfaces = append(interfaces, i.Name)
		}

		return fmt.Errorf("found %d firewall interfaces after destroy: %s", resp.Count, strings.Join(interfaces, ","))
	}

	return nil
}

func testAccFirewallInterfaceConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("testdata/firewall_interfaces/basic_%.2d.tf", step)
	b, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("cannot read file %s: %v", path, err)
	}

	return string(b)
}
