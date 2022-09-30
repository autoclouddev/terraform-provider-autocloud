package autocloud_sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/kjk/betterguid"
	"github.com/stoewer/go-strcase"
	"github.com/zclconf/go-cty/cty"
)

type Module struct {
	name          string
	variables     map[string]interface{}
	source        string // where the module is located in the registry
	version       string
	fileSystemDir string // where we wil process the files
}

func NewModule(source string, version string, name string) *Module {
	log.Printf("initializing module %s, source:%s  version: %s\n", name, source, version)
	// in here we should download a terraform module from a remote registry
	fileSystemDir := fmt.Sprintf("/tmp/%s", betterguid.New())
	DownloadModulePublicRegistry(fileSystemDir, source, version)
	variables, err := processVariables(fileSystemDir)
	DeleteDir(fileSystemDir)
	if err != nil {
		panic("Error reading variables")
	}
	return &Module{
		name:          name,
		source:        source,
		fileSystemDir: fileSystemDir,
		version:       version,
		variables:     variables,
	}
}

// supports only public registries for now
func DownloadModulePublicRegistry(fileSystemDir string, moduleSource string, version string) {
	//registry -> https://registry.terraform.io/v1/modules/
	// source-> terraform-aws-modules/s3-bucket/aws/3.4.0/download
	// get x-terraform-get header
	// use go-getter to fetch the source code
	// download source in the tmp dir and recieve the full path as input
	// get directory and use processVariables
	//fileSystemDir := "/tmp/gogetter"
	publicRegistry := "https://registry.terraform.io/v1/modules"
	terraformRegistryModuleUrl := moduleSource //"terraform-aws-modules/s3-bucket/aws"
	//version := m.version
	getSourceUrl := fmt.Sprintf("%s/%s/%s/download", publicRegistry, terraformRegistryModuleUrl, version)

	res, err := http.Get(getSourceUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	url := res.Header.Get("x-terraform-get")
	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: fileSystemDir,
		Dir: true,
		//the repository with a subdirectory I would like to clone only
		Src:  url, //"github.com/hashicorp/terraform/examples/cross-provider",
		Mode: getter.ClientModeDir,
		//define the type of detectors go getter should use, in this case only github is needed
		Detectors: []getter.Detector{
			&getter.GitHubDetector{},
		},
		//provide the getter needed to download the files
		Getters: map[string]getter.Getter{
			"git": &getter.GitGetter{},
		},
	}
	if err := client.Get(); err != nil {
		log.Fatalf("Error getting path %s: %v\n", client.Src, err)
		return
	}

}

func DeleteDir(fileSystemDir string) {
	err := os.RemoveAll(fileSystemDir)
	if err != nil {
		log.Fatal(err)
	}
}

func processVariables(source string) (map[string]interface{}, error) {
	module, diags := tfconfig.LoadModule(source)
	if diags != nil {
		fmt.Print(errors.New(diags.Error()))
		return nil, errors.New(diags.Error())
	}
	variables := make(map[string]interface{})
	// in here we could insert the AST processing
	for varName := range module.Variables {
		if module.Variables[varName].Type == "string" || module.Variables[varName].Type == "number" {
			// insert default value
			content := make(map[string]interface{})
			content["name"] = module.Variables[varName].Name
			content["default"] = "default value"
			content["description"] = module.Variables[varName].Description
			content["handlebars"] = fmt.Sprintf("{{%s}}", strcase.UpperCamelCase(varName))
			variables[varName] = content
		}
	}
	return variables, nil
}

func (m Module) ToString() string {
	hclFile := hclwrite.NewEmptyFile()
	rootBody := hclFile.Body()
	moduletxt := rootBody.AppendNewBlock("module",
		[]string{m.name})
	moduleBody := moduletxt.Body()

	// adding source and version
	moduleBody.SetAttributeValue("source", cty.StringVal(m.source))
	moduleBody.SetAttributeValue("version", cty.StringVal(m.version))
	for k, v := range m.variables {
		content := v.(map[string]interface{})
		str := (content["handlebars"]).(string)
		moduleBody.SetAttributeValue(k, cty.StringVal(str))
	}
	return fmt.Sprintf("%s", hclFile.Bytes())
}

func (m Module) ToForm() string {

	var formVariables []FormShape

	for _, v := range m.variables {
		f := FormShape{
			Id:     m.name,
			Type:   "string",
			Module: strcase.UpperCamelCase(m.name),
			FormQuestion: FormQuestion{
				FieldId:    "",
				FieldType:  "shortText",
				FieldLabel: "",
			},
		}
		content := v.(map[string]interface{})
		description := (content["description"]).(string)
		name := (content["name"]).(string)
		fieldId := fmt.Sprintf("%s.%s", strcase.UpperCamelCase(m.name), strcase.UpperCamelCase(name))
		f.Id = fieldId
		f.FormQuestion.FieldId = fieldId
		f.FormQuestion.FieldLabel = name
		f.FormQuestion.ExplainingText = description
		//f.FormQuestion.ValidationRules = append(f.FormQuestion.ValidationRules, v) //ref for future
		formVariables = append(formVariables, f)

	}
	j, _ := json.MarshalIndent(formVariables, "", "  ")
	return string(j)
}
