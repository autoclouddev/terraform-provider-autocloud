package blueprintconfiglow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/utils"
)

// dataSourceBlueprintConfig
type dsBlueprintConfig struct {
	// this struct can carry important elements for usage in the lifecycle, in this case the sdk is not needed
	// because we are not accessing anything from it
	// Leaving this as reference for the future
	//nolint:golint,unused
	autocloudClient *autocloudsdk.Client
}

func NewDataSourceBlueprintConfig() tfprotov5.DataSourceServer {
	return dsBlueprintConfig{}
}

func (d dsBlueprintConfig) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	objTypeDef, ok := GetBlueprintConfigLowLevelSchema().ValueType().(tftypes.Object)
	if !ok {
		return nil, errors.New("Cant get lowlevel attributes")
	}
	values, err := req.Config.Unmarshal(tftypes.Object{
		AttributeTypes: objTypeDef.AttributeTypes,
	})
	if err != nil {
		return nil, fmt.Errorf("Cant unmarshall config input, %v", err.Error())
	}

	var input blueprintConfigState
	err = values.As(&input)
	if err != nil {
		return nil, fmt.Errorf("Cant convert config input, %v", err.Error())
	}
	// TODO: refactor GetBlueprintConfigFromSchema to be used here
	bp := &blueprint_config.BluePrintConfig{}
	bp.Id = strconv.FormatInt(time.Now().Unix(), 10)
	bp.Children = make(map[string]blueprint_config.BluePrintConfig)
	for k, v := range input.source {
		bc := blueprint_config.BluePrintConfig{}
		err := json.Unmarshal([]byte(v), &bc)
		if err != nil {
			return nil, errors.New("invalid conversion to BluePrintConfig")
		}
		bp.Children[k] = bc
	}
	formVariables, err := blueprint_config.GetFormShape(*bp)
	if err != nil {
		return nil, err
	}
	jsonFormShape, err := utils.ToJsonString(formVariables)
	if err != nil {
		return nil, err
	}

	pretty, err := utils.PrettyStruct(bp)
	if err != nil {
		return nil, err
	}

	state, err := tfprotov5.NewDynamicValue(
		objTypeDef,
		tftypes.NewValue(tftypes.Object{
			AttributeTypes: objTypeDef.AttributeTypes,
		}, map[string]tftypes.Value{
			"id":               tftypes.NewValue(tftypes.String, strconv.FormatInt(time.Now().Unix(), 10)),
			"source":           tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, map[string]tftypes.Value{"hello": tftypes.NewValue(tftypes.String, "world")}),
			"config":           tftypes.NewValue(tftypes.String, jsonFormShape),
			"blueprint_config": tftypes.NewValue(tftypes.String, pretty),
			"variable":         tftypes.NewValue(objTypeDef.AttributeTypes["variable"], []tftypes.Value{}),
			"omit_variables":   tftypes.NewValue(objTypeDef.AttributeTypes["omit_variables"], []tftypes.Value{}),
		}))
	if err != nil {
		return nil, err
	}

	return &tfprotov5.ReadDataSourceResponse{
		State: &state,
	}, nil
}

func (d dsBlueprintConfig) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	return &tfprotov5.ValidateDataSourceConfigResponse{}, nil
}
