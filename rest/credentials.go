package rest

import "net/http"

func NewCredentials(apiToken string) Credentials {
	return &credentials{APIToken: apiToken}
}

// Credentials are able to modify a given HTTP Request in order to authenticate properly.
// One use case is e.g. adding the Authorization Header
type Credentials interface {
	// Authenticate modifies a given HTTP Request in order to ensure proper authentication on the server side
	Authenticate(request *http.Request) error
	// Configured tells whether actual values for authentication are available
	Configured() bool
}

type credentials struct {
	APIToken string `json:"api-token,omitempty"`
}

// Authenticate modifies a given HTTP Request in order to ensure proper authentication on the server side
func (credentials *credentials) Authenticate(request *http.Request) error {
	request.Header.Set("Authorization", "Api-Token "+credentials.APIToken)
	return nil
}

// Configured tells whether actual values for authentication are available
func (credentials *credentials) Configured() bool {
	return credentials.APIToken != ""
}
