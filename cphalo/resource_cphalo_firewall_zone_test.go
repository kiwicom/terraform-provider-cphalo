package cphalo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.com/kiwicom/cphalo-go"
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
					return testFirewallZoneAttributes("tf_acc_fw_zone", cphalo.IPList{"1.1.1.1", "2.2.2.2"}, "")
				}),
			},
			{
				Config: testAccFirewallZoneConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallZoneAttributes("tf_acc_fw_zone", cphalo.IPList{"3.3.3.3", "4.4.4.4"}, "fw zone")
				}),
			},
		},
	})
}

func testFirewallZoneAttributes(name string, ipAddress cphalo.IPList, description string) (err error) {
	var (
		client = testAccProvider.Meta().(*cphalo.Client)
		resp   cphalo.ListFirewallZonesResponse
	)

	if resp, err = client.ListFirewallZones(); err != nil {
		return fmt.Errorf("cannot fetch firewall zones: %v", err)
	}

	name = testID + name

	var found cphalo.FirewallZone
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

	if fmt.Sprintf("%s", found.IPAddress) != fmt.Sprintf("%s", ipAddress) {
		return fmt.Errorf("expected firewall zone %s to have ip_address %s; got: %s", name, ipAddress, found.IPAddress)
	}

	if found.Description != description {
		return fmt.Errorf("expected firewall zone %s to have description %s; got: %s", name, description, found.Description)
	}

	return nil
}

func testAccFirewallZoneCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListFirewallZones()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall services on destroy: %v", err)
	}

	var userCreatedZones []string

	for _, zone := range resp.Zones {
		if strings.HasPrefix(zone.Name, testID) {
			userCreatedZones = append(userCreatedZones, zone.Name)
		}
	}

	if len(userCreatedZones) > 0 {
		zones := strings.Join(userCreatedZones, ",")
		return fmt.Errorf("found %d user-created firewall zones after destroy: %s", resp.Count, zones)
	}

	return nil
}

func testAccFirewallZoneConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("firewall_zones/basic_%.2d.tf", step)

	data, err := readTestTemplateData(path, testID)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
