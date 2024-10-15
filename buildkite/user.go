package buildkite

import (
	"context"
	"fmt"
)

// UserService handles communication with the user related
// methods of the buildkite API.
//
// buildkite API docs: https://buildkite.com/docs/api
type UserService struct {
	client *Client
}

// User represents a buildkite user.
type User struct {
	ID        string     `json:"id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Email     string     `json:"email,omitempty"`
	CreatedAt *Timestamp `json:"created_at,omitempty"`
}

// Get the current user.
//
// buildkite API docs: https://buildkite.com/docs/api
func (us *UserService) Get(ctx context.Context) (User, *Response, error) {
	u := fmt.Sprintf("v2/user")
	req, err := us.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return User{}, nil, err
	}

	var user User
	resp, err := us.client.Do(req, &user)
	if err != nil {
		return User{}, resp, err
	}

	return user, resp, err
}
