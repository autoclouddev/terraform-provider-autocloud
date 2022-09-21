package autocloud_sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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
	HTTPClient  *http.Client
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
		HostURL:     "http://localhost:8080/api/v.0.0.1", //properties.ApiHost,
		AppClientID: properties.AppClientID,
		HTTPClient:  &http.Client{Timeout: 10 * time.Second},
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return nil, fmt.Errorf("username and password are empty")
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

//TMP hack to insert the token in the header, this is only used by the graphql client
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

func (c *Client) doRequest(req *http.Request, authToken *string) ([]byte, error) {
	token := c.Token

	if authToken != nil {
		token = *authToken
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300

	if !statusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
