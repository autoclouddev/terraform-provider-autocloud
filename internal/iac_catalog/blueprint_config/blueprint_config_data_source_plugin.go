package blueprint_config

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	autocloudsdk "gitlab.com/auto-cloud/infrastructure/public/terraform-provider-sdk"
)

type dataSourceDummy struct {
	autocloudClient *autocloudsdk.Client
}

type stateData struct {
	output *outputData
}

type outputData struct {
	Value cty.Value
}

func NewDataSourceDummy() tfprotov5.DataSourceServer {
	return dataSourceDummy{}
}

func (d dataSourceDummy) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	resp := &tfprotov5.ReadDataSourceResponse{
		Diagnostics: []*tfprotov5.Diagnostic{},
	}

	name, values, err := d.readConfigValues(req)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Error retrieving values from the config",
			Detail:   fmt.Sprintf("Error retrieving values from the config: %v", err),
		})
		return resp, nil
	}
	fmt.Println("VALUES FIRST")
	fmt.Println(values)
	remoteStateOutput, err := d.readStateOutput(ctx, values, name)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Error reading remote state output",
			Detail:   fmt.Sprintf("Error reading remote state output: %v", err),
		})
		return resp, nil
	}

	tfValue, tfType, err := parseStateOutput(remoteStateOutput)
	fmt.Println("VALUE")
	fmt.Println(tfValue)
	fmt.Println(tfValue.Type())
	fmt.Println("TYPE")
	fmt.Println(tfType)
	//fmt.Print(tfValue)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Error parsing remote state output",
			Detail:   fmt.Sprintf("Error parsing remote state output: %v", err),
		})
		return resp, nil
	}

	id := fmt.Sprintf("%s-%s", name, strconv.FormatInt(time.Now().Unix(), 10))
	state, err := tfprotov5.NewDynamicValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":   tftypes.String,
			"values": tftypes.DynamicPseudoType,
			"id":     tftypes.String,
		},
	}, tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name": tftypes.String,
			//"values": tftypes.Object{AttributeTypes: stateTypes},
			"values": tfType,
			"id":     tftypes.String,
		},
	}, map[string]tftypes.Value{
		"name": tftypes.NewValue(tftypes.String, name),
		//"values": tftypes.NewValue(tftypes.Object{AttributeTypes: stateTypes}, tftypesValues),
		"values": tfValue, //tftypes.NewValue(tfValue.Type(), tfValue.String()),
		"id":     tftypes.NewValue(tftypes.String, id),
	}))

	if err != nil {
		return &tfprotov5.ReadDataSourceResponse{
			Diagnostics: []*tfprotov5.Diagnostic{
				{
					Severity: tfprotov5.DiagnosticSeverityError,
					Summary:  "Error encoding state",
					Detail:   fmt.Sprintf("Error encoding state: %s", err.Error()),
				},
			},
		}, nil
	}

	return &tfprotov5.ReadDataSourceResponse{
		State: &state,
	}, nil
}

func (d dataSourceDummy) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	return &tfprotov5.ValidateDataSourceConfigResponse{}, nil
}

// https://github.com/hashicorp/terraform/blob/main/docs/resource-instance-change-lifecycle.md#validateresourceconfig

func (d dataSourceDummy) readConfigValues(req *tfprotov5.ReadDataSourceRequest) (string, map[string]interface{}, error) {
	var name string
	var err error
	var values map[string]interface{} = make(map[string]interface{})

	config := req.Config
	val, err := config.Unmarshal(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":   tftypes.String, //required
			"values": tftypes.DynamicPseudoType,
			"id":     tftypes.String,
		}})
	if err != nil {
		return name, values, fmt.Errorf("Error unmarshalling config: %w", err)
	}

	var valMap map[string]tftypes.Value
	err = val.As(&valMap)
	if err != nil {
		return name, values, fmt.Errorf("Error assigning configuration attributes to map: %w", err)
	}

	if valMap["name"].IsNull() {
		return name, values, fmt.Errorf("name cannot be nil: %w", err)
	}
	err = valMap["name"].As(&name)
	if err != nil {
		return name, values, fmt.Errorf("Error assigning 'name' value to string: %w", err)
	}
	fmt.Println("VALUES VALUE")
	//fmt.Println(valMap["values"].Type())
	//fmt.Printf("testing eq, %v\n", valMap["values"].Type().Is(tftypes.Number))
	var tmp map[string]tftypes.Value = make(map[string]tftypes.Value)
	var strVal string
	var boolVar bool
	var intVar int
	switch {
	case valMap["values"].Type().Is(tftypes.Object{}):
		fmt.Println("VALUES MARSHAL")
		bytes, err := valMap["values"].Type().MarshalJSON()
		fmt.Println(valMap["values"].Type())
		fmt.Println(string(bytes))
		err = valMap["values"].As(&tmp)
		for key, val := range tmp {
			fmt.Println("entry")
			fmt.Println(key)
			fmt.Println(val)
		}
		if err != nil {
			return name, values, fmt.Errorf("Error converting 'values' value to Object: %w", err)
		}
		values["val"] = tmp
		fmt.Println(values["val"])
	case valMap["values"].Type().Is(tftypes.Number):
		err := valMap["values"].As(&intVar)
		if err != nil {
			return name, values, fmt.Errorf("Error converting 'values' value to number: %w", err)
		}
		values["val"] = intVar
	case valMap["values"].Type().Is(tftypes.String):
		err := valMap["values"].As(&strVal)
		if err != nil {
			return name, values, fmt.Errorf("Error converting 'values' value to string: %w", err)
		}
		values["val"] = strVal
	case valMap["values"].Type().Is(tftypes.Bool):
		err := valMap["values"].As(&boolVar)
		if err != nil {
			return name, values, fmt.Errorf("Error converting 'values' value to boolean: %w", err)
		}
		values["val"] = boolVar

	}
	// err = valMap["values"].As(&tmp)
	// if err != nil {
	// 	return name, values, fmt.Errorf("Error converting 'values' value to json: %w", err)
	// }
	// fmt.Println(tmp)

	/*err = valMap["values"].As(&name)

	if err != nil {
		return name, values, fmt.Errorf("Error assigning 'values' value to interface: %w", err)
	}*/

	return name, values, nil
}

