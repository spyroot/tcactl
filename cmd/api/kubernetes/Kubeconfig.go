package kubernetes

type KubeconfigStruct struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `json:"certificate-authority-data" yaml:"certificate-authority-data"`
			Server                   string `json:"server" yaml:"server"`
		} `json:"cluster" yaml:"cluster"`
		Name string `json:"name" yaml:"name"`
	} `json:"clusters" yaml:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `json:"cluster" yaml:"cluster"`
			User    string `json:"user" yaml:"user"`
		} `json:"context" yaml:"context"`
		Name string `json:"name" yaml:"name"`
	} `json:"contexts" yaml:"contexts"`
	CurrentContext string `json:"current-context" yaml:"current-context"`
	Kind           string `json:"kind" yaml:"kind"`
	Preferences    struct {
	} `json:"preferences" yaml:"preferences"`
	Users []struct {
		Name string `json:"name" yaml:"name"`
		User struct {
			ClientCertificateData string `json:"client-certificate-data" yaml:"client-certificate-data"`
			ClientKeyData         string `json:"client-key-data" yaml:"client-key-data"`
		} `json:"user" yaml:"user"`
	} `json:"users" yaml:"users"`
}
