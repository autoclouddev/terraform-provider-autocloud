package autocloud_module

import (
	"context"
	"crypto/md5" // #nosec G501
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/generator"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk/service/iac_module"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

type moduleResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Source          types.String `tfsdk:"source"`
	Version         types.String `tfsdk:"version"`
	Name            types.String `tfsdk:"name"`
	BlueprintConfig types.String `tfsdk:"blueprint_config"`
	Config          types.String `tfsdk:"config"`
	Outputs         types.Map    `tfsdk:"outputs"`
	Header          types.String `tfsdk:"header"`
	Footer          types.String `tfsdk:"footer"`
	GeneratedName   types.String `tfsdk:"generated_name"`
}

type moduleResource struct {
	client *autocloudsdk.Client
}

var (
	_ resource.Resource               = &moduleResource{}
	_ resource.ResourceWithConfigure  = &moduleResource{}
	_ resource.ResourceWithModifyPlan = &moduleResource{}
)

func generateModuleID(moduleSource, orgID,
	moduleName string, //nolint:unparam
) string {
	// Concatenate the three unique strings
	combinedString := moduleSource + orgID  //+ moduleName
	hash := md5.Sum([]byte(combinedString)) // #nosec G401
	shortID := base64.URLEncoding.EncodeToString(hash[:])
	return shortID
}

func NewModuleResource() resource.Resource {
	return &moduleResource{}
}

func GenerateBlueprintConfig(moduleInput iac_module.ModuleInput, orgId string) (blueprint_config.BluePrintConfig, map[string]attr.Value, error) {
	var variables []generator.FormShape
	elems := make(map[string]attr.Value)
	moduleId := generateModuleID(moduleInput.Source, orgId, moduleInput.Name)
	module, err := moduleInput.Build()
	if err != nil {
		return blueprint_config.BluePrintConfig{}, elems, err
	}
	iac_module := module.ToIacModule(moduleId)

	modVars := iac_module.Variables
	err = json.Unmarshal([]byte(modVars), &variables)
	if err != nil {
		return blueprint_config.BluePrintConfig{}, elems, err
	}
	// ensure same order of questions
	compare := func(i, j int) bool {
		return variables[i].ID < variables[j].ID
	}
	sort.Slice(variables, compare)
	for idx := range variables {
		variables[idx].ModuleID = moduleId
	}
	bp := blueprint_config.BluePrintConfig{
		Id:                moduleId, //strconv.FormatInt(time.Now().Unix(), 10),
		Variables:         variables,
		OverrideVariables: make(map[string]blueprint_config.OverrideVariable),
		OmitVariables:     make([]string, 0),
	}
	outputsMap, err := utils.ToStringMap(iac_module.Outputs)

	for k, v := range outputsMap {
		elems[k] = types.StringValue(v)
	}
	if err != nil {
		return blueprint_config.BluePrintConfig{}, elems, err
	}
	return bp, elems, nil
}

func (r *moduleResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Fill in logic.

	// when is destroying this resource do nothing
	if req.Plan.Raw.IsNull() {
		return
	}
	// it should pregenerate blueprint_config and match on apply and plan time

	org, err := r.client.Organization.Get()
	if err != nil {
		resp.Diagnostics.AddError("Error while getting organization", err.Error())
		return
	}

	var localPlan moduleResourceModel
	req.Plan.Get(ctx, &localPlan)
	in := iac_module.ModuleInput{
		Name:    localPlan.Name.ValueString(),
		Source:  localPlan.Source.ValueString(),
		Version: localPlan.Version.ValueString(),
	}
	bp, outputs, err := GenerateBlueprintConfig(in, org.Id)
	if err != nil {
		resp.Diagnostics.AddError("Error while converting module input to form variables", err.Error())
		return
	}
	bpinBytes, _ := json.MarshalIndent(bp, "", "    ")
	localPlan.BlueprintConfig = types.StringValue(string(bpinBytes))
	variablesBytes, err := json.MarshalIndent(bp.Variables, "", "    ")
	if err != nil {
		resp.Diagnostics.AddError("Error while getting blueprint_config", err.Error())
		return
	}
	localPlan.Config = types.StringValue(string(variablesBytes))
	out, diag := types.MapValue(types.StringType, outputs)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}
	localPlan.Outputs = out

	localPlan.GeneratedName = types.StringValue(fmt.Sprintf("%s.%s", localPlan.Name.ValueString(), "generatedName"))

	resp.Plan.Set(ctx, localPlan)
}

