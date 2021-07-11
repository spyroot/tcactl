// Package app
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Mustafa mbayramo@vmware.com

package api

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spyroot/tcactl/lib/api_errors"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"strings"
	"time"
)

type VmTemplateFilterType string

const (

	//VmwareTemplateK8s default filter for k8s  templates
	VmwareTemplateK8s VmTemplateFilterType = "k8svm"

	// DefaultNamespace - default name space used for placement
	DefaultNamespace = "default"

	// DefaultFlavor default vdu flavor
	DefaultFlavor = "default"
)

// ClusterDeleteApiReq api request to delete cluster
type ClusterDeleteApiReq struct {
	//
	Cluster string
	//
	IsBlocking bool
	//
	IsVerbose bool
}

// ClusterCreateApiReq api request to create cluster
type ClusterCreateApiReq struct {
	//
	Spec *specs.SpecCluster
	//
	IsDryRun bool
	//
	IsBlocking bool
	//
	IsVerbose bool
	// isFixConflict resolve IP conflict, by check next IP
	IsFixConflict bool
}

// TcaApi - TCA Api interface
// Called need to use NewTcaApi to get instance before
type TcaApi struct {

	// rest client used to interact with tca
	rest *client.RestClient

	// specString validator.
}

// NewTcaApi - return instance for API.
//
func NewTcaApi(r *client.RestClient) (*TcaApi, error) {

	if r == nil {
		return nil, fmt.Errorf("rest client is nil, initilize rest client first")
	}

	a := &TcaApi{
		rest: r,
	}

	return a, nil
}

// CloudProviderNotFound error raised if tenant cloud not found
type CloudProviderNotFound struct {
	errMsg string
}

// Error - return if cloud provider not found
func (m *CloudProviderNotFound) Error() string {
	return m.errMsg + " cloud provider not found"
}

// UnsupportedCloudProvider error raised if tenant cloud not found
type UnsupportedCloudProvider struct {
	errMsg string
}

//
func (m *UnsupportedCloudProvider) Error() string {
	return m.errMsg + " cloud provider not supported"
}

func NewInstanceRequestSpec(cloudName string, clusterName string, vimType string, nfdName string,
	repo string, instanceName string, nodePoolName string) *specs.InstanceRequestSpec {
	i := &specs.InstanceRequestSpec{
		CloudName:        cloudName,
		ClusterName:      clusterName,
		VimType:          vimType,
		NfdName:          nfdName,
		Repo:             repo,
		InstanceName:     instanceName,
		NodePoolName:     nodePoolName,
		UseLinkedRepo:    true,
		AdditionalParams: specs.AdditionalParams{}}

	i.FlavorName = DefaultNamespace
	i.Description = ""
	i.Namespace = DefaultFlavor

	return i
}

// GetAllNodePool return a Node pool for particular cluster
// It generally useful to get list only if we need to display all
// in all other cases it efficient to use direct call for cluster.
func (a *TcaApi) GetAllNodePool(ctx context.Context) (*response.NodePool, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	clusters, err := a.rest.GetClusters(ctx)
	if err != nil {
		return nil, err
	}

	var pools response.NodePool
	// get all nodes pools
	for _, cluster := range clusters.Clusters {

		glog.Infof("Retrieving pool for a cluster name: '%v' uuid: '%v' state '%v'",
			cluster.ClusterName, cluster.Id, cluster.Status)

		// if cluster in failed state we have no pool.s
		if len(cluster.Id) == 0 {
			glog.Infof("SpecCluster id empty value")
			continue
		}

		clusterPool, poolErr := a.rest.GetClusterNodePools(cluster.Id)
		if poolErr != nil {
			glog.Error(err)
			return nil, err
		}

		if clusterPool != nil {
			glog.Infof("Got pool ids '%v'", clusterPool.GetIds())
			spec := clusterPool.Pools
			pools.Pools = append(pools.Pools, spec...)
		} else {
			glog.Error("Node pool is nil")
		}
	}

	return &pools, nil
}

// GetVimTenants method return vim tenants as response.Tenants
// collection.
func (a *TcaApi) GetVimTenants(ctx context.Context) (*response.Tenants, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetVimTenants(ctx)
}

