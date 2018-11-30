package cphalo

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCPHaloFirewallRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFirewallRuleRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceFirewallRuleRead(d *schema.ResourceData, meta interface{}) error {
	// TODO: future me will check if we need this
	//var (
	//	err          error
	//	client       = meta.(*api.Client)
	//	name         = d.Get("name").(string)
	//	rules        api.ListFirewallRulesResponse
	//	selectedRule api.FirewallRule
	//	parentID     = d.Get("parent_id").(string)
	//)
	//
	//if rules, err = client.ListFirewallRules(parentID); err != nil {
	//	return err
	//}
	//
	//for _, rule := range rules.Rules {
	//	if strings.TrimSpace(rule.Name) == name {
	//		selectedRule = rule
	//		break
	//	}
	//}
	//
	//if selectedRule.Name == "" {
	//	return fmt.Errorf("resource %s does not exists", name)
	//}
	//
	//d.SetId(selectedRule.ID)
	//d.Set("name", selectedRule.Name)

	return nil
}
