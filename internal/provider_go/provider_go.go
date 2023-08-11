package provider_go

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

type pluginProviderServer struct {
	providerSchema     *tfprotov5.Schema
	providerMetaSchema *tfprotov5.Schema
	resourceSchemas    map[string]*tfprotov5.Schema
	dataSourceSchemas  map[string]*tfprotov5.Schema
	autocloudClient    *autocloudsdk.Client

	resourceRouter
	dataSourceRouter map[string]func() tfprotov5.DataSourceServer
}

type errUnsupportedDataSource string

type ProviderMeta struct {
	Token    string
	Endpoint string
}

type errUnsupportedResource string

type resourceRouter map[string]tfprotov5.ResourceServer

func (e errUnsupportedDataSource) Error() string {
	return "unsupported data source: " + string(e)
}

func (e errUnsupportedResource) Error() string {
	return "unsupported resource: " + string(e)
}

func (p *pluginProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	return &tfprotov5.GetProviderSchemaResponse{
		Provider:          p.providerSchema,
		ProviderMeta:      p.providerMetaSchema,
		ResourceSchemas:   p.resourceSchemas,
		DataSourceSchemas: p.dataSourceSchemas,
	}, nil
}

func (p *pluginProviderServer) PrepareProviderConfig(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	return nil, nil
}

func (p *pluginProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	resp := &tfprotov5.ConfigureProviderResponse{
		Diagnostics: []*tfprotov5.Diagnostic{},
	}
	meta, err := RetrieveProviderMeta(req)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Error retrieving provider meta values for internal provider.",
			Detail:   fmt.Sprintf("This should never happen; please report it to https://github.com/hashicorp/terraform-provider-tfe/issues\n\nThe error received was: %q", err.Error()),
		})
		return resp, nil
	}

	client, err := autocloudsdk.NewClient(&meta.Endpoint, &meta.Token)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Error getting client",
			Detail:   fmt.Sprintf("Error getting client: %v", err),
		})
		return resp, nil
	}

	p.autocloudClient = client
	return resp, nil
}

func (p *pluginProviderServer) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	return &tfprotov5.StopProviderResponse{}, nil
}

func (p *pluginProviderServer) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	ds, ok := p.dataSourceRouter[req.TypeName]
	if !ok {
		return nil, errUnsupportedDataSource(req.TypeName)
	}
	return ds().ValidateDataSourceConfig(ctx, req)
}

func (p *pluginProviderServer) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	ds, ok := p.dataSourceRouter[req.TypeName]
	if !ok {
		return nil, errUnsupportedDataSource(req.TypeName)
	}
	return ds().ReadDataSource(ctx, req)
}

func (r resourceRouter) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	res, ok := r[req.TypeName]
	if !ok {
		return nil, errUnsupportedResource(req.TypeName)
	}
	return res.ValidateResourceTypeConfig(ctx, req)
}

func (r resourceRouter) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	res, ok := r[req.TypeName]
	if !ok {
		return nil, errUnsupportedResource(req.TypeName)
	}
	return res.UpgradeResourceState(ctx, req)
}

func (r resourceRouter) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	res, ok := r[req.TypeName]
	if !ok {
		return nil, errUnsupportedResource(req.TypeName)
	}
	return res.ReadResource(ctx, req)
}

func (r resourceRouter) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	res, ok := r[req.TypeName]
	if !ok {
		return nil, errUnsupportedResource(req.TypeName)
	}
	return res.PlanResourceChange(ctx, req)
}

func (r resourceRouter) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	res, ok := r[req.TypeName]
	if !ok {
		return nil, errUnsupportedResource(req.TypeName)
	}
	return res.ApplyResourceChange(ctx, req)
}

func (r resourceRouter) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	res, ok := r[req.TypeName]
	if !ok {
		return nil, errUnsupportedResource(req.TypeName)
	}
	return res.ImportResourceState(ctx, req)
}

// WithFlagGate returns a PluginProviderServer, the main entry point for muxing
// it receives a flag to allow resources to be included, so they cant collide with
// other framework implementations
func WithFlagGate(experimental bool) func() tfprotov5.ProviderServer {
	dataSourceSchemas := make(map[string]*tfprotov5.Schema)
	dataSourceRouter := make(map[string]func() tfprotov5.DataSourceServer)
	if experimental {
		// in here we include the resources we want to add to dataSourceSchema and router

		//dataSourceSchemas["autocloud_blueprint_config"] = blueprintconfiglow.GetBlueprintConfigLowLevelSchema()
		//dataSourceRouter["autocloud_blueprint_config"] = blueprintconfiglow.NewDataSourceBlueprintConfig
	}
	// PluginProviderServer returns the implementation of an interface for a lower
	// level usage of the Provider to Terraform protocol.
	// This relies on the terraform-plugin-go library, which provides low level
	// bindings for the Terraform plugin protocol.
	return func() tfprotov5.ProviderServer {
		return &pluginProviderServer{
			providerSchema: &tfprotov5.Schema{
				Block: &tfprotov5.SchemaBlock{
					Attributes: []*tfprotov5.SchemaAttribute{
						{
							Name:      "endpoint",
							Type:      tftypes.String,
							Optional:  true,
							Sensitive: true,
						},
						{
							Name:      "token",
							Type:      tftypes.String,
							Optional:  true,
							Sensitive: true,
						},
					},
				},
			},
			dataSourceSchemas: dataSourceSchemas,
			dataSourceRouter:  dataSourceRouter,
		}
	}
}

func RetrieveProviderMeta(req *tfprotov5.ConfigureProviderRequest) (ProviderMeta, error) {
	meta := ProviderMeta{}
	config := req.Config
	val, err := config.Unmarshal(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"token":    tftypes.String,
			"endpoint": tftypes.String,
		}})

	if err != nil {
		return meta, fmt.Errorf("could not unmarshal ConfigureProviderRequest %w", err)
	}
	var endpoint string
	var token string
	var valMap map[string]tftypes.Value
	err = val.As(&valMap)
	if err != nil {
		return meta, fmt.Errorf("could not set the schema attributes to map %w", err)
	}
	if !valMap["endpoint"].IsNull() {
		err = valMap["endpoint"].As(&endpoint)
		if err != nil {
			return meta, fmt.Errorf("could not set the hostname value to string %w", err)
		}
	}
	if !valMap["token"].IsNull() {
		err = valMap["token"].As(&token)
		if err != nil {
			return meta, fmt.Errorf("could not set the token value to string %w", err)
		}
	}

	meta.Endpoint = endpoint
	meta.Token = token

	return meta, nil
}
