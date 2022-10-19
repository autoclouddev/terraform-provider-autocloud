package autocloud_sdk

type IacCatalog struct {
	ID                      string              `json:"id"`
	Name                    string              `json:"name"`
	Author                  string              `json:"author"`
	Slug                    string              `json:"slug"`
	Source                  string              `json:"source"`
	Version                 string              `json:"version"`
	Template                string              `json:"template"`
	FormShape               string              `json:"formShape"`
	Description             string              `json:"description"`
	Instructions            string              `json:"instructions"`
	Labels                  []string            `json:"labels"`
	FileDefinitions         []IacCatalogFile    `json:"fileDefinitions"`
	GitConfig               IacCatalogGitConfig `json:"gitConfig"`
	GeneratorConfigLocation string              `json:"generatorConfigLocation"`
	GeneratorConfigJson     string              `json:"generatorConfigJson"`
	IacModuleIds            []string            `json:"iacModuleIds"`
}

type IacCatalogFile struct {
	Action           string            `json:"action"`
	PathFromRoot     string            `json:"pathFromRoot"`
	FilenameTemplate string            `json:"fileNameTemplate"`
	FilenameVars     map[string]string `json:"fileNameVars"`
}

type FieldOption struct {
	Label   string `json:"label"`
	FieldId string `json:"fieldId"`
	Value   string `json:"value"`
	Checked bool   `json:"checked"`
}

type FormQuestion struct {
	FieldId         string        `json:"fieldId"`
	FieldType       string        `json:"fieldType"`
	FieldLabel      string        `json:"fieldLabel"`
	ExplainingText  string        `json:"explainingText"`
	FieldOptions    []FieldOption `json:"fieldOptions"`
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
	Author          string              `json:"author"`
	Slug            string              `json:"slug"`
	Description     string              `json:"description"`
	Instructions    string              `json:"instructions"`
	Labels          []string            `json:"labels"`
	FileDefinitions []IacCatalogFile    `json:"fileDefinitions"`
	GitConfig       IacCatalogGitConfig `json:"gitConfig"`
	IacModuleIds    []string            `json:"iacModuleIds"`
}

type IacModule struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Source                  string `json:"source"`
	Version                 string `json:"version"`
	Template                string `json:"template"`
	Variables               string `json:"variables"`
	DbDefinitions           string `json:"dbDefinitions"`
	GeneratorConfigLocation string `json:"generatorConfigLocation"`
	GeneratorConfigJson     string `json:"generatorConfigJson"`
}

type IacModuleInput struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	Source                  string `json:"source"`
	Version                 string `json:"version"`
	Template                string `json:"template"`
	Variables               string `json:"variables"`
	DbDefinitions           string `json:"dbDefinitions"`
	GeneratorConfigLocation string `json:"generatorConfigLocation"`
	GeneratorConfigJson     string `json:"generatorConfigJson"`
}
