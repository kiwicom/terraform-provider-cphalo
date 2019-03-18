package cphalo

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
)

func resourceCPHaloServerGroupFirewallPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"linux_firewall_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Create: resourceCPHaloServerGroupFirewallPolicyUpdate,
		Read:   resourceCPHaloServerGroupFirewallPolicyRead,
		Update: resourceCPHaloServerGroupFirewallPolicyUpdate,
		Delete: resourceCPHaloServerGroupFirewallPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func resourceCPHaloServerGroupFirewallPolicyUpdate(d *schema.ResourceData, i interface{}) error {
	groupID := d.Get("group_id").(string)
	client := i.(*cphalo.Client)
	_, err := client.GetServerGroupFirewallPolicy(groupID)

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	d.Partial(true)
	if d.HasChange("linux_firewall_policy_id") {

		err = client.UpdateServerGroupFirewallPolicy(cphalo.ServerGroupFirewallPolicy{
			GroupID:               groupID,
			LinuxFirewallPolicyID: cphalo.NullableString(d.Get("linux_firewall_policy_id").(string)),
		})
		if err != nil {
			return fmt.Errorf("updating linux_firewall_policy_id of %s failed: %v", groupID, err)
		}
		d.SetPartial("linux_firewall_policy_id")
		logDebug("updated linux_firewall_policy_id")
	}
	d.Partial(false)

	d.SetId(groupID)

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetServerGroupFirewallPolicy(d.Id())

		if err != nil {
			return resp, "", err
		}

		if resp.Group.LinuxFirewallPolicyID != cphalo.NullableString(d.Get("linux_firewall_policy_id").(string)) {
			return resp, stateChangeWaiting, err
		}

		return resp, stateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for server group firewall policy %s to be updated: %v", d.Id(), err)
	}

	return resourceCPHaloServerGroupFirewallPolicyRead(d, i)
}

func resourceCPHaloServerGroupFirewallPolicyRead(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)

	resp, err := client.GetServerGroupFirewallPolicy(d.Id())

	if err != nil {
		return fmt.Errorf("cannot read server group %s: %v", d.Id(), err)
	}

	group := resp.Group

	_ = d.Set("group_id", group.GroupID)
	_ = d.Set("linux_firewall_policy_id", group.LinuxFirewallPolicyID)

	return nil
}

func resourceCPHaloServerGroupFirewallPolicyDelete(d *schema.ResourceData, i interface{}) error {
	group := cphalo.ServerGroupFirewallPolicy{
		GroupID:               d.Get("group_id").(string),
		LinuxFirewallPolicyID: cphalo.NullableString(""),
	}

	client := i.(*cphalo.Client)

	err := client.UpdateServerGroupFirewallPolicy(group)
	if err != nil {
		return fmt.Errorf("cannot update server group firewall policy: %v", err)
	}

	d.SetId(group.GroupID)

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetServerGroupFirewallPolicy(d.Id())

		if err != nil {
			return resp, "", err
		}

		if resp.Group.LinuxFirewallPolicyID != cphalo.NullableString("") {
			return resp, stateChangeWaiting, err
		}

		return resp, stateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for server group firewall policy %s to be deleted: %v", d.Id(), err)
	}

	logInfof("server group firewall policy %s deleted", d.Id())

	return nil
}
