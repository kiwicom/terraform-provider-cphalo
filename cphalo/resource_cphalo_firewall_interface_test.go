package cphalo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.com/kiwicom/cphalo-go"
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
		client = testAccProvider.Meta().(*cphalo.Client)
		resp   cphalo.ListFirewallInterfacesResponse
	)

	if resp, err = client.ListFirewallInterfaces(); err != nil {
		return fmt.Errorf("cannot fetch firewall interfaces: %v", err)
	}

	name = testID + name

	var found cphalo.FirewallInterface
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
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListFirewallInterfaces()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall interfaces on destroy: %v", err)
	}

	var userCreatedInterfaces []string

	for _, ifc := range resp.Interfaces {
		if strings.HasPrefix(ifc.Name, testID) {
			userCreatedInterfaces = append(userCreatedInterfaces, ifc.Name)
		}
	}

	if len(userCreatedInterfaces) > 0 {
		interfaces := strings.Join(userCreatedInterfaces, ",")
		return fmt.Errorf("found %d user-created firewall interfaces after destroy: %s", resp.Count, interfaces)
	}

	return nil
}

func testAccFirewallInterfaceConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("firewall_interfaces/basic_%.2d.tf", step)

	data, err := readTestTemplateData(path, testID)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
