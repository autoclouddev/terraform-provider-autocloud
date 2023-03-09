package blueprintconfiglow

import (
	"unsafe"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/iac_catalog/blueprint_config"
)

func GetBlueprintConfigLowLevelSchema() *tfprotov5.Schema {
	sch := blueprint_config.DataSourceBlueprintConfig().CoreConfigSchema()
	schema := *(*Block)(unsafe.Pointer(sch))
	return ConvertSchema(&schema)
}

func AttConverter(name string, att Attribute) *tfprotov5.SchemaAttribute {
	z := tfprotov5.SchemaAttribute{
		Name:            name,
		Type:            CtyToTftype(att.Type),
		Description:     att.Description,
		DescriptionKind: tfprotov5.StringKind(att.DescriptionKind),
		Required:        att.Required,
		Computed:        att.Computed,
		Optional:        att.Optional,
		Sensitive:       att.Sensitive,
		Deprecated:      false,
	}
	return &z
}

func CtyToTftype(in cty.Type) tftypes.Type {
	switch {
	case in.Equals(cty.String):
		return tftypes.String
	case in.Equals(cty.Number):
		return tftypes.Number
	case in.Equals(cty.Bool):
		return tftypes.Bool
	case in.IsSetType():
		tftype := CtyToTftype(in.ElementType())
		return tftypes.Set{ElementType: tftype}
	case in.IsListType():
		tftype := CtyToTftype(in.ElementType())
		return tftypes.List{ElementType: tftype}
	case in.IsMapType():
		tftype := CtyToTftype(in.ElementType())
		return tftypes.Map{ElementType: tftype}

	case in.IsObjectType():
		attrs := make(map[string]tftypes.Type)
		for k, v := range in.AttributeTypes() {
			attrs[k] = CtyToTftype(v)
		}
		return tftypes.Object{AttributeTypes: attrs}
	case in.IsTupleType():
		elems := make([]tftypes.Type, 0)

		for idx := 0; idx < in.Length(); idx++ {
			tftype := CtyToTftype(in.TupleElementType(idx))
			elems = append(elems, tftype)
		}
		return tftypes.Tuple{ElementTypes: elems}
	case in.HasDynamicTypes():
		return tftypes.DynamicPseudoType
	}

	return tftypes.DynamicPseudoType
}

func BlockConverter(bk Block) *tfprotov5.SchemaBlock {
	var att = make([]*tfprotov5.SchemaAttribute, 0)
	for name, v := range bk.Attributes {
		att = append(att, AttConverter(name, *v))
	}

	blockTypes := make([]*tfprotov5.SchemaNestedBlock, 0)
	for name, v := range bk.BlockTypes {
		b := *(*NestedBlock)(unsafe.Pointer(v))
		blockTypes = append(blockTypes, NestedBlockConverter(name, b))
	}

	z := tfprotov5.SchemaBlock{
		Description:     bk.Description,
		DescriptionKind: tfprotov5.StringKind(bk.DescriptionKind),
		Attributes:      att,

		BlockTypes: blockTypes,
	}

	return &z
}

func NestedBlockConverter(name string, nb NestedBlock) *tfprotov5.SchemaNestedBlock {
	z := tfprotov5.SchemaNestedBlock{
		TypeName: name,
		Block:    BlockConverter(nb.Block),
		Nesting:  tfprotov5.SchemaNestedBlockNestingMode(nb.Nesting - 1), //offset for iota
		MinItems: int64(nb.MinItems),
		MaxItems: int64(nb.MaxItems),
	}
	return &z
}

func ConvertSchema(schema *Block) *tfprotov5.Schema {
	z := tfprotov5.Schema{
		Version: 1,
		Block:   BlockConverter(*schema),
	}
	return &z
}
