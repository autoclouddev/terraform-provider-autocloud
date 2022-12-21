package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

func TestParseAndMapVariables(t *testing.T) {
	json := `
	[
  {
    "id": "s3bucket.attach_require_latest_tls_policy",
    "type": "string",
    "module": "S3bucket",
    "formQuestion": {
      "fieldId": "s3bucket.attach_require_latest_tls_policy",
      "fieldType": "radio",
      "fieldLabel": "attach_require_latest_tls_policy",
      "fieldOptions": [
        {
          "label": "Yes",
          "value": "true",
          "checked": false,
          "fieldId": "s3bucket.attach_require_latest_tls_policy-true"
        },
        {
          "label": "No",
          "value": "false",
          "checked": true,
          "fieldId": "s3bucket.attach_require_latest_tls_policy-false"
        }
      ],
      "explainingText": "Controls if S3 bucket should require the latest version of TLS",
      "validationRules": null
    }
  },
  {
    "id": "s3bucket.attach_elb_log_delivery_policy",
    "type": "string",
    "module": "S3bucket",
    "formQuestion": {
      "fieldId": "s3bucket.attach_elb_log_delivery_policy",
      "fieldType": "radio",
      "fieldLabel": "attach_elb_log_delivery_policy",
      "fieldOptions": [
        {
          "label": "Yes",
          "value": "true",
          "checked": false,
          "fieldId": "s3bucket.attach_elb_log_delivery_policy-true"
        },
        {
          "label": "No",
          "value": "false",
          "checked": true,
          "fieldId": "s3bucket.attach_elb_log_delivery_policy-false"
        }
      ],
      "explainingText": "Controls if S3 bucket should have ELB log delivery policy attached",
      "validationRules": null
    }
  }
]
`
	vars, err := utils.ParseVariables(json)
	assert.Nil(t, err)
	assert.NotNil(t, vars)
	assert.Equal(t, vars[1].ID, "s3bucket.attach_elb_log_delivery_policy")

	varsMap, err := utils.GetVariablesIdMap(json)
	assert.Nil(t, err)
	assert.Equal(t, "s3bucket.attach_require_latest_tls_policy", varsMap["attach_require_latest_tls_policy"])
	assert.Equal(t, "s3bucket.attach_elb_log_delivery_policy", varsMap["attach_elb_log_delivery_policy"])
}
