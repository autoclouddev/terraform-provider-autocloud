package autocloud_sdk

func GetIacCatalogInput(iacCatalog IacCatalog, moduleId string) IacCatalogInput {

	return IacCatalogInput{
		Name:            iacCatalog.Name,
		Author:          iacCatalog.Author,
		Slug:            iacCatalog.Slug,
		Description:     iacCatalog.Description,
		Instructions:    iacCatalog.Instructions,
		Labels:          iacCatalog.Labels,
		FileDefinitions: iacCatalog.FileDefinitions,
		GitConfig:       iacCatalog.GitConfig,
		IacModuleIds:    []string{moduleId},
	}
}

func GetIacModule(iacCatalog IacCatalog) IacModule {

	tfmodule := NewModule(iacCatalog.Source, iacCatalog.Version, iacCatalog.Name)

	iacModule := IacModule{
		Name:                    iacCatalog.Name,
		Variables:               tfmodule.ToForm(),
		Template:                tfmodule.ToString(),
		Version:                 iacCatalog.Version,
		Source:                  iacCatalog.Source,
		GeneratorConfigLocation: iacCatalog.GeneratorConfigLocation,
		GeneratorConfigJson:     iacCatalog.GeneratorConfigJson,
	}

	if len(iacCatalog.IacModuleIds) > 0 {
		iacModule.ID = iacCatalog.IacModuleIds[0]
	}

	return iacModule
}

func GetIacModuleInput(iacModule IacModule) IacModuleInput {

	return IacModuleInput{
		ID:                      iacModule.ID,
		Name:                    iacModule.Name,
		Variables:               iacModule.Variables,
		Template:                iacModule.Template,
		Version:                 iacModule.Version,
		Source:                  iacModule.Source,
		DbDefinitions:           "",
		GeneratorConfigLocation: iacModule.GeneratorConfigLocation,
		GeneratorConfigJson:     iacModule.GeneratorConfigJson,
	}
}
