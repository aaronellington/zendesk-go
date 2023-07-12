package zendesk

import (
	"fmt"
	"net/http"
)

type authentication interface {
	AddZendeskAuthentication(r *http.Request)
}

type AuthenticationPassword struct {
	Email    string
	Password string
}

func (auth AuthenticationPassword) AddZendeskAuthentication(r *http.Request) {
	r.SetBasicAuth(
		auth.Email,
		auth.Password,
	)
}

type AuthenticationToken struct {
	Email string
	Token string
}

func (auth AuthenticationToken) AddZendeskAuthentication(r *http.Request) {
	r.SetBasicAuth(
		fmt.Sprintf("%s/token", auth.Email),
		auth.Token,
	)
}

type ChatCredentials struct {
	ClientID     string
	ClientSecret string
}
