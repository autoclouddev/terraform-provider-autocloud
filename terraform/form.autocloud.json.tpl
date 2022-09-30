{
  "terraformModules": {
    "EKSGenerator": [
      {
        "id": "EKSGenerator.clusterName",
        "module": "EKSGenerator",
        "type": "string",
        "formQuestion": {
          "fieldId": "EKSGenerator.clusterName",
          "fieldType": "shortText",
          "fieldLabel": "Cluster name",
          "validationRules": [
            {
              "errorMessage": "This field is required",
              "rule": "isRequired"
            }
          ]
        }
      }
    ]
  }
}
