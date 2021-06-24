// Package specs
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
package specs

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/pkg/io"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ApiInterface /**
type ApiInterface struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
	Debug    int    `yaml:"debug_level" json:"debug_level"`
	ApiBase  string `yaml:"api_base" json:"api_base"`
}

// Cluster spec
type Cluster struct {
	Name string       `yaml:"name"`
	Rest ApiInterface `yaml:"rest"`
}

// ServerSpec Server spec
type ServerSpec struct {
	Formation struct {
		Cluster Cluster `yaml:"cluster"`
	} `yaml:"serverSpec"`
	BaseDir string
}

// Read Reads config.yml file and serialize everything
// in ServerSpec struct.
func Read(file string) (ServerSpec, error) {

	var spec ServerSpec
	var base string

	if !io.FileExists(file) {
		if pwd, err := os.Getwd(); err == nil {
			glog.Info("Reading config from ", pwd+"/"+file)
			base = filepath.Join(pwd, file)
		} else {
			return spec, err
		}
	}

	// if file exists, check if location current dir
	dir := filepath.Dir(file)
	if dir == "." {
		pwd, err := os.Getwd()
		if err != nil {
			return spec, err
		}
		base = filepath.Join(pwd, file)
	}

	base = file
	glog.Info("Reading config ", file)

	configYaml, err := ioutil.ReadFile(base)
	if err != nil {
		return spec, err
	}

	log.Println("Parsing server configuration specification file.")
	if err = yaml.Unmarshal(configYaml, &spec); err != nil {
		return spec, err
	}

	if dir == "." {
		_pwd, err := os.Getwd()
		if err != nil {
			return ServerSpec{}, err
		}
		_dir, err := filepath.Abs(_pwd)
		if err != nil {
			return ServerSpec{}, err
		}
		spec.BaseDir = _dir
		glog.Infof("Set current working dir as base %s", spec.BaseDir)
	} else {
		glog.Infof("Set base dir %s", filepath.Dir(file))
		spec.BaseDir = filepath.Dir(file)
	}

	ok, err := io.IsDir(spec.BaseDir)
	if err != nil {
		return spec, fmt.Errorf(err.Error())
	}
	if ok == false {
		return spec, fmt.Errorf("can't find a base dir")
	}

	return spec, nil
}
