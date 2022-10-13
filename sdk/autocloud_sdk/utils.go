package autocloud_sdk

func GetIacCatalogInput(iacCatalog IacCatalog) (*IacCatalogInput, error) {
	tfmodule, err := NewModule(iacCatalog.Source, iacCatalog.Version, iacCatalog.Name, "")

	if err != nil {
		return nil, err
	}
	return &IacCatalogInput{
		Name:                    iacCatalog.Name,
		Author:                  iacCatalog.Author,
		Version:                 iacCatalog.Version,
		Source:                  iacCatalog.Source,
		Slug:                    iacCatalog.Slug,
		Description:             iacCatalog.Description,
		Instructions:            iacCatalog.Instructions,
		Labels:                  iacCatalog.Labels,
		FileDefinitions:         iacCatalog.FileDefinitions,
		Template:                tfmodule.ToString(),
		FormShape:               tfmodule.ToForm(),
		GitConfig:               iacCatalog.GitConfig,
		GeneratorConfigLocation: iacCatalog.GeneratorConfigLocation,
		GeneratorConfigJson:     iacCatalog.GeneratorConfigJson,
	}, nil
}
