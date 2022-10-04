package autocloud_sdk

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTemplate(t *testing.T) {
	testDir := ".terraform"
	os.Mkdir(testDir, os.ModePerm)
	source := "terraform-aws-modules/s3-bucket/aws"
	version := "3.4.0"
	name := "s3_generator"

	//want := regexp.MustCompile(`module ` + moduleName + ` {`)
	m := NewModule(source, version, name)

	if len(m.variables) == 0 {
		t.Fatalf("no variables were processed")
	}
	//fmt.Println(m.ToString())
	fmt.Println(m.ToForm())
	dir, _ := os.Getwd()
	os.Remove(filepath.Join(dir, testDir))

}

// func TestDownloadModules(t *testing.T) {
// 	//source := "/Users/enciso/Documents/autocloud/infrastructure-live/modules/aws/storage/s3/bucket"
// 	DownloadModulePublicRegistry()
// 	//want := regexp.MustCompile(`module ` + moduleName + ` {`)
// 	//m := NewModule(source, "test")

// 	if false {
// 		t.Fatalf("no variables were processed")
// 	}
// 	//fmt.Println(m.ToString())
// 	// if !want.MatchString(template) || err != nil {
// 	// 	t.Fatalf(`"module test {"  = %q, %v, want match for %#q, nil`, template, err, want)
// 	// }
// }
