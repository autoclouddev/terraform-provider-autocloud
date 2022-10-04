package autocloud_sdk

type IacCatalog struct {
	ID              string              `json:"id"`
	Name            string              `json:"name"`
	ModuleName      string              `json:"moduleName"`
	Author          string              `json:"author"`
	Slug            string              `json:"slug"`
	Source          string              `json:"source"`
	Version         string              `json:"version"`
	Template        string              `json:"template"`
	FormShape       string              `json:"formShape"`
	Description     string              `json:"description"`
	Instructions    string              `json:"instructions"`
	Labels          []string            `json:"labels"`
	FileDefinitions []IacCatalogFile    `json:"fileDefinitions"`
	GitConfig       IacCatalogGitConfig `json:"gitConfig"`
}

type IacCatalogFile struct {
	Action           string            `json:"action"`
	PathFromRoot     string            `json:"pathFromRoot"`
	FilenameTemplate string            `json:"fileNameTemplate"`
	FilenameVars     map[string]string `json:"fileNameVars"`
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
	HtmlUrl     string `json:"htmlUrl"`
	Description string `json:"description"`
}
type IacCatalogGitConfig struct {
	DestinationBranch string                `json:"destinationBranch"`
	GitUrlOptions     []string              `json:"gitUrlOptions"`
	GitUrlDefault     string                `json:"gitUrlDefault"`
	PullRequest       IacCatalogGitConfigPR `json:"pullRequest"`
}

type IacCatalogGitConfigPR struct {
	Title                 string            `json:"title"`
	CommitMessageTemplate string            `json:"commitMessageTemplate"`
	Body                  string            `json:"body"`
	Variables             map[string]string `json:"variables"`
}

type IacCatalogInput struct {
	Name            string              `json:"name"`
	ModuleName      string              `json:"moduleName"`
	Author          string              `json:"author"`
	Slug            string              `json:"slug"`
	Source          string              `json:"source"`
	Version         string              `json:"version"`
	Description     string              `json:"description"`
	Instructions    string              `json:"instructions"`
	Labels          []string            `json:"labels"`
	FileDefinitions []IacCatalogFile    `json:"fileDefinitions"`
	Template        string              `json:"template"`
	FormShape       string              `json:"formShape"`
	GitConfig       IacCatalogGitConfig `json:"gitConfig"`
}
