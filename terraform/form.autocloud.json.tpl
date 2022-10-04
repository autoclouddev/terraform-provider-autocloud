{
  "terraformModules": {
    "ExampleS3": [
      {
        "id": "ExampleS3.clusterName",
        "module": "ExampleS3",
        "type": "string",
        "formQuestion": {
          "fieldId": "ExampleS3.clusterName",
          "fieldType": "shortText",
          "fieldLabel": "Cluster name",
          "validationRules": [
            {
              "errorMessage": "This field is required",
              "rule": "isRequired"
            }
          ]
        }
      },
      {
        "id": "ExampleS3.AccelerationStatus",
        "type": "string",
        "module": "ExampleS3",
        "formQuestion": {
          "fieldId": "ExampleS3.AccelerationStatus",
          "fieldType": "shortText",
          "fieldLabel": "[THIS IS AN OVERRIDE FROM TPL] Acceleration Status",
          "explainingText": "[THIS IS AN OVERRIDE FROM TPL] (Optional) Sets the accelerate configuration of an existing bucket. Can be Enabled or Suspended.",
          "validationRules": null
        }
      }
    ]
  }
}
