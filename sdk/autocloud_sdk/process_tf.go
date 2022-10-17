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
	"path"
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
	Name          string
	Variables     map[string]ProcVariable
	Source        string // where the module is located in the registry
	Version       string
	FileSystemDir string // where we wil process the files
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

func CreateFileDirectory() (string, error) {

	path, err := os.Getwd()
	if err != nil {
		return "", nil
	}
	//fileSystem -> <executedPath>/.terraform/<uniqueid>
	return filepath.Join(path, ".terraform", betterguid.New()), nil
}

func NewModule(source string, version string, name string, dir string) (*Module, error) {
	log.Printf("initializing module %s, source:%s  version: %s\n", name, source, version)

	var fileSystemDir string
	if fileSystemDir = dir; dir == "" {
		createdDir, err := CreateFileDirectory()
		fileSystemDir = createdDir
		if err != nil {
			return nil, err
		}
	}

	log.Printf("new dir created: %s", fileSystemDir)
	err := DownloadTFModule(fileSystemDir, source, version)
	if err != nil {
		return nil, err
	}
	variables, err := processVariables(fileSystemDir)
	DeleteDir(fileSystemDir)
	if err != nil {
		return nil, err
	}
	return &Module{
		Name:          name,
		Source:        source,
		FileSystemDir: fileSystemDir,
		Version:       version,
		Variables:     variables,
	}, nil
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
	log.Print("module_source: ", moduleSource, "\n")
	if moduleSource == "" {
		return moduleSource
	}
	res, err := http.Get(moduleSource)
	if err != nil {
		log.Fatal("Bad tf fetching: ", moduleSource, "\n")
		return ""
	}
	sourceUrl := res.Header.Get("x-terraform-get")
	log.Println("original terraform source code -> ", sourceUrl)
	//return sourceUrl
	//if host empyt add it
	if isUrl(sourceUrl) || strings.Contains(sourceUrl, "git::") {
		log.Println("valid terraform source code -> ", sourceUrl)
		return sourceUrl
	} else {
		clean_domain := strings.ReplaceAll(moduleSource, "https://", "")
		domain := strings.Split(clean_domain, "/")[0]
		sourceUrl := "https://" + path.Join(domain, sourceUrl)
		log.Println("calculated terraform source code -> ", sourceUrl)
		return sourceUrl
	}
}

func GetTFRegistryUrl(source string, version string) string {
	download := path.Join(version, "download")
	host_parts := strings.Split(source, ".")
	if len(host_parts) == 1 {
		return PublicRegistry + "/" + path.Join(source, download)
	}
	if len(host_parts) > 2 {
		domain := strings.Split(source, "/")[0]
		uri := strings.Join(strings.Split(source, "/")[1:], "/")
		return "https://" + domain + path.Join("/v1/modules/", uri, download)
	}
	return ""
}

func GetModuleUrl(source string) (string, error) {
	fmt.Println(source)
	u, err := url.Parse(source)
	if err == nil {
		if u.Path == "" {
			return "", errors.New("invalid parameter")
		} else {
			if u.Host == "" {
				host_parts := strings.Split(source, ".")
				if len(host_parts) > 2 {
					return "https://" + source, nil
				}
				return PublicRegistry + "/" + source, nil
			}
			log.Fatalln("has valid host")
			return source, nil
		}
	}
	return "", nil
}

func DownloadTFModule(fileSystemDir string, moduleSource string, version string) error {
	url := CleanSourceUrl(moduleSource, version)
	fmt.Println("download url ", url)
	if url == "" {
		return errors.New("invalid source for terraform modules")
	}
	// validate url, see if it has a missing the registry (case citizen)
	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst: fileSystemDir,
		Dir: true,
		//the repository with a subdirectory I would like to clone only
		Src:  url,
		Mode: getter.ClientModeDir,
	}
	if err := client.Get(); err != nil {
		log.Fatalf("Error getting path %s: %v\n", client.Src, err)
		return err
	}
	return nil
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
		[]string{m.Name})
	moduleBody := moduletxt.Body()

	// adding source and version
	moduleBody.SetAttributeValue("source", cty.StringVal(m.Source))
	moduleBody.SetAttributeValue("version", cty.StringVal(m.Version))
	for k, v := range m.Variables {
		moduleBody.SetAttributeValue(k, cty.StringVal(v.handlebars))
	}
	return fmt.Sprintf("%s", hclFile.Bytes())
}

func (m Module) ToForm() string {

	var formVariables []FormShape

	for _, v := range m.Variables {
		f := FormShape{
			Id:     m.Name,
			Type:   "string",
			Module: strcase.UpperCamelCase(m.Name),
			FormQuestion: FormQuestion{
				FieldId:    "",
				FieldType:  "shortText",
				FieldLabel: "",
			},
		}

		description := v.description
		name := v.name
		fieldId := fmt.Sprintf("%s.%s", strcase.UpperCamelCase(m.Name), strcase.UpperCamelCase(name))
		f.Id = fieldId
		f.FormQuestion.FieldId = fieldId
		f.FormQuestion.FieldLabel = name
		f.FormQuestion.ExplainingText = description
		formVariables = append(formVariables, f)

	}
	j, _ := json.MarshalIndent(formVariables, "", "  ")
	return string(j)
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		log.Fatalln(err)

	}
	return err == nil && u.Scheme != "" && u.Host != ""
}

func validateUrl(reg_url string, source_url string) (u string, e error) {
	if isUrl(source_url) {
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
