package specs

import (
	"github.com/asaskevich/govalidator"
	"io"
)

// InvalidInstanceSpec error if specs invalid
type InvalidInstanceSpec struct {
	errMsg string
}

//
func (m *InvalidInstanceSpec) Error() string {
	return m.errMsg
}

// InstanceRequestSpec new instance request
type InstanceRequestSpec struct {
	//
	SpecType SpecType `json:"kind,omitempty" yaml:"kind,omitempty" valid:"required~kind is mandatory spec field"`
	//InstanceName
	InstanceName string `json:"instance_name,omitempty" yaml:"instance_name,omitempty" valid:"required~instance_name is mandatory spec field"`
	// NfdName catalog name
	NfdName string `json:"catalog_name,omitempty" yaml:"catalog_name,omitempty" valid:"required~catalog_name is mandatory spec field"`
	// CloudName target cloud name
	CloudName string `json:"cloud_name,omitempty" yaml:"cloud_name,omitempty" valid:"required~cloud_name is mandatory spec field"`
	// ClusterName target cluster name
	ClusterName string `json:"cluster_name,omitempty" yaml:"cluster_name,omitempty" valid:"required~cluster_name is mandatory spec field"`
	// VimType VC or K8S vim name
	VimType string `json:"vim_type,omitempty" yaml:"vim_type,omitempty"`
	//UseAttached
	UseAttached bool `json:"default_repo,omitempty" yaml:"default_repo,omitempty"`
	//Repo
	Repo string `json:"repo_url,omitempty" yaml:"repo_url,omitempty"`
	//RepoUsername
	RepoUsername string `json:"repo_username,omitempty" yaml:"repo_username,omitempty"`
	//RepoPassword
	RepoPassword string `json:"repo_password,omitempty" yaml:"repo_password,omitempty"`
	//NodePoolName
	NodePoolName string `json:"node_pool,omitempty" yaml:"node_pool,omitempty"`
	// user linked Repo
	UseLinkedRepo bool `json:"use_linked_repo,omitempty" yaml:"use_linked_repo,omitempty"`
	// Namespace target Namespace
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// FlavorName flavor name
	FlavorName string `json:"flavor_name,omitempty" yaml:"flavor_name,omitempty"`
	//Description
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// AdditionalParams additional placement details
	AdditionalParams AdditionalParams `json:"additionalParams,omitempty" yaml:"additionalParams,omitempty"`
	// fix name conflict
	DoAutoName bool `json:"auto_name,omitempty" yaml:"auto_name,omitempty"`

	specError error
}

// Default set default optional fields
func (t *InstanceRequestSpec) Default() error {

	t.DoAutoName = true
	t.NodePoolName = "default"
	t.Namespace = "default"

	t.AdditionalParams = AdditionalParams{
		DisableGrant:        false,
		IgnoreGrantFailure:  false,
		DisableAutoRollback: false,
	}

	return nil
}

func (t *InstanceRequestSpec) Kind() SpecType {
	return t.SpecType
}

func (t *InstanceRequestSpec) GetCloudName() string {
	return t.CloudName
}

func (t *InstanceRequestSpec) SetCloudName(cloudName string) {
	t.CloudName = cloudName
}

func (t *InstanceRequestSpec) SetAutoName(f bool) {
	t.DoAutoName = f
}

func (t *InstanceRequestSpec) IsAutoName() bool {
	return t.DoAutoName
}

func (t *InstanceRequestSpec) GetClusterName() string {
	return t.ClusterName
}

func (t *InstanceRequestSpec) SetClusterName(clusterName string) {
	t.ClusterName = clusterName
}

func (t *InstanceRequestSpec) GetVimType() string {
	return t.VimType
}

func (t *InstanceRequestSpec) SetVimType(vimType string) {
	t.VimType = vimType
}

func (t *InstanceRequestSpec) GetNfdName() string {
	return t.NfdName
}

func (t *InstanceRequestSpec) SetNfdName(nfdName string) {
	t.NfdName = nfdName
}

func (t *InstanceRequestSpec) GetRepo() string {
	return t.Repo
}

func (t *InstanceRequestSpec) SetRepo(repo string) {
	t.Repo = repo
}

func (t *InstanceRequestSpec) GetInstanceName() string {
	return t.InstanceName
}

func (t *InstanceRequestSpec) SetInstanceName(instanceName string) {
	t.InstanceName = instanceName
}

func (t *InstanceRequestSpec) GetNodePoolName() string {
	return t.NodePoolName
}

func (t *InstanceRequestSpec) SetNodePoolName(nodePoolName string) {
	t.NodePoolName = nodePoolName
}

func (t *InstanceRequestSpec) GetNamespace() string {
	return t.Namespace
}

func (t *InstanceRequestSpec) SetNamespace(namespace string) {
	t.Namespace = namespace
}

func (t *InstanceRequestSpec) GetFlavorName() string {
	return t.FlavorName
}

func (t *InstanceRequestSpec) SetFlavorName(flavorName string) {
	t.FlavorName = flavorName
}

func (t *InstanceRequestSpec) GetDescription() string {
	return t.Description
}

func (t *InstanceRequestSpec) SetDescription(description string) {
	t.Description = description
}

func (t *InstanceRequestSpec) UseLinked() bool {
	return t.UseLinkedRepo
}

//SpecsFromFile method return instance form string
func (t InstanceRequestSpec) SpecsFromFile(fileName string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromFile(fileName, new(InstanceRequestSpec), f...)
}

func (t InstanceRequestSpec) SpecsFromString(s string, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpecFromFromString(s, new(InstanceRequestSpec), f...)
}

// SpecsFromReader create spec from reader
func (t InstanceRequestSpec) SpecsFromReader(r io.Reader, f ...SpecFormatType) (*RequestSpec, error) {
	return ReadSpec(r, new(InstanceRequestSpec), f...)
}

//Validate method validate specs
func (t *InstanceRequestSpec) Validate() error {

	if t == nil {
		return &InvalidInstanceSpec{errMsg: "nil instance"}
	}

	if t.Kind() != SpecKindInstance {
		return &InvalidInstanceSpec{errMsg: "spec must contain kind field"}
	}

	if len(t.Repo) > 0 {
		if len(t.RepoUsername) == 0 {
			return &InvalidInstanceSpec{errMsg: "for repo overwrite you must provide username."}
		}
		if len(t.RepoPassword) == 0 {
			return &InvalidInstanceSpec{errMsg: "for repo overwrite you must provide username."}
		}
	}

	_, err := govalidator.ValidateStruct(t)
	if err != nil {
		return err
	}

	return nil
}

// IsValid return false if validator set error
func (t *InstanceRequestSpec) IsValid() bool {
	if t.specError != nil {
		return false
	}
	return true
}
