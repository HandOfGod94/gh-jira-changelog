package api_token

import "net/http"

type apiTokeAuthenticator struct {
	token string
}

func NewApiAuthenticator(token string) *apiTokeAuthenticator {
	return &apiTokeAuthenticator{token: token}
}

func (a *apiTokeAuthenticator) Login() error {
	// No Op
	return nil
}

func (a *apiTokeAuthenticator) Client() *http.Client {
	return nil
}