// Configure adds the provider configured client to the resource.
func (r *moduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*autocloudsdk.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *autocloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *moduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_module"
}

// Schema defines the schema for the resource.
func (r *moduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Required: true,
			},
			"version": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"blueprint_config": schema.StringAttribute{
				Computed: true,
			},
			"config": schema.StringAttribute{
				Computed: true,
			},
			"outputs": schema.MapAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
			"header": schema.StringAttribute{
				Optional: true,
			},
			"footer": schema.StringAttribute{
				Optional: true,
			},
			"generated_name": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *moduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan moduleResourceModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	org, err := r.client.Organization.Get()
	if err != nil {
		resp.Diagnostics.AddError("Error while getting organization", err.Error())
		return
	}

	// compute blueprintconfig before creating it on the backend
	moduleId := generateModuleID(plan.Source.ValueString(), org.Id, plan.Name.ValueString())

	in := iac_module.ModuleInput{
		ID:      moduleId,
		Name:    plan.Name.ValueString(),
		Source:  plan.Source.ValueString(),
		Version: plan.Version.ValueString(),
		Header:  plan.Header.ValueString(),
		Footer:  plan.Footer.ValueString(),
	}
	bp, outputs, err := GenerateBlueprintConfig(in, org.Id)
	if err != nil {
		resp.Diagnostics.AddError("Error while preparing module", err.Error())
		return
	}
	module, err := r.client.Module.Create(in)
	if err != nil {
		resp.Diagnostics.AddError("Error while preparing module", err.Error())
		return
	}

	plan.ID = types.StringValue(module.ID)

	bpinBytes, err := json.MarshalIndent(bp, "", "    ")
	if err != nil {
		resp.Diagnostics.AddError("Error while getting blueprint_config", err.Error())
		return
	}
	variablesBytes, err := json.MarshalIndent(bp.Variables, "", "    ")
	if err != nil {
		resp.Diagnostics.AddError("Error while getting blueprint_config", err.Error())
		return
	}
	plan.Config = types.StringValue(string(variablesBytes))
	plan.BlueprintConfig = types.StringValue(string(bpinBytes))
	out, diag := types.MapValue(types.StringType, outputs)
	if diag.HasError() {
		resp.Diagnostics.Append(diag...)
		return
	}
	plan.Outputs = out
	plan.GeneratedName = types.StringValue(fmt.Sprintf("%s.%s", plan.Name.ValueString(), "generatedName"))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *moduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state moduleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*// pre generate module input data
	in := iac_module.ModuleInput{
		Name:    state.Name.ValueString(),
		Source:  state.Source.ValueString(),
		Version: state.Version.ValueString(),
	}
	moduleInputdata, err := r.client.Module.Prepare(in)
	if err != nil {
		resp.Diagnostics.AddError("Error while getting module, could not read Module configuration", err.Error())
		return
	}

	module, err := r.client.Module.Get(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error while getting module, could not read Module ID", err.Error())
		return
	}
	moduleInput := moduleInputdata.ToIacModule(module.ID)
	jsonString, _ := json.Marshal(module)
	// mix with pre-generated module input data
	err = json.Unmarshal(jsonString, &moduleInput)

	if err != nil {
		resp.Diagnostics.AddError("Error mixing module input", err.Error())
		return
	}
	formVariables := moduleInputdata.ToForm()
	var variables []generator.FormShape
	err = json.Unmarshal([]byte(formVariables), &variables)

	if err != nil {
		resp.Diagnostics.AddError("Error while converting module input to form", err.Error())
		return
	}
	*/
	//set refresed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *moduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan moduleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	in := iac_module.ModuleInput{
		ID:      plan.ID.ValueString(),
		Name:    plan.Name.ValueString(),
		Source:  plan.Source.ValueString(),
		Version: plan.Version.ValueString(),
		Header:  plan.Header.ValueString(),
		Footer:  plan.Footer.ValueString(),
	}

	module, err := r.client.Module.Update(in)

	if err != nil {
		resp.Diagnostics.AddError("Error while updating module", err.Error())
		return
	}

	var variables []generator.FormShape
	err = json.Unmarshal([]byte(module.Variables), &variables)
	if err != nil {
		resp.Diagnostics.AddError("Error while getting blueprint_config", err.Error())
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *moduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state moduleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Module.Delete(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error while deleting module", err.Error())
		return
	}
}
