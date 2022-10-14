package autocloud_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (c *Client) GetModules() ([]IacModule, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/iac_terraform_modules", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	modules := []IacModule{}
	err = json.Unmarshal(body, &modules)
	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (c *Client) GetModule(moduleId string) (*IacModule, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/iac_terraform_modules/%s", c.HostURL, moduleId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	module := IacModule{}
	err = json.Unmarshal(body, &module)
	if err != nil {
		return nil, err
	}

	return &module, nil
}

func (c *Client) CreateModule(module *IacModule) (*IacModule, error) {
	log.Printf("CreateModule IacModule: %+v\n\n", module)

	reqBody := GetIacModuleInput(module)
	log.Printf("CreateModule IacModuleInput: %+v\n\n", reqBody)

	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	fmt.Println(strings.NewReader(string(rb)))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/iac_terraform_modules", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	newModule := IacModule{}
	err = json.Unmarshal(body, &newModule)
	if err != nil {
		return nil, err
	}

	log.Printf("Create Module response: %+v\n", newModule)

	return &newModule, nil
}

func (c *Client) DeleteModule(moduleId string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/iac_terraform_modules/%s", c.HostURL, moduleId), nil)
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

func (c *Client) UpdateModule(module *IacModule) (*IacModule, error) {
	reqBody := GetIacModuleInput(module)
	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/iac_terraform_modules/%s", c.HostURL, module.ID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	updatedGen := IacModule{}
	err = json.Unmarshal(body, &updatedGen)
	if err != nil {
		return nil, err
	}

	return &updatedGen, nil
}

// IAC Module CRUD from IAC Catalog
func (c *Client) CreateCatalogModule(catalog IacCatalog) (*IacModule, error) {

	iacModule, err := GetIacModule(catalog)
	if err != nil {
		return nil, err
	}

	return c.CreateModule(iacModule)
}

func (c *Client) UpdateCatalogModule(catalog IacCatalog) (*IacModule, error) {

	iacModule, err := GetIacModule(catalog)
	if err != nil {
		log.Fatal("Error getting IacCatalogInput")
		return nil, err
	}
	iacModule, err = c.UpdateModule(iacModule)
	if err != nil {
		log.Fatal("Error getting IacCatalogInput")
		return nil, err
	}
	return iacModule, nil
}
