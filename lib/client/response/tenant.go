package response

import (
	"encoding/json"
	"github.com/spyroot/tcactl/lib/models"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

// TenantsDetails Tenant Cloud Details
type TenantsDetails struct {
	// tenant id in TCA spec
	TenantID string `json:"tenantId" yaml:"tenantId"`
	// vim name
	VimName                       string                  `json:"vimName" yaml:"vimName"`
	TenantName                    string                  `json:"tenantName" yaml:"tenantName"`
	HcxCloudURL                   string                  `json:"hcxCloudUrl" yaml:"hcxCloudUrl"`
	Username                      string                  `json:"username" yaml:"username"`
	Password                      string                  `json:"password,omitempty" yaml:"password"`
	VimType                       string                  `json:"vimType" yaml:"vimType"`
	VimURL                        string                  `json:"vimUrl" yaml:"vimUrl"`
	HcxUUID                       string                  `json:"hcxUUID" yaml:"hcxUUID"`
	HcxTenantID                   string                  `json:"hcxTenantId" yaml:"hcxTenantId"`
	Location                      *models.Location        `json:"location" yaml:"location"`
	VimID                         string                  `json:"vimId" yaml:"vimId"`
	Audit                         AuditField              `json:"audit" yaml:"audit"`
	VimConn                       *models.VimConnection   `json:"connection,omitempty" yaml:"connection"`
	Compatible                    bool                    `json:"compatible" yaml:"compatible"`
	ID                            string                  `json:"id" yaml:"id"`
	Name                          string                  `json:"name" yaml:"name"`
	AuthType                      string                  `json:"authType,omitempty" yaml:"authType"`
	ClusterName                   string                  `json:"clusterName,omitempty" yaml:"clusterName"`
	ClusterList                   []ClusterNodeConfigList `json:"clusterNodeConfigList" yaml:"clusterNodeConfigList"`
	HasSupportedKubernetesVersion bool                    `json:"hasSupportedKubernetesVersion" yaml:"hasSupportedKubernetesVersion"`
	ClusterStatus                 string                  `json:"clusterStatus" yaml:"clusterStatus"`
	IsCustomizable                bool                    `json:"isCustomizable" yaml:"isCustomizable"`
}

type TenantSpecs struct {
	CloudOwner    string           `json:"cloud_owner" yaml:"cloud_owner"`
	CloudRegionId string           `json:"cloud_region_id" yaml:"cloud_region_id"`
	VimId         string           `json:"vimId" yaml:"vimId"`
	VimName       string           `json:"vimName" yaml:"vimName"`
	Tenants       []TenantsDetails `json:"tenants" yaml:"tenants"`
}

func (t *TenantSpecs) NewInstance(r io.Reader) (*TenantSpecs, error) {
	return NewTenantSpecs(r)
}

func (t TenantSpecs) InstanceSpecsFromString(s string) (interface{}, error) {
	return TenantSpecsFromString(s)
}

// GetFields return all struct TenantSpecs fields
// in generic map.  It used for output filters to capture
// specific value
func (t *TenantSpecs) GetFields() (map[string]interface{}, error) {

	var m map[string]interface{}

	b, err := json.Marshal(t)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}

	return m, nil
}

type AuditField struct {
	// CreationUser who created uuid
	CreationUser      string `json:"creationUser" yaml:"creationUser"`
	CreationTimestamp string `json:"creationTimestamp" yaml:"creationTimestamp"`
}

type ClusterNodeConfigList struct {
	Labels           []string `json:"labels" yaml:"labels"`
	Id               string   `json:"id" yaml:"id"`
	Name             string   `json:"name" yaml:"name"`
	Status           string   `json:"status" yaml:"status"`
	ActiveTasksCount int      `json:"activeTasksCount" yaml:"activeTasksCount"`
	Compatible       bool     `json:"compatible" yaml:"compatible"`
}

// TenantSpecsFromFile - reads tenant spec from file
// and return TenantSpecs instance
func TenantSpecsFromFile(fileName string) (*TenantSpecs, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return ReadTenantSpec(file)
}

// TenantSpecsFromString take string that hold entire spec
// passed to reader and return TenantSpecs instance
func TenantSpecsFromString(str string) (*TenantSpecs, error) {
	r := strings.NewReader(str)
	return ReadTenantSpec(r)
}

// InvalidTenantSpec error if specs invalid
type InvalidTenantSpec struct {
	errMsg string
}

//
func (m *InvalidTenantSpec) Error() string {
	return m.errMsg
}

// ReadTenantSpec - Read tenant spec from io interface
// detects format and use either yaml or json parse
func ReadTenantSpec(b io.Reader) (*TenantSpecs, error) {

	var spec TenantSpecs

	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	err = json.Unmarshal(buffer, &spec)
	if err == nil {
		return &spec, nil
	}

	return nil, &InvalidTenantSpec{"unknown format"}
}

// NewTenantSpecs create spec from reader
func NewTenantSpecs(r io.Reader) (*TenantSpecs, error) {
	spec, err := ReadTenantSpec(r)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

// GetField - return struct field value
func (t *TenantsDetails) GetField(field string) string {

	r := reflect.ValueOf(t)
	fields, _ := t.GetFields()
	if _, ok := fields[field]; ok {
		f := reflect.Indirect(r).FieldByName(strings.Title(field))
		return f.String()
	}

	return ""
}

// GetFields return TenantsDetails fields name as
// map[string], each key is field name
func (t *TenantsDetails) GetFields() (map[string]interface{}, error) {

	var m map[string]interface{}

	b, err := json.Marshal(t)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(b, &m); err != nil {
		return m, err
	}

	return m, nil
}