// GetCurrentClusterTask get current cluster task
// taskId is operationId field.
func (a *TcaApi) GetCurrentClusterTask(ctx context.Context, taskId string) (*models.ClusterTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if IsValidUUID(taskId) == false {
		return nil, &InvalidTaskId{taskId}
	}

	clusters, err := a.rest.GetClusters(ctx)
	if err != nil {
		return nil, err
	}

	cid, err := clusters.GetClusterId(taskId)
	if err != nil {
		return nil, err
	}

	glog.Infof("Retrieving current task task list for cluster '%v'", cid)
	task, err := a.rest.GetClustersTask(ctx, specs.NewClusterTaskQuery(taskId))
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetVimVMTemplates - return compute cluster attached to cloud provider.
// caller need indicate template type and version.
func (a *TcaApi) GetVimVMTemplates(ctx context.Context, cloudName string,
	templateType VmTemplateFilterType, ver string) (*models.VcInventory, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(cloudName) == 0 {
		return nil, fmt.Errorf("empty cloud provider name")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	glog.Infof("Retrieving vm template list for cloud provider "+
		"uuid %s url %s", tenant.HcxUUID, tenant.VimURL)

	if len(tenant.HcxUUID) == 0 {
		return nil, fmt.Errorf("cloud provider hcx uuid is empty string")
	}

	if !tenant.IsVMware() {
		glog.Errorf("unsupported vim")
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	_filter := specs.NewVMwareTemplateQuery(tenant.HcxUUID, string(templateType), ver)
	t, err := a.rest.GetVMwareTemplates(ctx, _filter)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// GetVimFolders - return folders in target cloud provider.
// for VMware VC it list of VM Folders, models.Folders
func (a *TcaApi) GetVimFolders(ctx context.Context, cloudName string) (*models.Folders, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(cloudName) == 0 {
		return nil, fmt.Errorf("empty cloud provider name")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	glog.Infof("Retrieving vim folders list for cloud provider '%v' , '%v'",
		tenant.HcxUUID, tenant.VimURL)

	if len(tenant.HcxUUID) == 0 {
		return nil, fmt.Errorf("cloud provider hcx uuid is empty string")
	}

	if !tenant.IsVMware() {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	f := specs.NewVmwareFolderQuery(tenant.HcxUUID)
	if f == nil {
		return nil, fmt.Errorf("failed create folder filter")
	}

	t, err := a.rest.GetVMwareFolders(f)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// GetVimResourcePool - return resource pool in target VIM.
// caller need indicate template type and version.
func (a *TcaApi) GetVimResourcePool(ctx context.Context, cloudName string) (*models.ResourcePool, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	glog.Infof("Retrieving network list for cloud provider %v %v", tenant.HcxUUID, tenant.VimURL)
	if len(tenant.HcxUUID) == 0 {
		return nil, fmt.Errorf("cloud provider is empty")
	}

	if !tenant.IsVMware() {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	f := specs.NewVMwareResourcePoolQuery(tenant.HcxUUID)
	if f == nil {
		return nil, fmt.Errorf("failed create folder filter")
	}

	t, err := a.rest.GetVMwareResourcePool(ctx, f)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// validateKubeSpecCheck validate worker and master sub section
// of specs
func validateKubeSpecCheck(spec []models.TypeNode) error {

	for _, masterSpec := range spec {
		if len(masterSpec.Name) == 0 {
			return fmt.Errorf("cluster initialSpec master and worker must contains name")
		}
		if len(masterSpec.PlacementParams) == 0 {
			return fmt.Errorf("cluster initialSpec master and worker " +
				"must contains PlacementParams section and each type/name key value")
		}
		if len(masterSpec.Networks) == 0 {
			return fmt.Errorf("cluster initialSpec master " +
				"and worker must contains Networks section")
		}
		for _, network := range masterSpec.Networks {
			if len(network.NetworkName) == 0 {
				return fmt.Errorf("cluster initialSpec master and worker " +
					"each network must contains network name name, key (networkName)")
			}
			if len(network.Label) == 0 {
				return fmt.Errorf("masterNodes each network " +
					"must contains network label name, key (label)")
			}
		}
	}

	return nil
}

func validateClusterSpec(spec *specs.SpecCluster) error {

	if len(spec.Name) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains name key:value")
	}
	if len(spec.ClusterPassword) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains cluster password key:value")
	}
	if len(spec.ClusterTemplateId) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains ClusterTemplateId key:value")
	}
	if len(spec.ClusterType) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains clusterType and value management or workload key:value")
	}
	if len(spec.HcxCloudUrl) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains hcxCloudUrl key:value")
	}
	if len(spec.EndpointIP) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains EndpointIP key:value")
	}
	if len(spec.VmTemplate) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains vmTemplate key:value")
	}
	if len(spec.VmTemplate) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains vmTemplate key:value")
	}
	if len(spec.MasterNodes) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains masterNodes section")
	}
	if len(spec.WorkerNodes) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains WorkerNodes section")
	}
	if len(spec.PlacementParams) == 0 {
		return api_errors.NewInvalidSpec("cluster initialSpec must contains PlacementParams section")
	}
	if err := validateKubeSpecCheck(spec.MasterNodes); err != nil {
		return err
	}
	if err := validateKubeSpecCheck(spec.WorkerNodes); err != nil {
		return err
	}

	if spec.ClusterType == string(specs.ClusterWorkload) {
		if len(spec.ManagementClusterId) == 0 {
			return fmt.Errorf("workload cluster must contain ManagementClusterId key value")
		}
	}

	return nil
}

