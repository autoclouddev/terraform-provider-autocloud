package provider_go_test

import (
	"errors"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	provider "gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/provider_go"
)

func TestPluginProvider_providerMeta(t *testing.T) {
	cases := map[string]struct {
		endpoint string
		token    string
		err      error
	}{
		"has none": {},
		"has only endpoint": {
			endpoint: "http://localhost:8080/api/v.0.0.1",
		},
		"has only token": {
			token: "secret",
		},
		"has endpoint and token": {
			endpoint: "http://localhost:8080/api/v.0.0.1",
			token:    "secret",
		},
	}

	for name, tc := range cases {
		config, err := tfprotov5.NewDynamicValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"endpoint": tftypes.String,
				"token":    tftypes.String,
			},
		}, tftypes.NewValue(tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"endpoint": tftypes.String,
				"token":    tftypes.String,
			},
		}, map[string]tftypes.Value{
			"endpoint": tftypes.NewValue(tftypes.String, tc.endpoint),
			"token":    tftypes.NewValue(tftypes.String, tc.token),
		}),
		)
		if err != nil {
			log.Printf("error: %v", err)
		}

		req := &tfprotov5.ConfigureProviderRequest{
			Config: &config,
		}

		meta, err := provider.RetrieveProviderMeta(req)
		if !errors.Is(err, tc.err) {
			t.Fatalf("Test %s: should not be error", name)
		}

		if tc.endpoint == "" && meta.Endpoint != "" {
			t.Fatalf("Test %s: hostname was not set in config and meta hostname should be empty in this moment (in retrieveProviderMeta). It is parsed later in within the `getClient` function", name)
		}

		if tc.endpoint != "" && meta.Endpoint != tc.endpoint {
			t.Fatalf("Test %s: hostname was set in config and meta hostname %s  has not been set to what was given %s", name, meta.Endpoint, tc.endpoint)
		}

		if tc.token == "" && meta.Token != "" {
			t.Fatalf("Test %s: token was not set in config and meta.token %s has been incorrectly set", name, meta.Token)
		}

		if tc.token != "" && meta.Token != tc.token {
			t.Fatalf("Test %s: token was set in config and input token %s  does not have the same value in meta %s", name, tc.token, meta.Token)
		}

	}
}
