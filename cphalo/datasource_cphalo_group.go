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
	var g_id string

	client := meta.(*Client)

	gs, err := client.List()
	if err != nil {
		return err
	}
	n := d.Get("name").(string)
	log.Println(gs)

	g_id = ""
	for i := range gs {
		if strings.TrimSpace(gs[i].Name) == n {
			g_id = gs[i].Id
		}
	}

	if g_id == "" {
		return fmt.Errorf("Resouce does not exists")
	}

	g, err := client.Read(g_id)
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
	bodyString := string(bodyBytes)

	g := &listGroupsJsonResponse{}
	err = json.Unmarshal([]byte(bodyString), &g)

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
	bodyString := string(bodyBytes)
	err = json.Unmarshal([]byte(bodyString), &out)

	return &out, err
}

func (c *Client) GetGroup(name string) (*groupJsonResponse, error) {
	var out groupJsonResponse

	req, err := c.NewRequest("GET", "groups?search[name]="+name, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	err = json.Unmarshal([]byte(bodyString), &out)

	return &out, err

}
