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

func TestAccServerGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCPHaloCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCPHaloServerGroupConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testServerGroupAttributes("root group", "")
				}),
			},
			{
				Config: testAccCPHaloServerGroupConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					return testServerGroupAttributes("changed_name", "added_tag")
				}),
			},
			{
				Config: testAccCPHaloServerGroupConfig(t, 3),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					client := testAccProvider.Meta().(*api.Client)
					resp, err := client.ListServerGroups()

					if err != nil {
						return fmt.Errorf("cannot fetch server groups: %v", err)
					}

					expected := 7
					if resp.Count != expected {
						return fmt.Errorf("expected excatly %d server group, got %d", expected, resp.Count)
					}

					return nil
				}),
			},
		},
	})
}

func testServerGroupAttributes(nameExpected, tagExpected string) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListServerGroups()

	if err != nil {
		return fmt.Errorf("cannot fetch server groups: %v", err)
	}

	if resp.Count != 2 {
		return fmt.Errorf("expected excatly 2 server group, got %d", resp.Count)
	}

	var found api.ServerGroup
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
		return fmt.Errorf("expected server group %s to have tag %s; got: %s", nameExpected, tagExpected, found.Tag)
	}

	return nil
}

func testAccCPHaloCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListServerGroups()

	if err != nil {
		return fmt.Errorf("cannot fetch server groups on destroy: %v", err)
	}

	if resp.Count != 1 {
		var servers []string
		for _, g := range resp.Groups {
			servers = append(servers, g.Name)
		}

		return fmt.Errorf("found %d server groups after destroy: %s", resp.Count, strings.Join(servers, ","))
	}

	return nil
}

func testAccCPHaloServerGroupConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("testdata/server_groups/basic_%.2d.tf", step)
	b, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("cannot read file %s: %v", path, err)
	}

	return string(b)
}
