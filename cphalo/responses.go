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
}

type listGroupsJsonResponse struct {
	Groups []groupJsonResponse `json:"group"`
}

func (r *CpHaloError) Error() string {
	return fmt.Sprintf("%d %v: %v", r.StatusCode, r.StatusDesc, r.Message)
}
