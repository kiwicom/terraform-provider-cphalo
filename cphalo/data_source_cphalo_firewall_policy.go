package cphalo

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
)

func dataSourceCPHaloFirewallPolicy() *schema.Resource {
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
		client         = meta.(*cphalo.Client)
		name           = d.Get("name").(string)
		policies       cphalo.ListFirewallPoliciesResponse
		selectedPolicy cphalo.FirewallPolicy
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
	_ = d.Set("name", selectedPolicy.Name)

	return nil
}
