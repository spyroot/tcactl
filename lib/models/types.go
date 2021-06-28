package models

func VimType() []string {
	return []string{
		VimTypeVmware,
		VimTypeKubernetes,
	}
}

const (
	// VimTypeVmware vim type vc ( vmware)
	VimTypeVmware = "vc"

	// VimTypeKubernetes vim type kubernetes
	VimTypeKubernetes = "kubernetes"

	// CniTypeAntrea
	CniTypeAntrea = "antrea"

	// CniTypeMultus
	CniTypeMultus = "multus"

	// CniTypeCalico - calico cni
	CniTypeCalico = "calico"

	// CniTypeWhereAbouts - cni CniTypeWhereAbouts
	CniTypeWhereAbouts = "whereabouts"

	CsiVmware = "vsphere-csi"

	CsiNfs = "nfs_client"

	TemplateMgmt string = "MANAGEMENT"

	TemplateWorkload string = "WORKLOAD"
)
