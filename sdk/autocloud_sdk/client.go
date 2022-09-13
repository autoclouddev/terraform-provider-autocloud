package autocloud_sdk

import (
	"net/http"

	p "autocloud_sdk/properties"

	graphql "github.com/hasura/go-graphql-client"
)

// Client -
type Client struct {
	HostURL     string
	Token       string
	Auth        AuthStruct
	AppClientID string
	graphql     *graphql.Client
}

/*
	This is a tmp struct due to our api only supports username and password
	This should be changed later on to use another auth strategy
*/
// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id`
	Username string `json:"username`
	Token    string `json:"token"`
}

var Token = ""

// NewClient -
func NewClient(host, username, password *string) (*Client, error) {

	var properties p.Properties = p.LoadProperties()

	c := Client{
		HostURL:     properties.ApiHost,
		AppClientID: properties.AppClientID,
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token
	Token = ar.Token
	c.graphql = graphQLClient(properties.ApiHost)
	return &c, nil
}

type transport struct {
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {

	req.Header.Add("Authorization", Token)
	return t.underlyingTransport.RoundTrip(req)
}

func graphQLClient(apiHost string) *graphql.Client {
	client := graphql.NewClient(apiHost, &http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport}})
	return client
}
