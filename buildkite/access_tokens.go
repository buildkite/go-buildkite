package buildkite

import "fmt"

type AccessTokensService struct {
	client *Client
}

type AccessToken struct {
	UUID   *string  `json:"uuid,omitempty" yaml:"uuid,omitempty"`
	Scopes *[]Scope `json:"scopes,omitempty" yaml:"scopes,omitempty"`
}

type Scope string

// GetToken gets the current token which was used to authenticate the request
//
// buildkite API docs: https://buildkite.com/docs/rest-api/access-token
func (ats *AccessTokensService) GetToken() (*AccessToken, *Response, error) {

	var u string

	u = fmt.Sprintf("v2/access-token")

	req, err := ats.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var accessToken *AccessToken
	fmt.Print(req.Body)
	resp, err := ats.client.Do(req, accessToken)
	if err != nil {
		return nil, resp, err
	}

	return accessToken, resp, nil
}