func (d dataSourceDummy) readStateOutput(ctx context.Context, values map[string]interface{}, name string) (*stateData, error) {
	fmt.Printf("[DEBUG] Reading the Value %s\n", name)
	//fmt.Println(ctx)    //not using it currently
	//fmt.Println(values) //not using it currently
	/*opts := &tfe.WorkspaceReadOptions{
		Include: []tfe.WSIncludeOpt{tfe.WSOutputs},
	}
	ws, err := tfeClient.Workspaces.ReadWithOptions(ctx, orgName, wsName, opts)
	if err != nil {
		return nil, fmt.Errorf("Error reading workspace: %w", err)
	}
	values[val] = interface{}
	*/
	sd := &stateData{
		output: &outputData{},
	}
	/*
		operations := map[string]struct {
			Value interface{}
			Name  string
		}{
			"first": {
				Value: values,
				Name:  "hello",
			},
		}

		for _, op := range operations {
		}*/
	opValue := values["val"]
	fmt.Println("opValue")
	fmt.Println(opValue)
	buf, err := json.Marshal(opValue)
	if err != nil {
		return nil, fmt.Errorf("could not marshal output value: %w", err)
	}
	fmt.Println(string(buf))

	v := ctyjson.SimpleJSONValue{}
	err = v.UnmarshalJSON(buf)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal output value: %w", err)
	}
	sd.output = &outputData{
		Value: v.Value,
	}
	fmt.Println(sd)
	return sd, nil
}

func parseStateOutput(stateOutput *stateData) (tftypes.Value, tftypes.Type, error) {
	var emptyVar tftypes.Value
	var emptyType tftypes.Type
	fmt.Println("parseStateOutput init")
	fmt.Println(stateOutput.output)
	marshData, err := stateOutput.output.Value.Type().MarshalJSON()
	if err != nil {
		return emptyVar, emptyType, fmt.Errorf("could not marshal output type: %w", err)
	}
	fmt.Println(string(marshData))
	tfType, err := tftypes.ParseJSONType(marshData)
	if err != nil {
		return emptyVar, emptyType, fmt.Errorf("could not parse json type data: %w", err)
	}
	mByte, err := ctyjson.Marshal(stateOutput.output.Value, stateOutput.output.Value.Type())
	if err != nil {
		return emptyVar, emptyType, fmt.Errorf("could not marshal output value and output type: %w", err)
	}
	tfRawState := tfprotov5.RawState{
		JSON: mByte,
	}
	newVal, err := tfRawState.Unmarshal(tfType)
	if err != nil {
		return emptyVar, emptyType, err
	}
	fmt.Println("parseStateOutput")
	fmt.Println(newVal)
	fmt.Println(tfType)
	return newVal, tfType, nil
}

/*
func parseStateOutput(stateOutput *stateData) (map[string]tftypes.Value, map[string]tftypes.Type, error) {
	tftypesValues := map[string]tftypes.Value{}
	stateTypes := map[string]tftypes.Type{}
	//stateOutput.output
	for name, output := range stateOutput.output {
		marshData, err := output.Value.Type().MarshalJSON()
		if err != nil {
			return nil, nil, fmt.Errorf("could not marshal output type: %w", err)
		}
		tfType, err := tftypes.ParseJSONType(marshData)
		if err != nil {
			return nil, nil, fmt.Errorf("could not parse json type data: %w", err)
		}
		mByte, err := ctyjson.Marshal(output.Value, output.Value.Type())
		if err != nil {
			return nil, nil, fmt.Errorf("could not marshal output value and output type: %w", err)
		}
		tfRawState := tfprotov5.RawState{
			JSON: mByte,
		}
		newVal, err := tfRawState.Unmarshal(tfType)
		if err != nil {
			return nil, nil, fmt.Errorf("could not unmarshal tftype into value: %w", err)
		}

		tftypesValues[name] = newVal
		stateTypes[name] = tfType
	}

	// tftypes.NewValue(tftypes.Object{AttributeTypes: stateTypes}, tftypesValues),
	//  tftypes.NewValue(tftypes.String, id),


	return tftypesValues, stateTypes, nil
}
*/
