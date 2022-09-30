package autocloud_sdk

func GetIacCatalogInput(iacCatalog IacCatalog) IacCatalogInput {

	return IacCatalogInput{
		Name:            iacCatalog.Name,
		Author:          iacCatalog.Author,
		Slug:            iacCatalog.Slug,
		Description:     iacCatalog.Description,
		Instructions:    iacCatalog.Instructions,
		Labels:          iacCatalog.Labels,
		FileDefinitions: iacCatalog.FileDefinitions,
		GitConfig:       iacCatalog.GitConfig,
	}
}
