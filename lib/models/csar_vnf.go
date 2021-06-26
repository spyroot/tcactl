package models

import (
	"reflect"
	"strings"
)

const (
	PropertyDescriptorId           = "descriptorId"
	PropertyProvider               = "provider"
	PropertyDescriptorVersion      = "descriptorVersion"
	PropertyFlavourId              = "flavourId,"
	PropertyFlavourDescription     = "flavourDescription"
	PropertyProductName            = "productName"
	PropertyVersion                = "version"
	PropertyId                     = "id"
	PropertySoftwareVersion        = "softwareVersion"
	PropertyChartName              = "chartName"
	PropertyChartVersion           = "chartVersion"
	PropertyHelmVersion            = "helmVersion"
	PropertyName                   = "name"
	PropertyDescription            = "description"
	PropertyConfigurableProperties = "configurableProperties"
	PropertyVnfmInfo               = "vnfm_info"
)

type CSAR struct {
	ToscaDefinitionsVersion string                `yaml:"tosca_definitions_version"`
	Description             string                `yaml:"description"`
	Imports                 []string              `yaml:"imports"`
	Node_Type               map[string]ToscaNodes `yaml:"node_types"`
	TopologyTemplate        TopologyTemplate      `yaml:"topology_template"`
}

// SubstitutionMappings - csar file contains a substitution sub-section
// each section model based based on node_type
type SubstitutionMappings struct {
	NodeType string `yaml:"node_type"`
}

type Vnflcm struct {
	Type string `yaml:"type,omitempty"`
}

// ToscaInterface -
type ToscaInterface struct {
	DerivedFrom string     `yaml:"derived_from,omitempty"`
	Interfaces  Interfaces `yaml:"interfaces,omitempty"`
}

// Interfaces - Vnflcm interfaces
type Interfaces struct {
	Vnflcm Vnflcm `yaml:"Vnflcm,omitempty"`
}

type ToscaNodes struct {
	DerivedFrom string     `yaml:"derived_from"`
	Interfaces  Interfaces `yaml:"interfaces"`
}

// TopologyTemplate - topology template section of csar
type TopologyTemplate struct {
	SubstitutionMappings SubstitutionMappings      `yaml:"substitution_mappings"`
	NodeTemplates        map[string]*NodeTemplates `yaml:"node_templates"`
}

// ToscaProperties Properties
type ToscaProperties struct {
	DescriptorId           string                 `yaml:"descriptor_id,omitempty"`
	Provider               string                 `yaml:"provider,omitempty"`
	DescriptorVersion      string                 `yaml:"descriptor_version,omitempty"`
	FlavourId              string                 `yaml:"flavour_id,omitempty"`
	FlavourDescription     string                 `yaml:"flavour_description,omitempty"`
	ProductName            string                 `yaml:"product_name,omitempty"`
	Version                string                 `yaml:"version,omitempty"`
	Id                     string                 `yaml:"id,omitempty"`
	SoftwareVersion        string                 `yaml:"software_version,omitempty"`
	ChartName              string                 `yaml:"chartName,omitempty"`
	ChartVersion           string                 `yaml:"chartVersion,omitempty"`
	HelmVersion            string                 `yaml:"helmVersion,omitempty"`
	Name                   string                 `yaml:"name,omitempty"`
	Description            string                 `yaml:"description,omitempty"`
	ConfigurableProperties ConfigurableProperties `yaml:"configurable_properties,omitempty"`
	VnfmInfo               []string               `yaml:"vnfm_info,omitempty"`
}

func (t *ToscaProperties) GetField(field string) reflect.Value {
	r := reflect.ValueOf(t)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	return f
}

func (t *ToscaProperties) UpdateField(field string, val string) {
	r := reflect.ValueOf(t)
	indirect := reflect.Indirect(r)
	if indirect.IsValid() {
		f := reflect.Indirect(r).FieldByName(strings.Title(field))
		f.Set(reflect.ValueOf(val))
	}
}

// NodeTemplates - node template section
type NodeTemplates struct {
	NodeType          string            `yaml:"node_type,omitempty"`
	Properties        ToscaProperties   `yaml:"properties"`
	Type              string            `yaml:"type,omitempty"`
	InfraRequirements InfraRequirements `yaml:"infra_requirements,omitempty"`
}

// InfraRequirements - infrastructure requirements
type InfraRequirements struct {
	NodeComponents NodeComponents `yaml:"node_components"`
}

// NodeComponents - kernel key
type NodeComponents struct {
	Kernel Kernel `yaml:"kernel"`
}

// Kernel kernel type
type Kernel struct {
	KernelType KernelType `yaml:"kernel_type"`
}

// KernelType - kernel type name and version
type KernelType struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

//
type Infra struct {
	InfraRequirements InfraRequirements `yaml:"infra_requirements"`
}

// AdditionalVnfcScalableProperties
type AdditionalVnfcScalableProperties struct {
}

// ConfigurableProperties ConfigurableProperties
type ConfigurableProperties struct {
	AdditionalVnfcConfigurableProperties AdditionalVnfcConfigurableProperties `yaml:"additional_vnfc_configurable_properties"`
	AdditionalVnfcScalableProperties     AdditionalVnfcScalableProperties     `yaml:"additional_vnfc_scalable_properties"`
}

// AdditionalVnfcConfigurableProperties =
type AdditionalVnfcConfigurableProperties struct {
}
