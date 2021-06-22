package models

// Properties
type Properties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

// Policies
type Policies struct {
	PolicyScale          PolicyScale          `yaml:"policy_scale"`
	PolicyWorkflow       PolicyWorkflow       `yaml:"policy_workflow"`
	PolicyReconfigure    PolicyReconfigure    `yaml:"policy_reconfigure"`
	PolicyUpdate         PolicyUpdate         `yaml:"policy_update"`
	PolicyUpgrade        PolicyUpgrade        `yaml:"policy_upgrade"`
	PolicyUpgradePackage PolicyUpgradePackage `yaml:"policy_upgrade_package"`
}

// PolicyWorkflowProperties
type PolicyWorkflowProperties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

// PolicyReconfigureProperties
type PolicyReconfigureProperties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

// PolicyUpgradePackageProperties
type PolicyUpgradePackageProperties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

// PolicyWorkflow
type PolicyWorkflow struct {
	Type       string                   `yaml:"type"`
	Properties PolicyWorkflowProperties `yaml:"properties"`
}

// PolicyUpgrade
type PolicyUpgrade struct {
	Type       string                  `yaml:"type"`
	Properties PolicyUpgradeProperties `yaml:"properties"`
}

type ToscaPolicies struct {
	Policies []Policies `yaml:"policies"`
}

// PolicyReconfigure
type PolicyReconfigure struct {
	Type       string                      `yaml:"type"`
	Properties PolicyReconfigureProperties `yaml:"properties"`
}

// PolicyUpdate
type PolicyUpdate struct {
	Type       string                 `yaml:"type"`
	Properties PolicyUpdateProperties `yaml:"properties"`
}

// PolicyUpdateProperties
type PolicyUpdateProperties struct {
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
	InterfaceName string `yaml:"interface_name"`
}

// PolicyUpgradeProperties
type PolicyUpgradeProperties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

// PolicyUpgradePackage
type PolicyUpgradePackage struct {
	Type       string                         `yaml:"type"`
	Properties PolicyUpgradePackageProperties `yaml:"properties"`
}

// PolicyScale
type PolicyScale struct {
	Properties Properties `yaml:"properties"`
	Type       string     `yaml:"type"`
}
