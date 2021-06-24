package request

// CreateVnfLcm Vnf Lcm Action

const (
	LcmTypeScaleOut = "SCALE_OUT"
	AspectId        = "aspect1"
)

type CreateVnfLcm struct {
	VnfdId                 string `json:"vnfdId" yaml:"vnfdId"`
	VnfInstanceName        string `json:"vnfInstanceName" yaml:"vnfInstanceName"`
	VnfInstanceDescription string `json:"vnfInstanceDescription" yaml:"vnfInstanceDescription"`
	Metadata               struct {
		Tags []interface{} `json:"tags" yaml:"tags"`
	} `json:"metadata" yaml:"metadata"`
}

type VduParams struct {
	VduName             string `json:"vduName" yaml:"vdu_name"`
	DeploymentProfileId string `json:"deploymentProfileId" yaml:"deployment_profile_id"`
	ChartName           string `json:"chartName" yaml:"chart_name"`
	HelmName            string `json:"helmName" yaml:"helm_name"`
	Namespace           string `json:"namespace" yaml:"namespace"`
	RepoUrl             string `json:"repoUrl" yaml:"repo_url"`
	Username            string `json:"username" yaml:"username"`
	Password            string `json:"password" yaml:"password"`
	Overrides           string `json:"overrides" yaml:"overrides"`
	Metadata            []struct {
	} `json:"metadata" yaml:"metadata"`
	ImageName           string `json:"imageName" yaml:"image_name"`
	DisableAutoRollback bool   `json:"disableAutoRollback" yaml:"disable_auto_rollback"`
	DisableGrant        bool   `json:"disableGrant" yaml:"disable_grant"`
	IgnoreGrantFailure  bool   `json:"ignoreGrantFailure" yaml:"ignore_grant_failure"`
	CatalogName         string `json:"catalogName" yaml:"catalog_name"`
	CatalogId           string `json:"catalogId" yaml:"catalog_id"`
}

type CnfReconfigure struct {
	Type             string `json:"type"`
	AspectId         string `json:"aspectId"`
	NumberOfSteps    int    `json:"numberOfSteps"`
	AdditionalParams struct {
		SkipGrant     bool        `json:"skipGrant"`
		VduParams     []VduParams `json:"vduParams"`
		LcmInterfaces []struct {
			InterfaceName string `json:"interfaceName"`
			Parameters    []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"parameters"`
		} `json:"lcmInterfaces"`
	} `json:"additionalParams"`
}
