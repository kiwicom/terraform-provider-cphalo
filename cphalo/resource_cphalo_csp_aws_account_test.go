package cphalo

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.com/kiwicom/cphalo-go"
)

func TestAccCSPAWSAccount_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersWithAWS,
		CheckDestroy: testAccCSPAWSAccountCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCPHaloCSPAWSAccountConfig(t, 0),
				Check: resource.ComposeTestCheckFunc(func(_ *terraform.State) error {
					if isCI() {
						time.Sleep(time.Second * 5)
					}
					return nil
				}),
			},
			{
				Config: testAccCPHaloCSPAWSAccountConfig(t, 1),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testCSPAWSAccountAttributes("tf_aws_testacc_basic_01")
				}),
			},
			{
				Config: testAccCPHaloCSPAWSAccountConfig(t, 2),
				Check: resource.ComposeTestCheckFunc(func(s *terraform.State) error {
					return testCSPAWSAccountAttributes("tf_aws_testacc_basic_02")
				}),
			},
		},
	})
}

func testAccCPHaloCSPAWSAccountConfig(t *testing.T, step int) string {
	path := "csp_aws_accounts/basic_00_prerequisites.tf"

	preData, err := readTestTemplateData(path, testID)

	if err != nil {
		t.Fatal(err)
	}

	if step == 0 {
		return preData
	}

	path = fmt.Sprintf("csp_aws_accounts/basic_%.2d.tf", step)

	data, err := readTestTemplateData(path, testID)

	if err != nil {
		t.Fatal(err)
	}

	return preData + data
}

func testCSPAWSAccountAttributes(name string) error {
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListCSPAccounts()

	name = testID + name

	if err != nil {
		return fmt.Errorf("cannot fetch CSP AWS accounts: %v", err)
	}

	var found []string
	for _, c := range resp.CSPAccounts {
		found = append(found, c.AccountDisplayName)
		if c.AccountDisplayName == name {
			return nil
		}
	}

	return fmt.Errorf("expected CSP AWS account %s; found only: %v", name, strings.Join(found, ","))
}

func testAccCSPAWSAccountCheckDestroy(_ *terraform.State) error {
	client := testAccProvider.Meta().(*cphalo.Client)
	resp, err := client.ListCSPAccounts()

	if err != nil {
		return fmt.Errorf("cannot fetch CSP AWS accounts on destroy: %v", err)
	}

	var found []string
	for _, g := range resp.CSPAccounts {
		if strings.HasPrefix(g.AccountDisplayName, testID) {
			found = append(found, g.AWSExternalID)
		}
	}

	if len(found) > 0 {
		return fmt.Errorf("found %d CSP AWS accounts after destroy: %s", resp.Count, strings.Join(found, ","))
	}

	return nil
}
