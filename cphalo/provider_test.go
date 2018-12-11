package cphalo

import (
	"github.com/hashicorp/terraform/terraform"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvidersWithAWS map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"cphalo": testAccProvider,
	}
	testAccProvidersWithAWS = map[string]terraform.ResourceProvider{
		"cphalo": testAccProvider,
		"aws":    aws.Provider(),
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_implementation(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CP_APPLICATION_KEY"); v == "" {
		t.Fatal("CP_APPLICATION_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("CP_APPLICATION_SECRET"); v == "" {
		t.Fatal("CP_APPLICATION_SECRET must be set for acceptance tests")
	}
}