// NormalizeTemplateId - resolve template name to id
// for a give template type. Both must match in order
// method return ture
func (a *TcaApi) NormalizeTemplateId(IdOrName string, templateType string) (string, error) {

	if len(IdOrName) == 0 {
		return "", api_errors.NewTemplateNotFound(IdOrName)
	}

	if len(templateType) == 0 {
		return "", api_errors.NewTemplateInvalidType("template type is empty")
	}

	glog.Infof("Resolving cluster template '%s' to id", IdOrName)

	// resolve template id, in case client used name instead id
	clusterTemplates, err := a.rest.GetClusterTemplates()
	if err != nil {
		return "", err
	}

	t, err := clusterTemplates.GetTemplate(IdOrName)
	if err != nil {
		return "", err
	}

	glog.Infof("Resolved cluster template '%s' to template id %s", IdOrName, t.Id)

	// check template type
	if strings.ToLower(t.ClusterType) != strings.ToLower(templateType) {
		return "", fmt.Errorf("found template %s but template type mistmatch, "+
			"template type %v %v", IdOrName, t.ClusterType, templateType)
	}

	glog.Infof("Resolved template to template id %v", t.Id)
	return t.Id, nil
}

// doCheckCloudEndpoint
func (a *TcaApi) validateCloudEndpoint(ctx context.Context, cloud string) (*response.TenantsDetails, error) {

	// resolve template id, in case client used name instead id
	vimTenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(cloud, "https://") {
		tenant, err := vimTenants.FindCloudLink(cloud)
		if err != nil {
			return nil, err
		}

		return tenant, nil
	}

	t, err := vimTenants.FindCloudProvider(cloud)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// Validate cloud tenant state
func (a *TcaApi) validateTenant(tenant *response.TenantsDetails) error {

	glog.Infof("Validating cloud and tenant state.")

	if tenant.VimConn.Status != "ok" {
		return fmt.Errorf("cloud provider currently disconected")
	}

	if strings.ToLower(tenant.VimType) == models.VimTypeKubernetes {
		return fmt.Errorf("cloud provider already set to kubernetes")
	}

	return nil
}

// validatePlacement validate placement cluster placement attributes
func (a *TcaApi) validatePlacement(vim *VmwareVim,
	param models.PlacementParams) error {

	if models.VmwareTypeFolder(param.Type) == models.TypeFolder {
		glog.Infof("Validating vm folder '%s'", param.Name)
		if vim.folders.IsValidFolder(param.Name) == false {
			return api_errors.NewInvalidSpec("failed find a target folder")
		} else {
			glog.Infof("Resolved remote datastore folder.")
		}
	}
	if models.VmwareDatastore(param.Type) == models.TypeDataStore {
		glog.Infof("Validating datastore '%s'", param.Name)
		if vim.clusters.IsValidDatastore(param.Name) == false {
			return api_errors.NewInvalidSpec("failed find a target datastore")
		} else {
			glog.Infof("Resolved remote datastore name.")
		}
	}
	if models.VmwareResourcePool(param.Type) == models.TypeResourcePool {
		glog.Infof("Validating resource pool '%s'", param.Name)
		if vim.resourcePool.IsValidResource(param.Name) == false {
			return api_errors.NewInvalidSpec("failed find a target resource pool")
		} else {
			glog.Infof("Resolved remote resource pool.")
		}
	}
	if models.ClusterComputeResource(param.Type) == models.TypeClusterComputeResource {
		glog.Infof("Validating compute cluster '%s'", param.Name)
		if vim.clusters.IsValidClusterCompute(param.Name) == false {
			return api_errors.NewInvalidSpec("failed find a cluster compute resource")
		} else {
			glog.Infof("Resolved remote cluster name.")
		}
	}

	return nil
}

// normalizeDatastoreName adjust store name
func normalizeDatastoreName(ds string) string {

	fixedUrl := ""
	if ds == "ds:///" {
		fixedUrl = ds
	} else {
		if strings.HasPrefix(ds, "/vmfs") {
			fixedUrl = "ds:///" + ds
		}
		if strings.HasPrefix(ds, "vmfs") {
			fixedUrl = "ds:///" + ds
		}
	}

	// fixed '/' at the end
	if len(fixedUrl) > 0 && !strings.HasSuffix(fixedUrl, "/") {
		fixedUrl = fixedUrl + "/"
	}

	return fixedUrl
}

type VmwareVim struct {
	clusters     *models.VMwareClusters
	folders      *models.Folders
	resourcePool *models.ResourcePool
	networks     *models.CloudNetworks
}

func (a *TcaApi) getVmwareVimState(ctx context.Context, vimName string) (*VmwareVim, error) {
	// vc compute

	var (
		vimState VmwareVim
		err      error
	)

	vimState.clusters, err = a.GetVimComputeClusters(ctx, vimName)
	if err != nil {
		return nil, err
	}
	// vc folders
	vimState.folders, err = a.GetVimFolders(ctx, vimName)
	if err != nil {
		return nil, err
	}
	// vc resource pools
	vimState.resourcePool, err = a.GetVimResourcePool(ctx, vimName)
	if err != nil {
		return nil, err
	}

	vimState.networks, err = a.GetVimNetworks(ctx, vimName)
	if err != nil {
		return nil, err
	}

	return &vimState, nil
}

func (a *TcaApi) validateExtensions(ctx context.Context, spec *specs.SpecCluster) error {

	repos, err := a.GetRepos(ctx)
	if err != nil {
		return err
	}

	if spec.ClusterConfig != nil {
		for _, tool := range spec.ClusterConfig.Tools {
			if tool.Name == "harbor" {
				// adjust
				if tool.Properties.Type != "extension" {
					tool.Properties.Type = "extension"
				}
				extId, err := repos.GetRepoId(tool.Properties.ExtensionId)
				if err != nil {
					return err
				}

				tool.Properties.ExtensionId = extId
			}
		}
	}

	return nil
}

// validateVmwarePlacement method validate placement for VMware VIM
func (a *TcaApi) validateVmwarePlacement(ctx context.Context, spec *specs.SpecCluster, tenant *response.TenantsDetails) error {

	vmwareVim, err := a.getVmwareVimState(ctx, tenant.VimName)
	if err != nil {
		return err
	}

	// check for VMware vsphere-csi
	if spec.ClusterConfig != nil {
		for _, s := range spec.ClusterConfig.Csi {
			if s.Name == ClusterCsiVsphere {
				if s.Properties == nil {
					return api_errors.NewInvalidSpec("Invalid vsphere-csi property")
				}
				if !vmwareVim.clusters.IsValidDatastoreUrl(s.Properties.DatastoreUrl) {
					return api_errors.NewInvalidSpec("Invalid vsphere-csi property")
				}

				ds, err := vmwareVim.clusters.GetDatastoreByUrl(s.Properties.DatastoreUrl)
				if err != nil {
					return api_errors.NewInvalidSpec(err.Error())
				}
				s.Properties.DatastoreName = ds.Name
			}
		}
	}

	// validate global specString
	glog.Infof("Validating global specString placement parameters")
	for _, param := range spec.PlacementParams {
		if err := a.validatePlacement(vmwareVim, param); err != nil {
			return api_errors.NewInvalidSpec(err.Error())
		}
	}

	glog.Infof("Validating master node specString placement parameters")
	for i, worker := range spec.WorkerNodes {
		for j, n := range worker.Networks {
			// normalize port-group name
			normalized, err := vmwareVim.networks.NormalizeName(n.NetworkName)
			if err != nil {
				return api_errors.NewInvalidSpec(err.Error())
			}
			spec.WorkerNodes[i].Networks[j].NetworkName = normalized
		}
		for _, param := range worker.PlacementParams {
			if err := a.validatePlacement(vmwareVim, param); err != nil {
				return api_errors.NewInvalidSpec(err.Error())
			}
		}
	}

	glog.Infof("Validating worker node specString placement parameters")
	for i, master := range spec.MasterNodes {
		for j, n := range master.Networks {
			// normalize port-group name
			normalized, err := vmwareVim.networks.NormalizeName(n.NetworkName)
			if err != nil {
				return api_errors.NewInvalidSpec(err.Error())
			}
			spec.MasterNodes[i].Networks[j].NetworkName = normalized
		}

		for _, param := range master.PlacementParams {
			if err := a.validatePlacement(vmwareVim, param); err != nil {
				return api_errors.NewInvalidSpec(err.Error())
			}
		}
	}

	return nil
}

// validateCsi validate csi specString
func (a *TcaApi) validateCsi(spec *specs.SpecCluster) error {
	if spec.ClusterConfig != nil {
		for _, s := range spec.ClusterConfig.Csi {
			if s.Name == CLusterCsiNfs {
				if len(s.Properties.MountPath) == 0 {
					return api_errors.NewInvalidSpec("Invalid nfs client mount path")
				}
				if len(s.Properties.ServerIP) == 0 {
					return api_errors.NewInvalidSpec("Invalid nfs client server address")
				}
			}
		}
	}

	return nil
}

// validateCsi validate csi specString
func (a *TcaApi) validateVim(spec *specs.SpecCluster, tenant *response.TenantsDetails) error {

	if tenant == nil {
		return fmt.Errorf("tenant vim is nil")
	}

	if len(tenant.VimType) == 0 {
		return fmt.Errorf("unknown vim type")
	}

	if tenant.VimConn.Status != "ok" {
		return fmt.Errorf("cloud provider currently disconected")
	}

	if spec.Name == tenant.Name {
		spec.HcxCloudUrl = tenant.HcxCloudURL
	} else {
		return api_errors.NewInvalidSpec("Invalid cloud url, name")
	}

	return nil
}

// Validate cloud tenant state
func (a *TcaApi) validatePlacements(ctx context.Context, spec *specs.SpecCluster, tenant *response.TenantsDetails) error {

	glog.Infof("Validate placement details.")

	//if err := a.validateVim(specString, tenant); err != nil {
	//	return err
	//}

	if err := a.validateCsi(spec); err != nil {
		return err
	}

	err := a.validateExtensions(ctx, spec)
	if err != nil {
		return err
	}

	if tenant.IsVMware() {
		glog.Infof("Target cloud provider is VMware cluster, validating vc placements.")
		err := a.validateVmwarePlacement(ctx, spec, tenant)
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteTenantCluster - deletes tenant cluster
// it accept just name or id.
func (a *TcaApi) DeleteTenantCluster(ctx context.Context, tenantCluster string) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	vims, err := a.GetVims(ctx)
	if err != nil {
		return nil, err
	}

	clouds, err := vims.GetTenantClouds(tenantCluster, models.VimTypeKubernetes)
	if err != nil {
		return nil, err
	}

	return a.rest.DeleteTenant(clouds.TenantID)
}

// ResolveVim resolve vim name to id
func (a *TcaApi) ResolveVim(ctx context.Context, name string) (*response.TenantsDetails, error) {

	vims, err := a.GetVims(ctx)
	if err != nil {
		return nil, err
	}

	provider, err := vims.FindCloudProvider(name)
	if err != nil {
		return nil, err
	}

	if provider == nil {
		return nil, fmt.Errorf("nil vim")
	}

	return provider, nil
}

// ResolveVimName resolve name to id
func (a *TcaApi) ResolveVimName(ctx context.Context, name string) (string, error) {

	vim, err := a.ResolveVim(ctx, name)
	if err != nil {
		return "", err
	}

	return vim.ID, nil
}

func (a *TcaApi) ResolveVimId(ctx context.Context, name string) (string, error) {

	vim, err := a.ResolveVim(ctx, name)
	if err != nil {
		return "", err
	}

	return vim.VimID, nil
}

func (a *TcaApi) CreateSpecExample(
	name string,
	t string,
	h string,
	vm string,
	ip string,
	p string,
	netpath string,
	dns string,
) *specs.SpecCluster {

	var spec specs.SpecCluster
	spec.Name = name
	spec.ClusterType = string(specs.ClusterManagement)
	spec.ClusterTemplateId = t
	spec.HcxCloudUrl = h
	spec.VmTemplate = vm
	spec.EndpointIP = ip

	spec.PlacementParams = []models.PlacementParams{
		*models.NewPlacementParams("templates", "Folder"),
		*models.NewPlacementParams("vsanDatastore", "Datastore"),
		*models.NewPlacementParams("pod03", "ResourcePool"),
		*models.NewPlacementParams("core", "IsValidClusterCompute"),
	}
	spec.ClusterPassword = p

	net := models.NewNetworks(string(specs.ClusterManagement),
		netpath,
		[]string{dns})

	master := models.NewTypeNode("master", []models.Networks{*net}, []models.PlacementParams{
		*models.NewPlacementParams("Discovered virtual machine", "Folder"),
		*models.NewPlacementParams("vsanDatastore", "Datastore"),
		*models.NewPlacementParams("pod03", "ResourcePool"),
		*models.NewPlacementParams("core", "IsValidClusterCompute"),
	})

	worker := models.NewTypeNode("note-pool01", []models.Networks{*net}, []models.PlacementParams{
		*models.NewPlacementParams("Discovered virtual machine", "Folder"),
		*models.NewPlacementParams("vsanDatastore", "Datastore"),
		*models.NewPlacementParams("pod", "ResourcePool"),
		*models.NewPlacementParams("core", "IsValidClusterCompute"),
	})

	spec.MasterNodes = []models.TypeNode{*master}
	spec.WorkerNodes = []models.TypeNode{*worker}

	return &spec
}

func (a *TcaApi) GetClusters(ctx context.Context) (*response.Clusters, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	glog.Infof("Retrieving cluster list.")
	return a.rest.GetClusters(ctx)
}

func (a *TcaApi) GetClusterNodePools(Id string) (*response.NodePool, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	glog.Infof("Retrieving cluster list.")

	return a.rest.GetClusterNodePools(Id)
}

// GetAllNodePools - return all node pool for clusterId
func (a *TcaApi) GetAllNodePools(ctx context.Context) ([]response.NodesSpecs, error) {

	glog.Infof("Retrieving node pools.")

	var allSpecs []response.NodesSpecs
	if a.rest == nil {
		return allSpecs, fmt.Errorf("rest interface is nil")
	}

	clusters, err := a.GetClusters(ctx)
	if err != nil {
		return allSpecs, err
	}

	for _, c := range clusters.Clusters {
		pools, err := a.GetClusterNodePools(c.Id)
		if err != nil {
			return allSpecs, err
		}
		for _, p := range pools.Pools {
			r, err := a.rest.GetClusterNodePool(c.Id, p.Id)
			if err != nil {
				return allSpecs, err
			}
			allSpecs = append(allSpecs, *r)
		}
	}

	return allSpecs, nil
}

// GetTenant method return tenant as response.Tenants
// if tenant is name, method will lookup by name.
// if tenant is UUID it will lookup by id
// if it has prefix vmware it will lookup by VIM id.
func (a *TcaApi) GetTenant(ctx context.Context, tenant string) (*response.Tenants, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(tenant) == 0 {
		return nil, fmt.Errorf("empty tenant")
	}

	tenants, err := a.rest.GetVimTenants(ctx)
	if err != nil {
		return nil, err
	}

	r := response.Tenants{}

	if strings.Contains(tenant, "vmware") {
		r.TenantsList = tenants.Filter(response.FilterVimId, func(s string) bool {
			return s == tenant
		})
	} else if IsValidUUID(tenant) {
		r.TenantsList = tenants.Filter(response.FilterId, func(s string) bool {
			return s == tenant
		})
	} else {
		r.TenantsList = tenants.Filter(response.FilterName, func(s string) bool {
			return s == tenant
		})
	}

	return &r, err
}

// GetTenantsQuery query tenant information based on
// tenant id and package id
func (a *TcaApi) GetTenantsQuery(tenantId string, nfType string) (*response.Tenants, error) {

	glog.Infof("Retrieving tenants.")

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	reqFilter := specs.TenantsNfFilter{
		Filter: specs.TenantFilter{
			NfType: nfType,
		},
	}

	// attach tenant id if need
	if len(tenantId) > 0 {
		reqFilter.Filter.NfdId = tenantId
		//cid, vnfdId, err := a.GetCatalogId(tenantId)
		//if err != nil {
		//	return nil, nil
		//}
		//glog.Infof("Acquired catalog id '%v', for vnfId '%v'", cid, vnfdId)
		//reqFilter.Filter.NfdId = vnfdId
	}

	return a.rest.GetTenantsQuery(&reqFilter)
}

// GetNamedClusterTemplate - return template
// if name is string first resolve template
func (a *TcaApi) GetNamedClusterTemplate(name string) (*response.ClusterTemplateSpec, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	var templateId = name
	// resolve name first
	if IsValidUUID(name) == false {
		templates_, err := a.GetClusterTemplates()
		if err != nil {
			return nil, err
		}
		templateId, err = templates_.GetTemplateId(name)
		if err != nil {
			return nil, err
		}
	}

	return a.rest.GetClusterTemplate(templateId)
}

// CreateCnfNewInstance create a new instance of VNF or CNF.
// Dry run will validate request but will not create any CNF.
func (a *TcaApi) CreateCnfNewInstance(ctx context.Context, n *specs.InstanceRequestSpec, isDry bool, isBlocked bool) (*response.LcmInfo, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if n == nil {
		return nil, fmt.Errorf("nil request")
	}

	var (
		repoUsername = n.RepoUsername
		repoPassword = n.RepoPassword
	)

	if n.IsAutoName() {
		instance, err := a.GetInstance(ctx, n.InstanceName)
		if err != nil {
			return nil, err
		}
		if instance.VnfInstanceName == n.InstanceName {
			n.InstanceName = n.InstanceName + "-" + uuid.New().String()
			n.InstanceName = n.InstanceName[0:16]
		}
	}

	tenants, err := a.GetVimTenants(ctx)
	if err != nil {
		glog.Errorf("Failed acquire cloud tenant, error: %v", err)
		return nil, err
	}

	cloud, err := tenants.GetTenantClouds(n.CloudName, n.VimType)
	if err != nil {
		glog.Errorf("Failed acquire cloud provider details, error: %v", err)
		return nil, err
	}

	pkg, vnfd, err := a.GetCatalogAndVdu(n.NfdName)
	if err != nil || vnfd == nil {
		glog.Errorf("Failed acquire VDU information for %v", n.NfdName)
		return nil, err
	}

	// get linked Repo, if caller provide Repo that is not
	// linked nothing to do.
	reposUuid, err := a.rest.LinkedRepositories(cloud.TenantID, n.Repo)
	if err != nil {
		glog.Errorf("Failed acquire linked %v "+
			"repository to cloud provider %v. Indicate a Repo "+
			"linked to cloud provider.", n.Repo, cloud.TenantID)
		return nil, err
	}

	ext, err := a.rest.GetExtensions(ctx)
	if err != nil {
		glog.Errorf("Failed acquire extension information for %v", err)
		return nil, err
	}

	// if linked Repo resolve it
	if n.UseLinked() {
		linkedRepos, err := ext.FindRepo(reposUuid)
		if err != nil || linkedRepos == nil {
			glog.Errorf("Failed acquire extension information for %v", reposUuid)
			return nil, err
		}

		if linkedRepos.IsEnabled() == false {
			glog.Errorf("Repository %v is disabled", linkedRepos.Name)
			return nil, fmt.Errorf("repository %v is disabled", linkedRepos.Name)
		}

		glog.Infof("Found attached Repo %v and status %v", n.Repo, linkedRepos.State)

		// overwrite credential from linked Repo
		repoUsername = linkedRepos.AccessInfo.Username
		repoPassword = linkedRepos.AccessInfo.Password
	}

	// resolve nodePools
	nodePool, _, err := a.rest.GetNamedClusterNodePools(ctx, n.ClusterName)
	if err != nil || nodePool == nil {
		glog.Errorf("Failed acquire clusters node information for cluster %v, error %v", n.ClusterName, err)
		return nil, err
	}
	pool, err := nodePool.GetPool(n.NodePoolName)
	if err != nil {
		glog.Errorf("Failed acquire node pool information for node pool %v, error %v", n.NodePoolName, err)
		return nil, err
	}

	if isDry == true {
		return nil, nil
	}

	vnfLcm, err := a.rest.CreateInstance(ctx, &specs.LcmCreateRequest{
		VnfdId:                 pkg.VnfdID,
		VnfInstanceName:        n.InstanceName,
		VnfInstanceDescription: n.Description,
	})

	if err != nil {
		glog.Errorf("Failed create instance information %v", err)
		return nil, err
	}

	var flavorName = n.FlavorName
	if len(vnfd.Vnf.Properties.FlavourId) > 0 {
		flavorName = vnfd.Vnf.Properties.FlavourId
	}

	for _, vdu := range vnfd.Vdus {
		var req = specs.LcmInstantiateRequest{
			FlavourID: flavorName,
			VimConnectionInfo: []models.VimConnectionInfo{
				{
					Id:      cloud.VimID,
					VimType: "",
					Extra: &models.VimExtra{
						NodePoolId: pool.Id,
					},
				},
			},
			AdditionalVduParams: &specs.AdditionalParams{
				VduParams: []specs.VduParam{{
					Namespace: n.Namespace,
					RepoURL:   n.Repo,
					Username:  repoUsername,
					Password:  repoPassword,
					VduName:   vdu.VduId,
				}},
				DisableGrant:        n.AdditionalParams.DisableGrant,
				IgnoreGrantFailure:  n.AdditionalParams.IgnoreGrantFailure,
				DisableAutoRollback: n.AdditionalParams.DisableAutoRollback,
			},
		}

		glog.Infof("Instantiating %v", vnfLcm.Id)
		err := a.rest.InstanceInstantiate(ctx, vnfLcm.Id, req)
		if err != nil {
			glog.Errorf("Failed create cnf instance information %v", err)
			return nil, err
		}

	}

	instance, err := a.rest.GetRunningVnflcm(vnfLcm.Id)
	if err != nil {
		glog.Errorf("Failed create cnf instance information %v", err)
		return nil, err
	}

	if isBlocked {
		err := a.BlockWaitStateChange(ctx, vnfLcm.Id, StateInstantiate, DefaultMaxRetry, true)
		if err != nil {
			return instance, err
		}
	}
	return instance, nil
}

// BlockWaitStateChange - simple block and pull status
// instanceId is instance that method will pull and check
// waitFor is target status method waits.
// maxRetry a limit.
func (a *TcaApi) BlockWaitStateChange(ctx context.Context, instanceId string, waitFor string, maxRetry int, verbose bool) error {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for i := 1; i < maxRetry; i++ {
		{

			select {
			case <-ctx.Done():

				return nil
			default:

				instance, err := a.rest.GetRunningVnflcm(instanceId)
				if err != nil {
					return err
				}

				if verbose {
					fmt.Printf("Current state %s waiting for %s\n",
						instance.InstantiationState, waitFor)

					fmt.Printf("Current LCM Operation status %s target state %s\n\n",
						instance.Metadata.LcmOperationState,
						instance.Metadata.LcmOperation)
				}

				if instance.Metadata.LcmOperationState == "FAILED_TEMP" {
					break
				}

				if strings.HasPrefix(instance.InstantiationState, waitFor) {
					break
				}

				time.Sleep(TaskWaitSeconds * time.Second)
			}
		}
	}

	return nil
}

type TcaTaskFailed struct {
	ErrMsg string
}

func (e *TcaTaskFailed) Error() string {
	return "task failed on phase " + e.ErrMsg
}

type TaskCanceled struct {
	ErrMsg string
}

func (e *TaskCanceled) Error() string {
	return e.ErrMsg
}

type TaskNotFound struct {
	ErrMsg string
}

func (e *TaskNotFound) Error() string {
	return e.ErrMsg
}

// BlockWaitTaskFinish - simple block and pull status
// instanceId is instance that method will pull and check
// waitFor is target status method waits.
// maxRetry a limit.
func (a *TcaApi) BlockWaitTaskFinish(ctx context.Context, task *models.TcaTask, waitFor string, maxRetry int, verbose bool) error {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	req := specs.NewClusterTaskQuery(task.Id)

	if verbose {
		glog.Infof("Waiting task id=%s type=%s", task.Id, task.OperationId)
	}

	var (
		allTaskDone = false
		retry       = 0
	)

	for !allTaskDone || retry == maxRetry {
		{

			select {
			case <-ctx.Done():

				return nil
			default:

				task, err := a.rest.GetClustersTask(ctx, req)
				if err != nil {
					return err
				}

				if task == nil {
					return &TaskNotFound{"task not found"}
				}

				var numTaskWaiting = 0
				for _, item := range task.Items {
					if verbose {
						glog.Infof("Task id=%s type=%s progress=%d status=%s",
							item.TaskId, item.Type, item.Progress, item.Status)
					}

					if item.Status != waitFor {
						glog.Infof("Task still running id=%s type=%s progress=%d status=%s",
							item.TaskId, item.Type, item.Progress, item.Status)
						numTaskWaiting++
					}

					if item.Status == TaskStateFailed {
						glog.Errorf("Task failed.")
						return &TcaTaskFailed{item.Type}
					}
				}

				if numTaskWaiting == 0 {
					glog.Infof("All task finished.")
					allTaskDone = true
				}

				time.Sleep(TaskWaitSeconds * time.Second)
				retry++
			}
		}
	}

	return nil
}

// ResolvePoolId method resolves pool name to pool id
func (a *TcaApi) ResolvePoolId(ctx context.Context, poolName string, clusterScope ...string) (string, error) {

	// get cluster id for a pool
	if len(clusterScope) > 0 && len(clusterScope[0]) > 0 {
		var clusterId = clusterScope[0]
		if pool, restErr := a.rest.GetClusterNodePools(clusterId); restErr == nil && pool != nil {
			if spec, pollErr := pool.GetPool(poolName); pollErr == nil && spec != nil {
				return spec.Id, nil
			}
		}
	}

	// cluster scope not indicated, we get entire cluster.
	clusters, err := a.rest.GetClusters(ctx)
	if err != nil {
		glog.Error(err)
		return "", err
	}
	clusterIds, err := clusters.GetClusterIds()
	if err != nil {
		return "", err
	}

	// get all pools
	for _, cid := range clusterIds {
		if pool, restErr := a.rest.GetClusterNodePools(cid); restErr == nil && pool != nil {
			if nodeSpec, pollErr := pool.GetPool(poolName); pollErr == nil && nodeSpec != nil {
				return nodeSpec.Id, nil
			}
		}
	}

	return "", &response.PoolNotFound{ErrMsg: poolName}
}

// ResolvePoolName -  method resolves pool name to id  for a requested cluster.
// pool name is named pool and cluster is name or uuid
func (a *TcaApi) ResolvePoolName(ctx context.Context, poolName string, clusterName string) (string, string, error) {

	// empty name no ops
	if len(poolName) == 0 {
		return poolName, "", nil
	}

	if len(clusterName) == 0 {
		return "", "", fmt.Errorf("provide cluster name to resolve pool name")
	}

	nodePool, clusterId, err := a.rest.GetNamedClusterNodePools(ctx, clusterName)
	if err != nil || nodePool == nil {
		glog.Errorf("Failed acquire clusters node information %v", err)
		return poolName, "", err
	}

	pool, err := nodePool.GetPool(poolName)
	if err != nil {
		glog.Errorf("Failed acquire node pool information %v", err)
		return poolName, "", err
	}

	return pool.Id, clusterId, nil
}

// GetAuthorization retrieve API key from TCA
// and update internal state.
func (a *TcaApi) GetAuthorization() (bool, error) {
	return a.rest.GetAuthorization()
}

// ResolveClusterName - resolve cluster name to cluster id
// and return a string version. TCA use UUID format
// for ID.
func (a *TcaApi) ResolveClusterName(ctx context.Context, q string) (string, error) {

	clusters, err := a.rest.GetClusters(ctx)
	if err != nil {
		return "", nil
	}

	return clusters.GetClusterId(q)
}

// SetBaseUrl set tca base api url
func (a *TcaApi) SetBaseUrl(url string) {
	if a != nil && a.rest != nil {
		a.rest.BaseURL = url
	}
}

// GetBaseUrl set base url for TCA API client interface.
func (a *TcaApi) GetBaseUrl() string {

	if a != nil && a.rest != nil {
		return a.rest.BaseURL
	}

	return ""
}

// SetUsername - sets username for TCA API client interface.
func (a *TcaApi) SetUsername(username string) {

	if a != nil && a.rest != nil {
		a.rest.Username = username
	}
}

// SetPassword - set password for TCA API client interface.
func (a *TcaApi) SetPassword(password string) {

	if a != nil && a.rest != nil {
		a.rest.Password = password
	}

}

// GetApiKey returns API key used to connect to rest interface
func (a *TcaApi) GetApiKey() string {

	if a != nil && a.rest != nil {
		return a.rest.GetApiKey()
	}

	return ""
}

// SetTrace set trace that will output to stdout server responds
func (a *TcaApi) SetTrace(trace bool) {
	a.rest.SetTrace(trace)
}

// GetVims return all attached tenant vim
func (a *TcaApi) GetVims(ctx context.Context) (*response.Tenants, error) {
	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetVimTenants(ctx)
}
