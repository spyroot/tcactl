// Package testing
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Mustafa mbayramo@vmware.com
package tca

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spyroot/tcactl/cmd/api"
	"github.com/spyroot/tcactl/cmd/client"
	"log"
)

// Provider -
func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TCA_USERNAME", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("TCA_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"tca_cnfs": dataSourceCnfs(),
		},

		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	log.Println("Connecting")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if (username != "") && (password != "") {
		c, err := api.NewTcaApi(&client.RestClient{
			BaseURL:  "https://tca-vip03.cnfdemo.io",
			ApiKey:   "",
			SkipSsl:  true,
			Client:   nil,
			IsDebug:  true,
			Username: username,
			Password: password,
		},
		)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}

	c, err := api.NewTcaApi(&client.RestClient{
		BaseURL:  "",
		ApiKey:   "",
		SkipSsl:  true,
		Client:   nil,
		IsDebug:  true,
		Username: username,
		Password: password,
	},
	)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
