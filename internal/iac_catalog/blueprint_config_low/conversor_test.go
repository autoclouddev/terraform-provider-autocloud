package blueprintconfiglow_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
	blueprintconfiglow "gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config_low"
)

func TestConvertAttributesTypes(t *testing.T) {
	got := blueprint_config.DataSourceBlueprintConfig().CoreConfigSchema()
	for key, v := range got.Attributes {
		originalType := v.Type
		fmt.Println(originalType.GoString())
		att1 := *(*blueprintconfiglow.Attribute)(unsafe.Pointer(v))
		fmt.Println(att1.Type.GoString())
		att := blueprintconfiglow.AttConverter(key, att1)
		fmt.Println("helo", att)
		if v.Computed != att.Computed {
			t.Errorf("not equal, v: %v, att: %v", v.Computed, att.Computed)
		}
		if v.Sensitive != att.Sensitive {
			t.Errorf("not equal, v: %v, att: %v", v.Sensitive, att.Sensitive)
		}
		if v.Description != att.Description {
			t.Errorf("not equal, v: %v, att: %v", v.Description, att.Description)
		}
		if v.Required != att.Required {
			t.Errorf("not equal, v: %v, att: %v", v.Required, att.Required)
		}
	}
}

func TestConvertBlock(t *testing.T) {
	got := blueprint_config.DataSourceBlueprintConfig().CoreConfigSchema()
	for key, v := range got.BlockTypes {
		b := *(*blueprintconfiglow.Block)(unsafe.Pointer(v))
		bk := blueprintconfiglow.BlockConverter(b)
		fmt.Println(key)
		fmt.Println(bk)
	}
}

func TestConvertSchema(t *testing.T) {
	got := blueprint_config.DataSourceBlueprintConfig().CoreConfigSchema()
	schema := *(*blueprintconfiglow.Block)(unsafe.Pointer(got))
	lowLvSchema := blueprintconfiglow.ConvertSchema(&schema)
	fmt.Println(lowLvSchema)
	for k, v := range lowLvSchema.Block.Attributes {
		fmt.Println(k, v.Name, v.Type.String())
	}
	types := lowLvSchema.ValueType()
	obj, ok := types.(tftypes.Object)
	if ok {
		fmt.Println(obj.AttributeTypes)
	}
}

func TestGetLowLevelSchema(t *testing.T) {
	schema := blueprintconfiglow.GetBlueprintConfigLowLevelSchema()
	fmt.Println(schema)
}

func TestTypeConversion(t *testing.T) {
	type testCase struct {
		tftype  tftypes.Type
		ctyType cty.Type
	}
	testCases := map[string]testCase{
		"string": {
			tftype:  tftypes.String,
			ctyType: cty.String,
		},
		"number": {
			tftype:  tftypes.Number,
			ctyType: cty.Number,
		},
		"bool": {
			tftype:  tftypes.Bool,
			ctyType: cty.Bool,
		},
		"list": {
			tftype:  tftypes.List{ElementType: tftypes.String},
			ctyType: cty.List(cty.String),
		},
		"set": {
			tftype:  tftypes.Set{ElementType: tftypes.String},
			ctyType: cty.Set(cty.String),
		},
		"map": {
			tftype:  tftypes.Map{ElementType: tftypes.String},
			ctyType: cty.Map(cty.String),
		},
		"tuple": {
			tftype:  tftypes.Tuple{ElementTypes: []tftypes.Type{tftypes.String, tftypes.Bool, tftypes.Tuple{ElementTypes: []tftypes.Type{tftypes.Bool, tftypes.Number}}}},
			ctyType: cty.Tuple([]cty.Type{cty.String, cty.Bool, cty.Tuple([]cty.Type{cty.Bool, cty.Number})}),
		},
		"object": {
			tftype: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
				"st":    tftypes.String,
				"tuple": tftypes.Tuple{ElementTypes: []tftypes.Type{tftypes.Bool, tftypes.Number}},
				"obj": tftypes.Object{AttributeTypes: map[string]tftypes.Type{
					"inn": tftypes.Number,
				}},
			}},
			ctyType: cty.Object(map[string]cty.Type{
				"st":    cty.String,
				"tuple": cty.Tuple([]cty.Type{cty.Bool, cty.Number}),
				"obj": cty.Object(map[string]cty.Type{
					"inn": cty.Number,
				}),
			}),
		},
		"dynamic": {
			tftype:  tftypes.DynamicPseudoType,
			ctyType: cty.DynamicPseudoType,
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			r := blueprintconfiglow.CtyToTftype(test.ctyType)
			if !r.Equal(test.tftype) {
				t.Fatalf("failed conversion -> %s, t.cty: %s, t.tf: %s, r.tf: %s", name, test.ctyType.GoString(), test.tftype.String(), r.String())
			}
		})
	}
}
