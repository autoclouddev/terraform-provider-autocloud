package autocloud_sdk

func GetIacCatalogInput(iacCatalog IacCatalog) IacCatalogInput {
	tfmodule := NewModule(iacCatalog.Source, iacCatalog.Version, iacCatalog.Name)

	return IacCatalogInput{
		Name:            iacCatalog.Name,
		ModuleName:      iacCatalog.ModuleName,
		Author:          iacCatalog.Author,
		Slug:            iacCatalog.Slug,
		Description:     iacCatalog.Description,
		Instructions:    iacCatalog.Instructions,
		Labels:          iacCatalog.Labels,
		FileDefinitions: iacCatalog.FileDefinitions,
		Template:        tfmodule.ToString(),
		FormShape:       tfmodule.ToForm(),
		GitConfig:       iacCatalog.GitConfig,
		GeneratorConfigLocation: iacCatalog.GeneratorConfigLocation,
		GeneratorConfigJson:     iacCatalog.GeneratorConfigJson,
	}
}