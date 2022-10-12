package autocloud_sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/kjk/betterguid"
	"github.com/stoewer/go-strcase"
	"github.com/zclconf/go-cty/cty"
)

type Module struct {
	name          string
	variables     map[string]ProcVariable
	source        string // where the module is located in the registry
	version       string
	fileSystemDir string // where we wil process the files
}

type ProcVariable struct {
	name        string
	description string
	handlebars  string
}

const PublicRegistry = "https://registry.terraform.io/v1/modules"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	HttpClient HTTPClient
)

func init() {
	HttpClient = &http.Client{}
}

func NewModule(source string, version string, name string) *Module {
	log.Printf("initializing module %s, source:%s  version: %s\n", name, source, version)
	// in here we should download a terraform module from a remote registry

	// this is executed from /path where terraform apply is called
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//fileSystem -> <executedPath>/.terraform/<uniqueid>
	fileSystemDir := filepath.Join(path, ".terraform", betterguid.New())

	log.Printf("new dir created: %s", fileSystemDir)
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

func CleanSourceUrl(source string, version string) string {
	if strings.Contains(source, "git::") && strings.Index(source, "git::") == 0 {
		if strings.Contains(source, "?ref=") || version == "" {
			return source
		} else {
			return source + "?ref=" + version
		}
	}
	//
	moduleSource := GetTFRegistryUrl(source, version)
	fmt.Print("module_source", moduleSource, "\n")
	if moduleSource == "" {
		panic("eerror")
	}
	res, err := http.Get(moduleSource)
	fmt.Print("res", res)
	if err != nil {
		fmt.Print("err", err)
		return ""
	}
	source_url := res.Header.Get("x-terraform-get")
	//if host empyt add it
	//src_code, _ := GetModuleUrl(source_url)
	return source_url
}

func GetTFRegistryUrl(source string, version string) string {
	download := "/" + version + "/download"
	host_parts := strings.Split(source, ".")

	if len(host_parts) == 1 {
		return PublicRegistry + "/" + source + download
	}
	if len(host_parts) > 2 { ///v1/modules/general/database/postgres
		domain := strings.Split(source, "/")[0]
		path := strings.Join(strings.Split(source, "/")[1:], "/")
		return "https://" + domain + "/v1/modules/" + path + download
	}
	return ""
}

func GetModuleUrl(source string) (string, error) {
	fmt.Println(source)
	u, err := url.Parse(source)
	fmt.Println("error:", err)
	fmt.Println("hostname", u.Hostname())

	fmt.Println("scheme", u.Scheme)

	if err == nil {
		if u.Path == "" {
			fmt.Printf("helllo, invalid parameter")
			return "", errors.New("invalid parameter")
		} else {
			if u.Host == "" {
				host_parts := strings.Split(source, ".")
				if len(host_parts) > 2 {
					return "https://" + source, nil
				}
				return PublicRegistry + "/" + source, nil
			}
			fmt.Println("has valid host")
			return source, nil
		}

	}

	return "", nil
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

	//https://registry.terraform.io/v1/modules/terraform-aws-modules/s3-bucket/aws/1.0.0/download
	//https://app.terraform.io/api/registry/v1/modules/enciso/s3/aws/1.0.0/download

	// terraformRegistryModuleUrl := moduleSource //"terraform-aws-modules/s3-bucket/aws"
	// getSourceUrl, _ := GetModuleUrl(terraformRegistryModuleUrl)

	// //version := m.version
	// fmt.Println(getSourceUrl)
	// //res, err := http.Get(getSourceUrl)
	// request, _ := http.NewRequest(http.MethodGet, getSourceUrl, nil)
	// request.Header.Add("Accept", "application/json")
	// res, err := HttpClient.Do(request)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// source_url := res.Header.Get("x-terraform-get")
	//	url, err := ComplementDownloadUrl(GetDomainFromUrl(getSourceUrl), source_url)
	url := CleanSourceUrl(moduleSource, version)
	if url == "" {
		// log error somewhere
		panic("invalid url")
	}
	// validate url, see if it has a missing the registry (case citizen)
	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: fileSystemDir,
		Dir: true,
		//the repository with a subdirectory I would like to clone only
		Src:  url, //"github.com/hashicorp/terraform/examples/cross-provider",
		Mode: getter.ClientModeDir,
		//define the type of detectors go getter should use, in this case only github is needed
		// Detectors: []getter.Detector{
		// 	&getter.GitHubDetector{},
		// },
		// //provide the getter needed to download the files
		// Getters: map[string]getter.Getter{
		// 	"git": &getter.GitGetter{},
		// },
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

func processVariables(source string) (map[string]ProcVariable, error) {
	module, diags := tfconfig.LoadModule(source)
	if diags != nil {
		fmt.Print(errors.New(diags.Error()))
		return nil, errors.New(diags.Error())
	}
	variables := make(map[string]ProcVariable)
	// in here we could insert the AST processing
	for varName := range module.Variables {
		if module.Variables[varName].Type == "string" || module.Variables[varName].Type == "number" {
			// insert default value

			content := ProcVariable{
				name:        module.Variables[varName].Name,
				description: module.Variables[varName].Description,
				handlebars:  fmt.Sprintf("{{%s}}", strcase.UpperCamelCase(varName)),
			}
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
		moduleBody.SetAttributeValue(k, cty.StringVal(v.handlebars))
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

		description := v.description
		name := v.name
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

func isUrl(str string) bool {
	u, err := url.Parse(str)
	fmt.Print(err)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func validateUrl(reg_url string, source_url string) (u string, e error) {
	if isUrl(source_url) {
		fmt.Println("yei")
		return source_url, nil
	} else if isUrl(filepath.Join(reg_url, source_url)) {
		return filepath.Join(reg_url, source_url), nil
	}
	return "", errors.New("Not a valid url")

}

func GetDomainFromUrl(url string) string {
	res := strings.ReplaceAll(url, "https://", "")
	return strings.Split(res, "/")[0]
}

func ComplementDownloadUrl(host string, path string) (string, error) {
	if host == "" {
		return "", errors.New("invalid url")
	}
	return "https://" + host + path, nil
}
