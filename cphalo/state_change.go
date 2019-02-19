package cphalo

import (
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/kiwicom/cphalo-go"
)

const (
	stateChangeWaiting = "waiting"
	stateChangeChanged = "changed"
)

func baseStateChange(actionTimeout string, d *schema.ResourceData, f resource.StateRefreshFunc) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{stateChangeWaiting},
		Target:     []string{stateChangeChanged},
		MinTimeout: time.Second,
		Timeout:    d.Timeout(actionTimeout),
		Refresh:    f,
	}

	_, err := stateConf.WaitForState()

	return err
}

func deleteStateChange(d *schema.ResourceData, f resource.StateRefreshFunc) error {
	return baseStateChange(schema.TimeoutDelete, d, f)
}

func createStateChange(d *schema.ResourceData, f resource.StateRefreshFunc) error {
	return baseStateChange(schema.TimeoutCreate, d, f)
}

func updateStateChange(d *schema.ResourceData, f resource.StateRefreshFunc) error {
	return baseStateChange(schema.TimeoutUpdate, d, f)
}

// createStateChangeDefault makes a default checker function to evaluate API create calls
// it waits until 404 error responses are gone - this policy is inline with official API
// recommendations
func createStateChangeDefault(d *schema.ResourceData, f func() (interface{}, error)) error {
	checkerFunc := func() (result interface{}, state string, err error) {
		resp, err := f()

		if err == nil {
			return resp, stateChangeChanged, nil
		}

		if _, ok := err.(*cphalo.ResponseError404); ok {
			return resp, stateChangeWaiting, nil
		}

		return resp, "", err
	}

	return createStateChange(d, checkerFunc)
}

// deleteStateChangeDefault makes a default checker function to evaluate API delete calls
// it waits until a 404 error is received - this policy is inline with official API
// recommendations
func deleteStateChangeDefault(d *schema.ResourceData, f func() (interface{}, error)) error {
	checkerFunc := func() (result interface{}, state string, err error) {
		resp, err := f()

		if err == nil {
			return resp, stateChangeWaiting, nil
		}

		if _, ok := err.(*cphalo.ResponseError404); ok {
			return resp, stateChangeChanged, nil
		}

		return resp, "", err
	}

	return deleteStateChange(d, checkerFunc)
}
