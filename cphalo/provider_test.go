package cphalo

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"cphalo": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigure(t *testing.T) {
	var expectedApplicationKey string
	var expectedApplicationSecret string

	if v := os.Getenv("CP_APPLICATION_KEY"); v != "" {
		expectedApplicationKey = v
	} else {
		expectedApplicationKey = "foo"
	}

	if v := os.Getenv("CP_APPLICATION_SECRET"); v != "" {
		expectedApplicationSecret = v
	} else {
		expectedApplicationSecret = "foo"
	}

	raw := map[string]interface{}{
		"application_key":    expectedApplicationKey,
		"application_secret": expectedApplicationSecret,
	}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	rp := Provider().(*schema.Provider)
	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	config := rp.Meta().(*Client)
	if config.AppKey != expectedApplicationKey {
		t.Fatalf("bad: %#v", config)
	}

	if config.AppSecret != expectedApplicationSecret {
		t.Fatalf("bad: %#v", config)
	}
}
