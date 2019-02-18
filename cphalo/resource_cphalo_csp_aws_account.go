package cphalo

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
)

func resourceCPHaloCSPAWSAccount() *schema.Resource {
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
		Create: resourceCPHaloCSPAWSAccountCreate,
		Read:   resourceCPHaloCSPAWSAccountRead,
		Update: resourceCPHaloCSPAWSAccountUpdate,
		Delete: resourceCPHaloCSPAWSAccountDelete,
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

func resourceCPHaloCSPAWSAccountCreate(d *schema.ResourceData, i interface{}) error {
	account := cphalo.CreateCSPAccountAWSRequest{
		ExternalID:         d.Get("external_id").(string),
		RoleArn:            d.Get("role_arn").(string),
		GroupID:            d.Get("group_id").(string),
		AccountDisplayName: d.Get("account_display_name").(string),
		CSPAccountType:     "aws",
	}

	client := i.(*cphalo.Client)

	resp, err := client.CreateCSPAccount(account)
	if err != nil {
		return fmt.Errorf("cannot create CSP AWS account: %v", err)
	}

	d.SetId(resp.CSPAccount.ID)

	err = createStateChangeDefault(d, func() (interface{}, error) {
		return client.GetCSPAccount(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for CSP AWS account %s to be created: %v", d.Id(), err)
	}

	return resourceCPHaloCSPAWSAccountRead(d, i)
}

func resourceCPHaloCSPAWSAccountRead(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)

	resp, err := client.GetCSPAccount(d.Id())

	if err != nil {
		return fmt.Errorf("cannot read CSP AWS account %s: %v", d.Id(), err)
	}

	account := resp.CSPAccount

	_ = d.Set("external_id", account.AWSExternalID)
	_ = d.Set("role_arn", account.AWSRoleArn)
	_ = d.Set("group_id", account.GroupID)
	_ = d.Set("account_display_name", account.AccountDisplayName)

	return nil
}

func resourceCPHaloCSPAWSAccountUpdate(d *schema.ResourceData, i interface{}) error {
	client := i.(*cphalo.Client)
	_, err := client.GetCSPAccount(d.Id())

	if err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	if d.HasChange("account_display_name") {
		cspAccount := cphalo.CSPAccount{
			ID:                 d.Id(),
			AccountDisplayName: d.Get("account_display_name").(string),
			AWSRoleArn:         d.Get("role_arn").(string),
		}
		if err := client.UpdateCSPAccount(cspAccount); err != nil {
			return fmt.Errorf("updating csp AWS account of %s failed: %v", d.Id(), err)
		}
		logDebug("updated csp AWS account")
	}

	err = updateStateChange(d, func() (result interface{}, state string, err error) {
		resp, err := client.GetCSPAccount(d.Id())

		if err != nil {
			return resp, "", err
		}

		matches := []bool{
			resp.CSPAccount.AccountDisplayName == d.Get("account_display_name").(string),
			resp.CSPAccount.AWSRoleArn == d.Get("role_arn").(string),
		}

		for _, match := range matches {
			if !match {
				return resp, StateChangeWaiting, err
			}
		}

		return resp, StateChangeChanged, nil
	})

	if err != nil {
		return fmt.Errorf("error waiting for CSP AWS account %s to be updated: %v", d.Id(), err)
	}

	return resourceCPHaloCSPAWSAccountRead(d, i)
}

func resourceCPHaloCSPAWSAccountDelete(d *schema.ResourceData, i interface{}) (err error) {
	client := i.(*cphalo.Client)

	if err := client.DeleteCSPAccount(d.Id()); err != nil {
		return fmt.Errorf("failed to delete %s: %v", d.Id(), err)
	}

	err = deleteStateChangeDefault(d, func() (interface{}, error) {
		return client.GetCSPAccount(d.Id())
	})

	if err != nil {
		return fmt.Errorf("error waiting for CSP AWS account %s to be deleted: %v", d.Id(), err)
	}

	logInfof("CSP AWS account %s deleted", d.Id())

	return nil
}
