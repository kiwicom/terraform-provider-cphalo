package cphalo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCpHaloGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCpHaloGroupRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceCpHaloGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client)

	gs, err := client.List()
	if err != nil {
		return err
	}
	n := d.Get("name").(string)
	log.Println(gs)

	gId := ""
	for i := range gs {
		if strings.TrimSpace(gs[i].Name) == n {
			gId = gs[i].Id
		}
	}

	if gId == "" {
		return fmt.Errorf("Resouce does not exists")
	}

	g, err := client.Read(gId)
	if err != nil {
		return err
	}

	log.Println(g.Group)

	d.SetId(g.Group.Id)
	d.Set("name", g.Group.Name)

	return nil
}

func (cs *Client) List() ([]groupJsonResponse, error) {
	req, err := cs.NewRequest("GET", "groups", nil)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	g := &listGroupsJsonResponse{}
	err = json.Unmarshal(bodyBytes, g)

	return g.Groups, err
}

func (cs *Client) Read(id string) (*listGroupJsonReponse, error) {
	var out listGroupJsonReponse

	req, err := cs.NewRequest("GET", "groups/"+id, nil)
	if err != nil {
		return nil, err
	}

	resp, err := cs.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(bodyBytes, &out)

	return &out, err
}
