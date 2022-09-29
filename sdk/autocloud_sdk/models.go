package autocloud_sdk

type IacCatalog struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Author          string           `json:"author"`
	Slug            string           `json:"slug"`
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

type IacCatalogGitConfig struct {
	DestinationBranch string                `json:"destinationBranch"`
	PullRequest       IacCatalogGitConfigPR `json:"pullRequest"`
}

type IacCatalogGitConfigPR struct {
	Title                 string                      `json:"title"`
	CommitMessageTemplate string                      `json:"commitMessageTemplate"`
	Body                  string                      `json:"body"`
	Variables             IacCatalogGitConfigTemplate `json:"variables"`
}

type IacCatalogGitConfigTemplate struct {
	AuthorName  string `json:"authorName"`
	ClusterName string `json:"clusterName"`
}

type IacCatalogInput struct {
	Name            string              `json:"name"`
	Author          string              `json:"author"`
	Slug            string              `json:"slug"`
	Description     string              `json:"description"`
	Instructions    string              `json:"instructions"`
	Labels          []string            `json:"labels"`
	GitConfig       IacCatalogGitConfig `json:"gitConfig"`
	FileDefinitions []IacCatalogFile    `json:"fileDefinitions"`
}
