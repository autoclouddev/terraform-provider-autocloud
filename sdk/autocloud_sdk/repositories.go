package autocloud_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Url         string `json:htmlUrl`
	Description string `json:description`
}

func (c *Client) GetRepositories(organizationName string) ([]Repository, error) {
	endpoint := fmt.Sprintf(organizationName + "/repositories")
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.HostURL, endpoint), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, nil)
	if err != nil {
		return nil, err
	}

	repositories := []Repository{}
	err = json.Unmarshal(body, &repositories)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}
