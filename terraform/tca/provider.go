// Package tca
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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spyroot/tcactl/lib/api"
	"github.com/spyroot/tcactl/lib/client"
	"log"
)

// New - return new instance
func New() *schema.Provider {
	//
	return &schema.Provider{
		//
		Schema: map[string]*schema.Schema{
			"tca_username": {
				Type:        schema.TypeString,
				Optional:    false,
				DefaultFunc: schema.EnvDefaultFunc("TCA_USERNAME", nil),
			},
			"tca_password": {
				Type:        schema.TypeString,
				Optional:    false,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("TCA_PASSWORD", nil),
			},
			"tca_url": {
				Type:        schema.TypeString,
				Optional:    false,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("TCA_URL", nil),
			},
			"tca_ssl_verify": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   false,
				DefaultFunc: schema.EnvDefaultFunc("TCA_SSL_VERIFY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"tca_templates": resourceClusterTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"instances":     dataSourceCnfs(),
			"tca_templates": dataSourceTemplates(),
			"tca_clusters":  dataSourceClusters(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

//providerConfigure configure tca provider
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	username := d.Get("tca_username").(string)
	password := d.Get("tca_password").(string)
	sslVerifyEnv := d.Get("tca_ssl_verify")
	sslVerify := false
	if sslVerifyEnv != nil {
		sslVerify = sslVerifyEnv.(bool)
	}

	tcaUrl := d.Get("tca_url").(string)

	log.Println("[INFO] Connecting")

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if len(username) == 0 {
		return nil, diag.FromErr(fmt.Errorf("TCA_USERNAME not set"))
	}
	if len(password) == 0 {
		return nil, diag.FromErr(fmt.Errorf("TCA_PASSWORD not set"))
	}
	if len(tcaUrl) == 0 {
		return nil, diag.FromErr(fmt.Errorf("TCA_URL not set"))
	}

	restClient, err := client.NewRestClient(tcaUrl, sslVerify, username, password)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	c, err := api.NewTcaApi(restClient)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, diags
}
