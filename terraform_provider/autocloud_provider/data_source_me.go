package autocloud_provider

import (
	"autocloud_sdk"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMe() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMeRead,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*autocloud_sdk.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	me, err := c.GetMe()
	if err != nil {
		return diag.FromErr(err)
	}

	// user := &autocloud_sdk.User{}
	// user.Me.Name = "enrique"
	// user.Me.Email = "someemail"
	// user.Me.ID = "someID"

	//me := make(map[string]interface{})
	//b, err := json.Marshal(user)
	//err = json.NewDecoder(user).Decode(&me)
	tflog.Trace(ctx, "geting the user"+me.Me.Email)
	fmt.Println(me)
	err = d.Set("id", me.Me.ID)
	if err != nil {
		fmt.Println(err)
		return diag.FromErr(err)

	}
	err = d.Set("email", me.Me.Email)
	if err != nil {
		fmt.Println(err)
		return diag.FromErr(err)

	}
	err = d.Set("name", me.Me.Name)
	if err != nil {
		fmt.Println(err)
		return diag.FromErr(err)

	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
