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

func TestAccCSPAccount_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCSPAccountCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCPHaloCSPAccountConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testCSPAccountAttributes("tf_aws_testacc_basic_01")
				}),
			},
			{
				Config: testAccCPHaloCSPAccountConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testCSPAccountAttributes("tf_aws_testacc_basic_02")
				}),
			},
		},
	})
}

func testAccCPHaloCSPAccountConfig(t *testing.T, step int) string {
	path := fmt.Sprintf("testdata/csp_accounts/basic_%.2d.tf", step)
	b, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatalf("cannot read file %s: %v", path, err)
	}

	return string(b)
}

func testCSPAccountAttributes(name string) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListCSPAccounts()

	if err != nil {
		return fmt.Errorf("cannot fetch CSP accounts: %v", err)
	}

	if resp.Count != 1 {
		return fmt.Errorf("expected excatly 1 CSP account, got %d", resp.Count)
	}

	if resp.CSPAccounts[0].AccountDisplayName != name {
		return fmt.Errorf("expected display name %s; got %s", name, resp.CSPAccounts[0].AccountDisplayName)
	}

	return nil
}

func testAccCSPAccountCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*api.Client)
	resp, err := client.ListCSPAccounts()

	if err != nil {
		return fmt.Errorf("cannot fetch CSP accounts on destroy: %v", err)
	}

	if resp.Count != 0 {
		var found []string
		for _, g := range resp.CSPAccounts {
			found = append(found, g.ExternalID)
		}

		return fmt.Errorf("found %d CSP accounts after destroy: %s", resp.Count, strings.Join(found, ","))
	}

	return nil
}
