package blueprintconfiglow

import "github.com/hashicorp/terraform-plugin-go/tftypes"

type blueprintConfigState struct {
	id              string
	source          map[string]string
	config          string
	blueprintConfig string
}

func (in *blueprintConfigState) FromTerraform5Value(val tftypes.Value) error {
	// this is an object type, so we're always going to get a
	// `tftypes.Value` that coerces to a map[string]tftypes.Value
	// as input
	v := map[string]tftypes.Value{}
	err := val.As(&v)
	if err != nil {
		return err
	}

	// now that we can get to the tftypes.Value for each field,
	// call its As method and assign the result to the appropriate
	// variable.

	err = v["id"].As(&in.id)
	if err != nil {
		return err
	}

	err = v["config"].As(&in.config)
	if err != nil {
		return err
	}

	err = v["blueprint_config"].As(&in.blueprintConfig)
	if err != nil {
		return err
	}

	s := map[string]tftypes.Value{}
	err = v["source"].As(&s)
	if err != nil {
		return err
	}
	in.source = make(map[string]string, len(s))
	for k, v := range s {
		var str string
		err = v.As(&str)
		if err != nil {
			return err
		}
		in.source[k] = str
	}

	return nil
}
