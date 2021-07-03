package request

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type SpecType string

type RegisterVimSpec struct {
	SpecType    *SpecKind `json:"kind,omitempty" yaml:"kind,omitempty" validate:"required"`
	HcxCloudUrl string    `json:"hcxCloudUrl" yaml:"hcxCloudUrl" validate:"required,url"`
	VimName     string    `json:"vimName" yaml:"vimName" validate:"required"`
	TenantName  string    `json:"tenantName,omitempty" yaml:"tenantName,omitempty"`
	Username    string    `json:"username" yaml:"username" validate:"required"`
	Password    string    `json:"password" yaml:"password" validate:"required"`
}

func (p *RegisterVimSpec) GetKind() *SpecKind {
	return p.SpecType
}

// ProviderSpecsFromFile - reads tenant spec from file
// and return TenantSpecs instance
func ProviderSpecsFromFile(fileName string) (*RegisterVimSpec, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadProviderSpec(file)
}

// ProviderSpecsFromFromString take string that holdw entire spec
// passed to reader and return TenantSpecs instance
func ProviderSpecsFromFromString(str string) (*RegisterVimSpec, error) {
	r := strings.NewReader(str)
	return ReadProviderSpec(r)
}

// ReadProviderSpec - Read spec from io reader
// detects format and use either yaml or json parse
func ReadProviderSpec(b io.Reader) (*RegisterVimSpec, error) {

	var spec RegisterVimSpec

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	return nil, &InvalidClusterSpec{"unknown format"}
}
