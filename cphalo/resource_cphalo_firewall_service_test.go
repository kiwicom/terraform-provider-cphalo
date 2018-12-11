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
		svc             api.FirewallService
		svcResp         api.GetFirewallServiceResponse
		client          = testAccProvider.Meta().(*api.Client)
		resp            api.ListFirewallServicesResponse
		fetchedServices = make(map[string]api.FirewallService)
	)

	if resp, err = client.ListFirewallServices(); err != nil {
		return fmt.Errorf("cannot fetch firewall services: %v", err)
	}

	for _, expectedService := range expectedServices {
		var found api.FirewallService

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
				svc.Name == expectedService.name,
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
