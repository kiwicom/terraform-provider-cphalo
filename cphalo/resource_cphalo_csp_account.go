package cphalo

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.skypicker.com/terraform-provider-cphalo/api"
	"time"
)

func resourceCPHaloCSPAccount() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"external_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"account_display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Create: resourceCPHaloCSPAccountCreate,
		Read:   resourceCPHaloCSPAccountRead,
		Update: resourceCPHaloCSPAccountUpdate,
		Delete: resourceCPHaloCSPAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 2),
			Update: schema.DefaultTimeout(time.Minute * 2),
			Delete: schema.DefaultTimeout(time.Minute * 5),
		},
	}
}

func resourceCPHaloCSPAccountCreate(d *schema.ResourceData, i interface{}) error {
	account := api.CreateCSPAccountRequest{
		ExternalID:         d.Get("external_id").(string),
		RoleArn:            d.Get("role_arn").(string),
		GroupID:            d.Get("group_id").(string),
		AccountDisplayName: d.Get("account_display_name").(string),
	}

	client := i.(*api.Client)

	resp, err := client.CreateCSPAccount(account)
	if err != nil {
		return fmt.Errorf("cannot create CSP account: %v", err)
	}

	d.SetId(resp.CSPAccount.ID)

	err = createStateChangeDefault(d, func() (interface{}, error) {
		return client.GetCSPAccount(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for CSP account %s to be created: %v", d.Id(), err)
	}

	return resourceCPHaloCSPAccountRead(d, i)
}

func resourceCPHaloCSPAccountRead(d *schema.ResourceData, i interface{}) error {
	client := i.(*api.Client)

	resp, err := client.GetCSPAccount(d.Id())

	if err != nil {
		return fmt.Errorf("cannot read CSP account %s: %v", d.Id(), err)
	}

	account := resp.CSPAccount

	d.Set("external_id", account.ExternalID)
	d.Set("role_arn", account.RoleArn)
	d.Set("group_id", account.GroupID)
	d.Set("account_display_name", account.AccountDisplayName)

	return nil
}

func resourceCPHaloCSPAccountUpdate(d *schema.ResourceData, i interface{}) error {
	client := i.(*api.Client)
	_, err := client.GetCSPAccount(d.Id())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	d.Partial(true)
	if d.HasChange("account_display_name") {
		if err := client.UpdateCSPAccount(api.CSPAccount{ID: d.Id(), AccountDisplayName: d.Get("account_display_name").(string)}); err != nil {
			return fmt.Errorf("updating account_display_name of %s failed: %v", d.Id(), err)
		}
		d.SetPartial("account_display_name")
		logDebug("updated account_display_name")
	}
	d.Partial(false)

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetCSPAccount(d.Id())

		if err != nil {
			return resp, "", err
		}

		matches := []bool{
			resp.CSPAccount.AccountDisplayName == d.Get("account_display_name").(string),
		}

		for _, match := range matches {
			if !match {
				return resp, StateChangeWaiting, err
			}
		}

		return resp, StateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for CSP account %s to be updated: %v", d.Id(), err)
	}

	return resourceCPHaloCSPAccountRead(d, i)
}

func resourceCPHaloCSPAccountDelete(d *schema.ResourceData, i interface{}) (err error) {
	client := i.(*api.Client)

	if err := client.DeleteCSPAccount(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	err = deleteStateChangeDefault(d, func() (interface{}, error) {
		return client.GetCSPAccount(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for CSP account %s to be deleted: %v", d.Id(), err)
	}

	logInfof("CSP account %s deleted", d.Id())

	return nil
}
