package blueprint_config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/hashicorp/go-tfe"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_references"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

type Credentials struct {
	Credentials map[string]Credential `json:"credentials"`
}

type Credential struct {
	Token string `json:"token"`
}

func getTFCCredential() (string, error) {
	// Open the credentials file
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	file, err := os.Open(path.Join(userHomeDir, ".terraform.d", "credentials.tfrc.json"))
	if err != nil {
		return "", fmt.Errorf("error opening credentials file: %s", err)
	}
	defer file.Close()

	// Parse the JSON from the file into a Credentials struct
	var creds Credentials
	err = json.NewDecoder(file).Decode(&creds)
	if err != nil {
		return "", fmt.Errorf("error parsing credentials file: %s", err)
	}

	// Get the token from the "app.terraform.io" credential
	token := creds.Credentials["app.terraform.io"].Token
	return token, nil
}

func readState(fileData interface{}) (map[string]string, error) {
	state := fileData.(map[string]interface{})
	resources := state["resources"].([]interface{})
	references := make(map[string]string)

	for _, v := range resources {
		data := v.(map[string]interface{})
		if data["type"] == "autocloud_blueprint_config" {
			instances := data["instances"].([]interface{})
			for _, v := range instances {
				rawData := v.(map[string]interface{})
				attributes := rawData["attributes"].(map[string]interface{})

				if attributes["references"] != nil {
					rawAliases := attributes["references"].(string)
					storedAliases := make(map[string]string)
					err := json.Unmarshal([]byte(rawAliases), &storedAliases)
					if err != nil {
						return nil, errors.New("error reading references")
					}
					utils.MergeMaps(&references, &storedAliases)
				} else {
					return nil, errors.New("references attribute not found on current schema. please destroy and recreate again.")
				}
			}
		}
	}

	return references, nil
}

func retrieveLocalState() (map[string]string, error) {
	log.Println("LOCAL STATE")
	fileData, err := utils.LoadData[interface{}](STATE_FILE)
	if err != nil {
		return nil, err
	}
	references, err := readState(fileData)
	if err != nil {
		return nil, err
	}

	return references, nil // read from local terraform.tfstate fle
}

func donwloadStateFromTFC(ctx context.Context, config map[string]string) (interface{}, error) {
	// Set up a TFE client
	client, err := tfe.NewClient(&tfe.Config{
		Token: config["tfcToken"],
	})
	if err != nil {
		return nil, errors.New("Error creating TFE client")
	}

	// Get the latest state for the workspace
	states, err := client.StateVersions.List(ctx, &tfe.StateVersionListOptions{
		Workspace:    config["workspaceName"],
		Organization: config["orgName"],
	})
	if err != nil {
		return nil, errors.New("Error getting terraform states")
	}

	// Get latest state version
	var latestStateVersion *tfe.StateVersion
	for _, sv := range states.Items {
		if latestStateVersion == nil || sv.Serial > latestStateVersion.Serial {
			latestStateVersion = sv
		}
	}

	if latestStateVersion != nil {
		// Print latest state version
		fmt.Printf("Latest state version is %d\n", latestStateVersion.Serial)

		latestState, err := client.StateVersions.Download(ctx, latestStateVersion.DownloadURL)
		if err != nil {
			return nil, errors.New("Error downloading the latest state")
		}

		var stateData interface{}
		err = json.Unmarshal(latestState, &stateData)
		if err != nil {
			return nil, errors.New("Error parsing to json")
		}
		return stateData, nil
	}

	return nil, nil // no states found
}

func retrieveRemoteState(ctx context.Context) (map[string]string, error) {
	log.Println("REMOTE STATE")
	fileData, err := utils.LoadData[interface{}](path.Join(".terraform", STATE_FILE))
	if err != nil {
		return nil, err
	}
	config := make(map[string]string)
	state := fileData.(map[string]interface{})
	backend := state["backend"].(map[string]interface{})
	backendType := backend["type"].(string)
	if backendType == "cloud" {
		token, err := getTFCCredential()
		if err != nil {
			return nil, err
		}
		cloudConfig := backend["config"].(map[string]interface{})
		workSpacesConfig := cloudConfig["workspaces"].(map[string]interface{})
		config["tfcToken"] = token
		config["orgName"] = cloudConfig["organization"].(string)
		config["workspaceName"] = workSpacesConfig["name"].(string)

		fileData, err := donwloadStateFromTFC(ctx, config)
		if err != nil {
			return nil, err
		}
		if fileData != nil {
			references, err := readState(fileData)
			if err != nil {
				return nil, err
			}
			return references, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, fmt.Errorf("Backend not supported")
	}
}

func LoadReferencesFromState(ctx context.Context) error {
	references := make(map[string]string)
	aliases := blueprint_config_references.GetInstance()

	// looking for local tfstate
	if v, err := os.Stat(STATE_FILE); err == nil && v.Size() > 0 {
		references, err = retrieveLocalState()
		if err != nil {
			return err
		}
	} else {
		// looking for remote tfstate
		if v, err := os.Stat(path.Join(".terraform", STATE_FILE)); err == nil && v.Size() > 0 {
			remoteReferences, err := retrieveRemoteState(ctx)
			if err != nil {
				return err
			}
			if remoteReferences != nil {
				references = remoteReferences
			}
		}
	}

	for key, value := range references {
		aliases.SetValue(key, value)
	}
	return nil
}
