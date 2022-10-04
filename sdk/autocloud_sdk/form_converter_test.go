package autocloud_sdk

import (
	"testing"
)

func TestFormConverter(t *testing.T) {
	t.Skip("skip until the ast is implemented")
	//source := "/Users/enciso/Documents/autocloud/infrastructure-live/modules/aws/storage/s3/bucket"

	//want := regexp.MustCompile(`module ` + moduleName + ` {`)
	rawString := `variable_name = list(object({
    allowed_headers = optional(list(string))
		max_age_seconds = optional(number)
  }))`

	/*
		variable_name = [{
			allowed_headers: ["hola", "qadios"],
			max_age_seconds : 11
		}]

	*/
	output := `variable_name = [
	{{#each variable_name}}
		{
			{{#check allowed_headers}}
			allowed_headers = [
				{{#allowed_headers}}
			 "{{this}}"
			 {{/allowed_headers}}
			]
			{{#check allowed_headers}}
			{{#check max_age_seconds}}
			max_age_seconds = {{max_age_seconds}}
			{{/check max_age_seconds}}
		}
	{{/each variable_name}}`
	if ConvertToForm(rawString) != output {
		t.Fatalf("the output is not matching the form")
	}
}

/*


	HCL -> AST -> Handlebars
         AST -> FormDefinition

TEST CASES TO DEAL LATER

	list(object({
    id      = string
    enabled = bool

    abort_incomplete_multipart_upload_days = optional(number)

    expiration = optional(object({
      date                         = optional(string)
      days                         = optional(number)
      expired_object_delete_marker = optional(bool)
    }))

    filter = optional(object({
      object_size_greater_than = optional(number)
      object_size_less_than    = optional(number)
      prefix                   = optional(string)
      tag = optional(object({
        key   = string
        value = string
      }))
    }))

    noncurrent_version_expiration = optional(object({
      newer_noncurrent_versions = optional(number)
      noncurrent_days           = optional(number)
    }))

    noncurrent_version_transition = optional(object({
      newer_noncurrent_versions = optional(number)
      noncurrent_days           = optional(number)
      storage_class             = string
    }))

    transition = optional(list(object({
      date          = optional(string)
      days          = optional(number)
      storage_class = string
    })))
  }))

	object({
    enabled             = bool
    namespace           = string
    cloud_provider      = string
    account             = string
    region              = string
    environment         = string
    stage               = string
    name                = string
    delimiter           = string
    attributes          = list(string)
    tags                = map(string)
    additional_tag_map  = map(string)
    regex_replace_chars = string
    label_order         = list(string)
    id_length_limit     = number
  })


	list(string)


	map(string)

	bool

	string



*/
