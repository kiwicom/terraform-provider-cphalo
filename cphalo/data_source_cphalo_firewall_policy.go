package cphalo

import (
	"fmt"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFirewallPolicyRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceFirewallPolicyRead(d *schema.ResourceData, meta interface{}) error {
	var (
		err            error
		client         = meta.(*api.Client)
		name           = d.Get("name").(string)
		policies       api.ListFirewallPoliciesResponse
		selectedPolicy api.FirewallPolicy
	)

	if policies, err = client.ListFirewallPolicies(); err != nil {
		return err
	}

	for _, policy := range policies.Policies {
		if strings.TrimSpace(policy.Name) == name {
			selectedPolicy = policy
			break
		}
	}

	if selectedPolicy.Name == "" {
		return fmt.Errorf("resource %s does not exists", name)
	}

	d.SetId(selectedPolicy.ID)
	d.Set("name", selectedPolicy.Name)

	return nil
}
