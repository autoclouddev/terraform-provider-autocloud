package blueprint_config_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"gitlab.com/auto-cloud/infrastructure/public/terraform-provider/internal/acctest"
)

func TestAccDummyOutputs(t *testing.T) {
	//rInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	//fileName := "test-fixtures/state-versions/terraform.tfstate"
	name := "something"
	//t.Cleanup(orgCleanup)

	//waitForOutputs(t, client, orgName, wsName)

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() { acctest.TestAccPreCheck(t) },
		//ProviderFactories: acctest.ProviderFactories,
		ProtoV5ProviderFactories: acctest.CreateMuxFactories(),
		/*ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			"autocloud": func() (tfprotov5.ProviderServer, error) {
				ctx := context.Background()
				providers := []func() tfprotov5.ProviderServer{
					// Example terraform-plugin-sdk/v2 providers
					//provider.New("dev")().GRPCProvider,
					provider_go.PluginProviderServer,
				}

				muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)

				if err != nil {
					return nil, err
				}

				return muxServer.ProviderServer(), nil
			},
		},*/
		Steps: []resource.TestStep{
			{
				Config: testAccTFEOutputs_dataSource(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.autocloud_dummy.foobar", "name", name),
					testAccDummy("autocloud_dummy"),
					// These outputs rely on the values in test-fixtures/state-versions/terraform.tfstate
					// testCheckOutputState("test_output_list_string", &terraform.OutputState{Value: []interface{}{"us-west-1a"}}),
					// testCheckOutputState("test_output_string", &terraform.OutputState{Value: "9023256633839603543"}),
					// testCheckOutputState("test_output_tuple_number", &terraform.OutputState{Value: []interface{}{"1", "2"}}),
					// testCheckOutputState("test_output_tuple_string", &terraform.OutputState{Value: []interface{}{"one", "two"}}),
					// testCheckOutputState("test_output_object", &terraform.OutputState{Value: map[string]interface{}{"foo": "bar"}}),
					// testCheckOutputState("test_output_number", &terraform.OutputState{Value: "5"}),
					// testCheckOutputState("test_output_bool", &terraform.OutputState{Value: "true"}),
				),
			},
		},
	})
}

func testAccTFEOutputs_dataSource(name string) string {
	return fmt.Sprintf(`
  data "autocloud_dummy" "foobar" {
	name = "%s"
	values = {
		hello = "word"
	    foo ="bar"
		alice = {
			bob = "marlyn"
		}
    }
	//values = "hello"

  }

  output "final" {
	value = data.autocloud_dummy.foobar.values
  }
`, name)
}

/*

   [
	"object",
		{
			"alice":
				[
					"object",
					{"bob":"string"}
				],
			"foo":"string",
			"hello":"string"
		}
	]
	map[
			alice:tftypes.Object[
				"bob":tftypes.String
			]
			<"bob":tftypes.String<"marlyn">>
			foo:tftypes.String<"bar">
			hello:tftypes.String<"word">
		]

*/
/*
map[

	alice:tftypes.Object[
		"bob":tftypes.String
	]<"bob":tftypes.String<"marlyn">>
	foo:tftypes.String<"bar">
	hello:tftypes.String<"word">
	]
*/
func testAccDummy(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fmt.Println("outputs")
		fmt.Println(s.RootModule().Outputs)

		fmt.Println("root")
		fmt.Println(s.RootModule())
		// _, ok := s.RootModule().Resources[resourceName]
		// if !ok {
		// 	return fmt.Errorf("Not found: %s", resourceName)
		// }
		/*rawConf := rs.Primary.Attributes["values"]

		fmt.Println(rawConf)*/
		return nil
	}
}
