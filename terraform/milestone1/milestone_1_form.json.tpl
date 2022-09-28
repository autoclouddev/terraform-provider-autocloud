{
  "terraformModules" : {
    "EKSGenerator": {
        "name": "EKSGenerator",
        "variables": [
            {
                "id": "EKSGenerator.clusterName",
                "module": "EKSGenerator",
                "type": "string",
                "formQuestion": {
                    "fieldId": "EKSGenerator.clusterName",
                    "fieldType": "shortText",
                    "fieldLabel": "Cluster name:",
                    "explainingText": "",
                    "validationRules": [
                        {
                            "rule": "isRequired",
                            "errorMessage": "This field is required"
                        }
                    ]
                }
            }
        ],
        "dbDefinitions": { // we can handle this duplication in the sdk, and refacto the backend to avoid this in the first place
            "clusterName": {
                "id": "EKSGenerator.clusterName",
                "module": "EKSGenerator",
                "type": "string",
                "formQuestion": {
                    "fieldId": "EKSGenerator.clusterName",
                    "fieldType": "shortText",
                    "fieldLabel": "Cluster name:",
                    "explainingText": "",
                    "validationRules": [
                        {
                            "rule": "isRequired",
                            "errorMessage": "This field is required"
                        }
                    ]
                }
            }
        }
    }

  }

}
