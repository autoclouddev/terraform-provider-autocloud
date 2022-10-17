package autocloud_sdk

import (
	"context"
	"fmt"
	"os"
	"testing"

	getter "github.com/hashicorp/go-getter"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestTemplateFromGit(t *testing.T) {
	source := "git::ssh://git@gitlab.com/auto-cloud/infrastructure/infrastructure-modules//aws/storage/s3/bucket"
	version := "0.91.0"
	name := "s3_generator"
	destination := "/tmp/module_source"
	m, _ := NewModule(source, version, name, destination)

	assert.Greater(t, len(m.Variables), 0, "no variables were processed")
	fmt.Println(m.ToForm())
	assert.Greater(t, len(m.ToForm()), 0, "A form was created")
	os.Remove(destination)
}

func TestTemplateFromRegistry(t *testing.T) {
	source := "citizen.tools.autocloud.dev/general/database/postgres"
	version := "0.1.0"
	name := "postgres_db"
	destination := "/tmp/module_source"
	m, _ := NewModule(source, version, name, destination)

	assert.Greater(t, len(m.Variables), 0, "no variables were processed")
	fmt.Println(m.ToForm())
	assert.Greater(t, len(m.ToForm()), 0, "A form was created")
	os.Remove(destination)
}

func TestValidateSourceUrl(t *testing.T) {

	source := "terraform-aws-modules/s3-bucket/aws"
	url, _ := GetModuleUrl(source)
	expected := PublicRegistry + "/" + source
	if expected != url {
		t.Fatalf("url was not converted using the public registry")
	}

	source = "citizen.tools.autocloud.dev/general/database/postgres"
	url, _ = GetModuleUrl(source)
	expected = "https://" + source
	if expected != url {
		t.Fatalf("url was not detected")
	}

}

func TestComplementDownloadUrl(t *testing.T) {
	registryUrl := "citizen.tools.autocloud.dev"
	downloadUrl := "/v1/modules/tarball/general/database/postgres/0.1.0/module.tar.gz"
	fullUrl, _ := ComplementDownloadUrl(registryUrl, downloadUrl)
	fmt.Println("fullurl", fullUrl)
	if fullUrl != "https://"+registryUrl+downloadUrl {
		t.Fatalf("not valid url source")
	}

}
func TestDownloadModuleGit(t *testing.T) {
	t.Skip("this is a debug tool")
	/*
		THIS FUNCTION IS TO TEST THE DOWNLOAD FUNCTION FOR GOGETTER
	*/
	//source := "https://archivist.terraform.io/v1/object/dmF1bHQ6djI6YnN5c0lBS1lmdEhMV0tjVXcrQUZvNmtTZmp1NzJBYnVvTG8zM1pVcXNHZm9mTnpwZVd4WURGNUhBZVB3N2tqQVBmR1hxemZnV2lrejFNVXdQL0h6dW85UHBWTi9ka2lMZ0RIS0xBY3lha0pNQzlvMGsraTA4WXVPRGZ2MTFPRU1JUjBTNXRyUktLcWM1Nnk0RmxZM2NyVWVWekQ0dUtydE15aTkyYi9uTXRFKzhJcU5Fb2tSY3RoTWt0T0NJTzdHYVJIbGZTT2p1b2huTysrRDA1N1ZNZ0NocGw3WkVmV3p5YmtLaVZhREsyOFplb3daVGt5QTV4eXJ5MTJZSVF4SkN2aDVPMW5LZkNRRjQzSXh2QitRR1RBbW45UER6QlRDRHhibWQ5Rml4QkFtYVJLUW84RnlUR21YSEhkMGU5ZkhEWUxLZjdVUHhvZXE5VFA3OUhHMzZjbW5FcnM9?archive=tgz"
	//source := "git::https://github.com/terraform-aws-modules/terraform-aws-s3-bucket?ref=v1.0.0"
	//source := "https://citizen.tools.autocloud.dev/v1/modules/general/database/postgres/0.1.0/download" get download
	//source := "git::ssh://git@gitlab.com/auto-cloud/infrastructure/infrastructure-modules//aws/storage/s3"
	source := "https://citizen.tools.autocloud.dev/v1/modules/tarball/general/database/postgres/0.1.0/module.tar.gz"
	fileSystemDir := "/tmp/gogetter"
	client := &getter.Client{
		Ctx: context.Background(),
		//define the destination to where the directory will be stored. This will create the directory if it doesnt exist
		Dst:  fileSystemDir,
		Dir:  true,
		Src:  source,
		Mode: getter.ClientModeDir,
	}
	if err := client.Get(); err != nil {
		fmt.Printf("error %v", err)
	}
	//assert.Equal(t, gitsource, CleanSourceUrl(gitsource, ""), "they should be equal")
	// check if dir exists
}

func TestCleanSourceUrl(t *testing.T) {
	assert.Equal(t, true, true, "they should be equal")

	defer gock.Off()
	gitsource := "git::ssh://git@gitlab.com/auto-cloud/infrastructure/infrastructure-modules//aws/storage/s3"
	// git
	assert.Equal(t, gitsource, CleanSourceUrl(gitsource, ""), "they should be equal")
	publicModule := "terraform-aws-modules/s3-bucket/aws"
	publicSrcModule := "git::https://github.com/terraform-aws-modules/terraform-aws-s3-bucket?ref=v3.2.3"
	gock.New("https://registry.terraform.io").
		Get("/v1/modules/"+publicModule).
		Reply(200).SetHeader("x-terraform-get", publicSrcModule)
	// public registry
	assert.Equal(t, publicSrcModule, CleanSourceUrl(publicModule, "3.2.3"), "they should be equal")

	//private registry (CITIZEN)
	citizenModule := "citizen.tools.autocloud.dev/general/database/postgres"
	citizenSrcModule := "https://citizen.tools.autocloud.dev/v1/modules/tarball/general/database/postgres/0.1.0/module.tar.gz"
	gock.New("https://citizen.tools.autocloud.dev").
		Get("/v1/modules/general/database/postgres").
		Reply(200).SetHeader("x-terraform-get", citizenSrcModule)
	assert.Equal(t, citizenSrcModule, CleanSourceUrl(citizenModule, "0.1.0"), "they should be equal")

}

func TestGetTFRegistryUrl(t *testing.T) {
	privateModule := "citizen.tools.autocloud.dev/general/database/postgres"
	version := "0.1.0"
	result := GetTFRegistryUrl(privateModule, version)
	fmt.Println(result)
	assert.Equal(t, result, "https://citizen.tools.autocloud.dev/v1/modules/general/database/postgres/0.1.0/download")
	publicModule := "terraform-aws-modules/s3-bucket/aws"
	version = "3.2.3"
	result = GetTFRegistryUrl(publicModule, version)
	fmt.Println(result)
	assert.Equal(t, result, "https://registry.terraform.io/v1/modules/terraform-aws-modules/s3-bucket/aws/3.2.3/download")
}

func TestIsUrl(t *testing.T) {

	assert.Equal(t, isUrl(""), false)
	assert.Equal(t, isUrl("https://registry.terraform.io/v1/modules/terraform-aws-modules/s3-bucket/aws/3.2.3/download"), true)
	assert.Equal(t, isUrl("/v1/modules/terraform-aws-modules/s3-bucket/aws/3.2.3/download"), false)

}
