package autocloud_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) GetGenerators() ([]IacCatalog, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/iac_generators", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	generators := []IacCatalog{}
	err = json.Unmarshal(body, &generators)
	if err != nil {
		return nil, err
	}

	return generators, nil
}

func (c *Client) GetGenerator(generatorID string) (*IacCatalog, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/iac_generators/%s", c.HostURL, generatorID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	generator := IacCatalog{}
	err = json.Unmarshal(body, &generator)
	if err != nil {
		return nil, err
	}

	return &generator, nil
}

func (c *Client) CreateGenerator(generator IacCatalog) (*IacCatalog, error) {
	fmt.Println(generator)
	reqBody := struct {
		Name string `json:"name"`
	}{
		Name: generator.Name,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	fmt.Println(strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/iac_generators", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	newGenerator := IacCatalog{}
	err = json.Unmarshal(body, &newGenerator)
	if err != nil {
		return nil, err
	}

	return &newGenerator, nil
}

func (c *Client) DeleteGenerator(generatorID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/iac_generators/%s", c.HostURL, generatorID), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return err
	}

	if string(body) != "OK" {
		return errors.New(string(body))
	}

	return nil
}

func (c *Client) UpdateGenerator(generator IacCatalog) (*IacCatalog, error) {
	reqBody := struct {
		Name string `json:"name"`
	}{
		Name: generator.Name,
	}
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/iac_generators/%s", c.HostURL, generator.ID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	updatedGen := IacCatalog{}
	err = json.Unmarshal(body, &updatedGen)
	if err != nil {
		return nil, err
	}

	return &updatedGen, nil
}
