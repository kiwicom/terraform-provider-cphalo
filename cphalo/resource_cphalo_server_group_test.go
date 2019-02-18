package cphalo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.com/kiwicom/cphalo-go"
)

func TestAccServerGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCPHaloCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCPHaloServerGroupConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testServerGroupAttributes("root group", "", "")
				}),
			},
			{
				Config: testAccCPHaloServerGroupConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testServerGroupAttributes("changed_name", "added_tag", "and added some interesting description")
				}),
			},
			{
				Config: testAccCPHaloServerGroupConfig(t, 3),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					client := testAccProvider.Meta().(*cphalo.Client)
					resp, err := client.ListServerGroups()

					if err != nil {
						return fmt.Errorf("cannot fetch server groups: %v", err)
					}

					expected := 6

					var servers []string
					for _, g := range resp.Groups {
						if strings.HasPrefix(g.Name, testID) {
							servers = append(servers, g.Name)
						}
					}

					if len(servers) != expected {
						return fmt.Errorf("found %d server groups, expected %d: %s", len(servers), expected, strings.Join(servers, ","))
					}

					return nil
				}),
			},
			{
				Config: testAccCPHaloServerGroupConfig(t, 4),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					client := testAccProvider.Meta().(*cphalo.Client)
					resp, err := client.ListServerGroups()
					nameExpected := testID + "changed_name"

					if err != nil {
						return fmt.Errorf("cannot fetch server groups: %v", err)
					}

					var found cphalo.ServerGroup
					var servers []string
					for _, g := range resp.Groups {
						servers = append(servers, g.Name)
						if g.Name == nameExpected {
							found = g
						}
					}

					if len(found.ID) == 0 {
						return fmt.Errorf("could not find server group %s", nameExpected)
					}

					if len(found.AlertProfileIDs) != 1 {
						return fmt.Errorf("expected 1 alert profile; got %d", len(found.AlertProfileIDs))
					}

					return nil
				}),
			},
		},
	})
}

func testServerGroupAttributes(nameExpected, tagExpected, descriptionExpected string) error {
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListServerGroups()

	if err != nil {
		return fmt.Errorf("cannot fetch server groups: %v", err)
	}

	nameExpected = testID + nameExpected

	if tagExpected != "" {
		tagExpected = testID + tagExpected
	}

	var found cphalo.ServerGroup
	var servers []string
	for _, g := range resp.Groups {
		servers = append(servers, g.Name)
		if g.Name == nameExpected {
			found = g
		}
	}

	if found.Name == "" {
		return fmt.Errorf("could not found server group %s; found only: %s", nameExpected, strings.Join(servers, ","))
	}

	if found.Name != nameExpected {
		return fmt.Errorf("expected server group %s; found only: %s", nameExpected, strings.Join(servers, ","))
	}

	if found.Tag != tagExpected {
		return fmt.Errorf("expected server group %s to have tag '%s'; got: '%s'", nameExpected, tagExpected, found.Tag)
	}

	if found.Description != descriptionExpected {
		return fmt.Errorf("expected server group %s to have description '%s'; got: '%s'", nameExpected, descriptionExpected, found.Description)
	}

	if len(found.AlertProfileIDs) != 0 {
		return fmt.Errorf("expected 0 alert profiles; found %d", len(found.AlertProfileIDs))
	}

	return nil
}

func testAccCPHaloCheckDestroy(_ *terraform.State) error {
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

func testAccCPHaloServerGroupConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("server_groups/basic_%.2d.tf", step)

	data, err := readTestTemplateData(path, testID)

	if err != nil {
		t.Fatal(err)
	}

	return data
}
