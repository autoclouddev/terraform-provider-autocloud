package autocloud_sdk

type IacCatalog struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Author       string   `json:"author"`
	Slug         string   `json:"slug"`
	Description  string   `json:"description"`
	Instructions string   `json:"instructions"`
	Labels       []string `json:"labels"`
}
