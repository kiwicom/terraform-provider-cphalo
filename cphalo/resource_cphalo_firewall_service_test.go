package cphalo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.com/kiwicom/cphalo-go"
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
					return testFirewallServiceAttributes(
						[]expectedFirewallService{
							{
								name:     "tf_acc_custom_ssh",
								protocol: "TCP",
								port:     "2222",
							},
							{
								name:     "tf_acc_ping",
								protocol: "ICMP",
							},
						},
					)
				}),
			},
			{
				Config: testAccFirewallServiceConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testFirewallServiceAttributes(
						[]expectedFirewallService{
							{
								name:     "tf_acc_custom_ssh",
								protocol: "TCP",
								port:     "2223",
							},
						},
					)
				}),
			},
		},
	})
}

func testFirewallServiceAttributes(expectedServices []expectedFirewallService) (err error) {
	var (
		svc             cphalo.FirewallService
		svcResp         cphalo.GetFirewallServiceResponse
		client          = testAccProvider.Meta().(*cphalo.Client)
		resp            cphalo.ListFirewallServicesResponse
		fetchedServices = make(map[string]cphalo.FirewallService)
	)

	if resp, err = client.ListFirewallServices(); err != nil {
		return fmt.Errorf("cannot fetch firewall services: %v", err)
	}

	for _, expectedService := range expectedServices {
		var found cphalo.FirewallService

		for _, s := range resp.Services {
			if s.System {
				continue
			}

			var ok bool
			if svc, ok = fetchedServices[s.ID]; !ok {
				if svcResp, err = client.GetFirewallService(s.ID); err != nil {
					return fmt.Errorf("cannot fetch details for firewall service (%s): %v", s.ID, err)
				}

				fetchedServices[s.ID] = svcResp.Service
				svc = svcResp.Service
			}

			matches := []bool{
				svc.Name == testID+expectedService.name,
				svc.Protocol == expectedService.protocol,
				svc.Port == expectedService.port,
			}

			allMatch := true
			for _, match := range matches {
				if !match {
					allMatch = false
				}
			}

			if allMatch {
				found = s
			}
		}

		if found.Name == "" {
			return fmt.Errorf("could not find correct firewall service on position %s", expectedService.name)
		}
	}

	return nil
}

func testAccFirewallServiceCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListFirewallServices()

	if err != nil {
		return fmt.Errorf("cannot fetch firewall services on destroy: %v", err)
	}

	var userCreatedServices []string

	for _, svc := range resp.Services {
		if strings.HasPrefix(svc.Name, testID) {
			userCreatedServices = append(userCreatedServices, svc.Name)
		}
	}

	if len(userCreatedServices) > 0 {
		services := strings.Join(userCreatedServices, ",")
		return fmt.Errorf("found %d user-created firewall services after destroy: %s", resp.Count, services)
	}

	return nil
}

func testAccFirewallServiceConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("firewall_services/basic_%.2d.tf", step)

	data, err := readTestTemplateData(path, testID)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
