// Package app
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

package vc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/spyroot/tcactl/pkg/os"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
	"net/url"
)

const (
	EnvURL      = "VC_URL"
	EnvUserName = "VC_USERNAME"
	EnvPassword = "VC_PASSWORD"
	EnvInsecure = "VC_INSECURE"
)

var urlDescription = fmt.Sprintf("ESX or vCenter URL [%s]", EnvURL)
var urlFlag = flag.String("url", os.GetEnvString(EnvURL, ""), urlDescription)

// ExtractUrl naive url extractor
func ExtractUrl(hostname string) (*url.URL, error) {
	if len(hostname) > 0 {
		return soap.ParseURL(hostname)
	}
	return soap.ParseURL(*urlFlag)
}

// StringOrEnv check if string not empty return , or get from env,
// otherwise error
func StringOrEnv(str string, envName string) (string, error) {

	if len(str) > 0 {
		return str, nil
	}

	str = os.GetEnvString(envName, "")
	if len(str) > 0 {
		return str, nil
	}

	return "", fmt.Errorf("empty %s string", envName)
}

// Connect Function open connection to vCenter or ESXi.
func Connect(ctx context.Context, hostname string, username string, password string) (*govmomi.Client, error) {

	flag.Parse()
	var vcUrl *url.URL

	_username, err := StringOrEnv(username, EnvUserName)
	if err != nil {
		return nil, err
	}

	_password, err := StringOrEnv(password, EnvPassword)
	if err != nil {
		return nil, err
	}

	_hostname, err := StringOrEnv(hostname, EnvURL)
	if err != nil {
		return nil, err
	}

	vcUrl, err = ExtractUrl(_hostname)
	if err != nil {
		return nil, err
	}

	if len(_username) > 0 && len(_password) > 0 {
		vcUrl.User = url.UserPassword(_username, _password)
	}

	if vcUrl == nil {
		return nil, errors.New("failed extract VC url https://vc_fqdn")
	}

	client, err := govmomi.NewClient(ctx, vcUrl, true)
	if err != nil {
		return nil, err
	}

	if client == nil {
		return nil, errors.New("nil vc client")
	}

	return client, nil
}
