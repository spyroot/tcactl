package models

type Properties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

type Policies struct {
	PolicyScale          PolicyScale          `yaml:"policy_scale"`
	PolicyWorkflow       PolicyWorkflow       `yaml:"policy_workflow"`
	PolicyReconfigure    PolicyReconfigure    `yaml:"policy_reconfigure"`
	PolicyUpdate         PolicyUpdate         `yaml:"policy_update"`
	PolicyUpgrade        PolicyUpgrade        `yaml:"policy_upgrade"`
	PolicyUpgradePackage PolicyUpgradePackage `yaml:"policy_upgrade_package"`
}

type PolicyWorkflowProperties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

type PolicyReconfigureProperties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

type PolicyUpgradePackageProperties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

type PolicyWorkflow struct {
	Type       string                   `yaml:"type"`
	Properties PolicyWorkflowProperties `yaml:"properties"`
}

type PolicyUpgrade struct {
	Type       string                  `yaml:"type"`
	Properties PolicyUpgradeProperties `yaml:"properties"`
}

type ToscaPolicies struct {
	Policies []Policies `yaml:"policies"`
}

type PolicyReconfigure struct {
	Type       string                      `yaml:"type"`
	Properties PolicyReconfigureProperties `yaml:"properties"`
}

type PolicyUpdate struct {
	Type       string                 `yaml:"type"`
	Properties PolicyUpdateProperties `yaml:"properties"`
}

type PolicyUpdateProperties struct {
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
	InterfaceName string `yaml:"interface_name"`
}

type PolicyUpgradeProperties struct {
	InterfaceName string `yaml:"interface_name"`
	InterfaceType string `yaml:"interface_type"`
	IsEnabled     bool   `yaml:"isEnabled"`
}

type PolicyUpgradePackage struct {
	Type       string                         `yaml:"type"`
	Properties PolicyUpgradePackageProperties `yaml:"properties"`
}

type PolicyScale struct {
	Properties Properties `yaml:"properties"`
	Type       string     `yaml:"type"`
}
