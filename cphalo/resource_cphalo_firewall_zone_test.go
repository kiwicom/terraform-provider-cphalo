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
	cpHaloExistingFirewallZoneCount = 1
)

func TestAccFirewallZone_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccFirewallZoneCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallZoneConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testFirewallZoneAttributes("tf_acc_fw_zone", "1.1.1.1")
				}),
			},
			{
				Config: testAccFirewallZoneConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallZoneAttributes("tf_acc_fw_zone", "2.2.2.2")
				}),
			},
		},
	})
}

func testFirewallZoneAttributes(name, ipAddress string) (err error) {
	var (
		client = testAccProvider.Meta().(*api.Client)
		resp   api.ListFirewallZonesResponse
	)

	if resp, err = client.ListFirewallZones(); err != nil {
		return fmt.Errorf("cannot fetch firewall zones: %v", err)
	}

	var found api.FirewallZone
	var zones []string
	for _, i := range resp.Zones {
		zones = append(zones, i.Name)
		if i.Name == name {
			found = i
		}
	}

	if found.Name == "" {
		return fmt.Errorf("could not find firewall zone %s; found only: %s", name, strings.Join(zones, ","))
	}

	if found.Name != name {
		return fmt.Errorf("expected firewall zone %s; found only: %s", name, strings.Join(zones, ","))
	}

	if found.IpAddress != ipAddress {
		return fmt.Errorf("expected firewall zone %s to have ip_address %s; got: %s", name, ipAddress, found.IpAddress)
	}

	return nil
}

func testAccFirewallZoneCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListFirewallZones()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall services on destroy: %v", err)
	}

	if resp.Count != cpHaloExistingFirewallZoneCount {
		var zones []string
		for _, i := range resp.Zones {
			zones = append(zones, i.Name)
		}

		return fmt.Errorf("found %d firewall zones after destroy: %s", resp.Count, strings.Join(zones, ","))
	}

	return nil
}

func testAccFirewallZoneConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("testdata/firewall_zones/basic_%.2d.tf", step)
	b, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("cannot read file %s: %v", path, err)
	}

	return string(b)
}
