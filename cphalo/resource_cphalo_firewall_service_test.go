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

const (
	cpHaloExistingFirewallServiceCount = 27
)

func TestAccFirewallService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFirewallServiceCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallServiceConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					// FIXME: refactor test to accept multiple service checks
					return testFirewallServiceAttributes("tf_acc_custom_ssh", "TCP", "2222")
				}),
			},
			{
				Config: testAccFirewallServiceConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallServiceAttributes("tf_acc_custom_ssh", "TCP", "2223")
				}),
			},
		},
	})
}

func testFirewallServiceAttributes(name, protocol, port string) (err error) {
	var (
		client = testAccProvider.Meta().(*api.Client)
		resp   api.ListFirewallServicesResponse
	)

	if resp, err = client.ListFirewallServices(); err != nil {
		return fmt.Errorf("cannot fetch firewall services: %v", err)
	}

	var found api.FirewallService
	var services []string
	for _, i := range resp.Services {
		services = append(services, i.Name)
		if i.Name == name {
			found = i
		}
	}

	if found.Name == "" {
		return fmt.Errorf("could not find firewall service %s; found only: %s", name, strings.Join(services, ","))
	}

	if found.Name != name {
		return fmt.Errorf("expected firewall service %s; found only: %s", name, strings.Join(services, ","))
	}

	if found.Protocol != protocol {
		return fmt.Errorf("expected firewall service %s to have protocol %s; got: %s", name, protocol, found.Protocol)
	}

	if found.Port != port {
		return fmt.Errorf("expected firewall service %s to have port %s; got: %s", name, port, found.Port)
	}

	return nil
}

func testAccFirewallServiceCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListFirewallServices()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall services on destroy: %v", err)
	}

	if resp.Count != cpHaloExistingFirewallServiceCount {
		var services []string
		for _, i := range resp.Services {
			services = append(services, i.Name)
		}

		return fmt.Errorf("found %d firewall services after destroy: %s", resp.Count, strings.Join(services, ","))
	}

	return nil
}

func testAccFirewallServiceConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("testdata/firewall_services/basic_%.2d.tf", step)
	b, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("cannot read file %s: %v", path, err)
	}

	return string(b)
}
