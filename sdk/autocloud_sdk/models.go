package autocloud_sdk

type IacCatalog struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Author          string           `json:"author"`
	Slug            string           `json:"slug"`
	Source          string           `json:"source"`
	Version         string           `json:"version"`
	Template        string           `json:"template"`
	FormShape       string           `json:"formShape"`
	Description     string           `json:"description"`
	Instructions    string           `json:"instructions"`
	Labels          []string         `json:"labels"`
	FileDefinitions []IacCatalogFile `json:"fileDefinitions"`
}

type IacCatalogFile struct {
	Action           string            `json:"action"`
	PathFromRoot     string            `json:"pathFromRoot"`
	FilenameTemplate string            `json:"fileNameTemplate"`
	FilenameVars     map[string]string `json:"fileNameVars"`
}

type IacCatalogInput struct {
	Name            string           `json:"name"`
	Author          string           `json:"author"`
	Slug            string           `json:"slug"`
	Source          string           `json:"source"`
	Version         string           `json:"version"`
	Description     string           `json:"description"`
	Instructions    string           `json:"instructions"`
	Labels          []string         `json:"labels"`
	FileDefinitions []IacCatalogFile `json:"fileDefinitions"`
	Template        string           `json:"template"`
	FormShape       string           `json:"formShape"`
}

type FormQuestion struct {
	FieldId         string        `json:"fieldId"`
	FieldType       string        `json:"fieldType"`
	FieldLabel      string        `json:"fieldLabel"`
	ExplainingText  string        `json:"explainingText"`
	ValidationRules []interface{} `json:"validationRules"`
}

type FormShape struct {
	Id           string       `json:"id"`
	Type         string       `json:"type"`
	Module       string       `json:"module"`
	FormQuestion FormQuestion `json:"formQuestion"`
}

type Repository struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	HtmlUrl     string `json:htmlUrl`
	Description string `json:description`
}
