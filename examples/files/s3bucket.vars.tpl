[
  {
    "id": "S3Bucket.env",
    "type": "radio",
    "module": "S3Bucket",
    "fieldValue": "null",
    "formQuestion": {
      "fieldId": "S3Bucket.env",
      "fieldType": "radio",
      "fieldLabel": "bucket env (appended question)",
      "fieldOptions": [
        {
          "label": "dev",
          "value": "dev",
          "checked": false,
          "fieldId": "S3Bucket.env-dev"
        },
        {
          "label": "nonprod",
          "value": "nonprod",
          "checked": true,
          "fieldId": "S3Bucket.env-nonprod"
        },
        {
          "label": "prod",
          "value": "prod",
          "checked": false,
          "fieldId": "S3Bucket.env-prod"
        }
      ],
      "explainingText": "bucket env helper text (appended question)",
      "validationRules": [
        {
          "rule": "isRequired",
          "value": "",
          "errorMessage": "invalid"
        }
      ]
    },
    "fieldDataType": "string",
    "fieldDefaultValue": "null"
  }
]
