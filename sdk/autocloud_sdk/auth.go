package autocloud_sdk

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// SignIn - Get a new token for user
func (c *Client) SignIn() (*AuthResponse, error) {

	conf := &aws.Config{Region: aws.String("us-east-1")}
	sess, err := session.NewSession(conf)
	CognitoClient := cognito.New(sess)
	if err != nil {
		panic(err)
	}
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return nil, fmt.Errorf("define username and password")
	}

	params := map[string]*string{
		"USERNAME": aws.String(c.Auth.Username),
		"PASSWORD": aws.String(c.Auth.Password),
	}

	authTry := &cognito.InitiateAuthInput{
		AuthFlow:       aws.String("CUSTOM_AUTH"),
		AuthParameters: params,
		ClientId:       aws.String(c.AppClientID),
	}

	req, resp := CognitoClient.InitiateAuth(authTry)
	if resp != nil {
		fmt.Printf("error loggin in %s", resp)
	}

	ar := AuthResponse{
		Token: *req.AuthenticationResult.IdToken,
	}
	return &ar, nil
}
