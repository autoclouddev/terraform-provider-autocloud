package autocloud_sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)



func (c *Client) GetRepositories() ([]Repository, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/source_control/repositories", c.HostURL), nil)
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
