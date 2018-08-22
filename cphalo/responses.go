package cphalo

import "fmt"

type CpHaloError struct {
	StatusCode int    `json:"statuscode"`
	StatusDesc string `json:"statusdesc"`
	Message    string `json:"errormessage"`
}

type errorJsonResponse struct {
	Error *CpHaloError `json:"error"`
}

type apiKeyJsonResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type groupJsonResponse struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type listGroupJsonReponse struct {
	Group groupDetailJsonResponse
}

type groupDetailJsonResponse struct {
	Name                    string         `json:"name"`
	Id                      string         `json:"id"`
	Url                     string         `json:"url"`
	Description             string         `json:"description"`
	LinuxFirewallPolicyId   string         `json:"linux_firewall_policy_id"`
	WindowsFirewallPolicyId string         `json:"windows_firewall_policy_id"`
	ServerCounts            map[string]int `json:"server_counts"`
}

type listGroupsJsonResponse struct {
	Groups []groupJsonResponse `json:""`
}

func (r *CpHaloError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}
