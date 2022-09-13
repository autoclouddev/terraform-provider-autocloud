package main

import (
	"fmt"

	"autocloud_sdk"
)

func main() {
	username := ""
	password := ""
	client, err := autocloud_sdk.NewClient(nil, &username, &password)
	if err != nil {
		fmt.Println("error intializing client")
	}
	fmt.Println(client.HostURL)
	user, err := client.GetMe()
	if err != nil {
		fmt.Println("error calling api", err)
	}
	fmt.Println(user)
}
