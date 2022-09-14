package autocloud_sdk

import (
	r "autocloud_sdk/rest"
	"encoding/json"
	"fmt"
)

type IacCatalog struct {
	id   string `json:"id"`
	name string `json:"name"`
}

func (*Client) CreateGenerator(name string) (*IacCatalog, error) {

	newGenerator := IacCatalog{
		name: "hola",
	}
	json, err := json.Marshal(newGenerator)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp := r.Post("iac_generator", string(json), "API KEY")
	fmt.Printf(resp)

	return &newGenerator, nil
}
