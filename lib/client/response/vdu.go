package response

import (
	"encoding/json"
	"reflect"
	"strings"
)

// InfraRequirements csar section
type InfraRequirements struct {
	NodeComponents struct {
		Kernel struct {
			KernelType struct {
				Name    string `json:"name" yaml:"name"`
				Version string `json:"version" yaml:"version"`
			} `json:"kernel_type" yaml:"kernel_type"`
		} `json:"kernel" yaml:"kernel"`
	}
}

type ToscaProperties struct {
	// DescriptorId main id to identify vdu
	DescriptorId       string             `json:"descriptor_id" yaml:"descriptor_id"`
	Provider           string             `json:"provider" yaml:"provider"`
	ProductName        string             `json:"product_name" yaml:"product_name"`
	Version            string             `json:"version" yaml:"version"`
	Id                 string             `json:"id" yaml:"id"`
	SoftwareVersion    string             `json:"software_version" yaml:"software_version"`
	DescriptorVersion  string             `json:"descriptor_version" yaml:"descriptor_version"`
	FlavourId          string             `json:"flavour_id" yaml:"flavour_id"`
	FlavourDescription string             `json:"flavour_description" yaml:"flavour_description"`
	VnfmInfo           []string           `json:"vnfm_info" yaml:"vnfm_info"`
	InfraRequirements  *InfraRequirements `json:"infra_requirements,omitempty" yaml:"infra_requirements,omitempty"`
}

// GetField - return struct field value
func (t *ToscaProperties) GetField(field string) string {

	r := reflect.ValueOf(t)
	fields, _ := t.GetFields()
	if _, ok := fields[field]; ok {
		f := reflect.Indirect(r).FieldByName(strings.Title(field))
		return f.String()
	}

	return ""
}

// GetFields return VduPackage fields name as
// map[string], each key is field name
func (t *ToscaProperties) GetFields() (map[string]interface{}, error) {

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

// VduPackage - vdu package
type VduPackage struct {
	Description string `json:"description" yaml:"description"`
	Vnf         struct {
		Properties   *ToscaProperties `json:"properties" yaml:"properties"`
		Requirements struct {
		} `json:"requirements" yaml:"requirements"`
		Capabilities struct {
		} `json:"capabilities" yaml:"capabilities"`
		Metadata struct {
		} `json:"metadata" yaml:"metadata"`
		Interfaces struct {
			Vnflcm struct {
			} `json:"Vnflcm" yaml:"vnflcm"`
		} `json:"interfaces" yaml:"interfaces"`
		NodeType string `json:"nodeType" yaml:"node_type"`
		NodeSpec struct {
		} `json:"node_spec" yaml:"node_spec"`
	} `json:"vnf" yaml:"vnf"`
	Metadata struct {
	} `json:"metadata" yaml:"metadata"`
	Inputs struct {
	} `json:"inputs" yaml:"inputs"`
	Policies []struct {
		Type       string `json:"type" yaml:"type"`
		Properties struct {
			InterfaceName string `json:"interface_name" yaml:"interface_name"`
			InterfaceType string `json:"interface_type" yaml:"interface_type"`
			IsEnabled     bool   `json:"isEnabled" yaml:"is_enabled"`
		} `json:"properties" yaml:"properties"`
	} `json:"policies" yaml:"policies"`
	Groups struct {
	} `json:"groups" yaml:"groups"`
	Basepath       string        `json:"basepath" yaml:"basepath"`
	VolumeStorages []interface{} `json:"volume_storages" yaml:"volume_storages"`
	Vdus           []struct {
		VduId       string `json:"vdu_id" yaml:"vdu_id"`
		Type        string `json:"type" yaml:"type"`
		Description string `json:"description" yaml:"description"`
		Properties  struct {
			Name         string `json:"name" yaml:"name"`
			Description  string `json:"description" yaml:"description"`
			ChartName    string `json:"chartName" yaml:"chart_name"`
			ChartVersion string `json:"chartVersion" yaml:"chart_version"`
			HelmVersion  string `json:"helmVersion" yaml:"helm_version"`
		} `json:"properties" yaml:"properties"`
		VirtualStorages []interface{} `json:"virtual_storages" yaml:"virtual_storages"`
		Dependencies    []interface{} `json:"dependencies" yaml:"dependencies"`
		Vls             []interface{} `json:"vls" yaml:"vls"`
		Cps             []interface{} `json:"cps" yaml:"cps"`
		Artifacts       []interface{} `json:"artifacts" yaml:"artifacts"`
	} `json:"vdus" yaml:"vdus"`
	Vls        []interface{} `json:"vls" yaml:"vls"`
	Cps        []interface{} `json:"cps" yaml:"cps"`
	VnfExposed struct {
		ExternalCps []interface{} `json:"external_cps" yaml:"external_cps"`
		ForwardCps  []interface{} `json:"forward_cps" yaml:"forward_cps"`
	} `json:"vnf_exposed" yaml:"vnf_exposed"`
	Graph struct {
		MwcApp02 []interface{} `json:"mwc_app02" yaml:"mwc_app_02"`
	} `json:"graph" yaml:"graph"`
	VduDependencyDetails struct {
		Field1 []string `json:"0" yaml:"field_1"`
	} `json:"vduDependencyDetails" yaml:"vdu_dependency_details"`
}

// GetField - return struct field value
func (t *VduPackage) GetField(field string) string {

	r := reflect.ValueOf(t)
	fields, _ := t.GetFields()
	if _, ok := fields[field]; ok {
		f := reflect.Indirect(r).FieldByName(strings.Title(field))
		return f.String()
	}

	return ""
}

// GetFields return VduPackage fields name as
// map[string], each key is field name
func (t *VduPackage) GetFields() (map[string]interface{}, error) {

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
