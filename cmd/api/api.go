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
	b64 "encoding/base64"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spyroot/hestia/cmd/client"
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/client/response"
	"github.com/spyroot/hestia/cmd/csar"
	"github.com/spyroot/hestia/cmd/models"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type VmTemplateFilterType string

const (
	VmwareTemplateK8s VmTemplateFilterType = "k8svm"
)

// TcaApi - TCA Api interface
type TcaApi struct {
	rest *client.RestClient
}

// NewTcaApi - return instance for API.
func NewTcaApi(rest *client.RestClient) (*TcaApi, error) {

	if rest == nil {
		return nil, fmt.Errorf("nil rest client argument")
	}

	return &TcaApi{rest: rest}, nil
}

// CloudProviderNotFound error raised if tenant cloud not found
type CloudProviderNotFound struct {
	errMsg string
}

//
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

// GetVimComputeClusters - return compute cluster attached to VIM
// For example VMware VIM
func (a *TcaApi) GetVimComputeClusters(cloudName string) (*models.VMwareClusters, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	if tenant.IsVMware() {
		//
		glog.Infof("Retrieving list for cloud provider %v '%v'",
			tenant.HcxUUID, tenant.VimURL)

		f := request.NewClusterFilterQuery(tenant.HcxUUID)
		clusterInventory, err := a.rest.GetVmwareCluster(f)

		if err != nil {
			return nil, err
		}

		return clusterInventory, nil

	} else {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	return nil, &CloudProviderNotFound{errMsg: cloudName}
}

// GetVimNetworks - return compute cluster attached to vim
func (a *TcaApi) GetVimNetworks(cloudName string) (*models.CloudNetworks, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	glog.Infof("Retrieving network list for cloud provider %v", tenant.HcxUUID, tenant.VimURL)

	if !tenant.IsVMware() {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	// get cluster id
	f := request.NewClusterFilterQuery(tenant.HcxUUID)
	clusterInventory, err := a.rest.GetVmwareCluster(f)
	if err != nil {
		return nil, err
	}

	var networks models.CloudNetworks

	// get all network for all clusters
	for _, item := range clusterInventory.Items {
		networkFilter := request.VMwareNetworkQuery{}
		networkFilter.Filter.TenantId = tenant.HcxUUID
		if strings.HasPrefix(item.EntityId, "domain") {
			networkFilter.Filter.ClusterId = item.EntityId
			net, err := a.rest.GetVmwareNetworks(&networkFilter)
			if err != nil {
				return nil, err
			}

			networks.Network = append(networks.Network, net.Network...)
		}
	}

	return &networks, nil
}

// GetNodePool return a Node pool for particular cluster
func (a *TcaApi) GetNodePool(clusterId string) (*response.NodePool, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	clusters, err := a.rest.GetClusters()
	if err != nil || clusters == nil {
		return nil, err
	}

	clusterId, err = clusters.GetClusterId(clusterId)
	if err != nil {
		return nil, err
	}

	pool, err := a.rest.GetClusterNodePools(clusterId)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// GetAllNodePool return a Node pool for particular cluster
// It generally useful to get list only if want to display
// in all other cases it efficient to use direct call for cluster.
func (a *TcaApi) GetAllNodePool() (*response.NodePool, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	clusters, err := a.rest.GetClusters()
	if err != nil {
		return nil, err
	}

	var pools response.NodePool
	// get all nodes pools
	for _, cluster := range clusters.Clusters {

		glog.Infof("Retrieving pool for a cluster name: '%v' uuid: '%v' state '%v'",
			cluster.ClusterName,
			cluster.Id,
			cluster.Status)

		// if cluster in failed state we have no pool.s
		if len(cluster.Id) == 0 {
			glog.Infof("Cluster id empty value")
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

// TenantsCloudProvider return a tenant attached to cloud provide for lookup query string
func (a *TcaApi) TenantsCloudProvider(query string) (*response.Tenants, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	r, err := tenants.FindCloudProvider(query)
	if err != nil {
		return nil, err
	}

	return &response.Tenants{
		TenantsList: []response.TenantsDetails{*r},
	}, nil
}

// GetVimTenants return vim tenants
func (a *TcaApi) GetVimTenants() (*response.Tenants, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetVimTenants()
}

// GetCurrentClusterTask get current cluster task
// taskId is operationId field.
func (a *TcaApi) GetCurrentClusterTask(taskId string) (*models.ClusterTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if IsValidUUID(taskId) == false {
		return nil, &InvalidTaskId{taskId}
	}
	clusters, err := a.rest.GetClusters()
	if err != nil {
		return nil, err
	}

	cid, err := clusters.GetClusterId(taskId)
	if err != nil {
		return nil, err
	}

	glog.Infof("Retrieving current task task list for cluster '%v'", cid)
	task, err := a.rest.GetClusterTask(request.NewClusterTaskQuery(taskId))
	if err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteCluster get current cluster task
func (a *TcaApi) DeleteCluster(clusterId string) (bool, error) {

	if a.rest == nil {
		return false, fmt.Errorf("rest interface is nil")
	}

	clusters, err := a.rest.GetClusters()
	if err != nil {
		return false, err
	}

	cid, clusterErr := clusters.GetClusterId(clusterId)
	if clusterErr != nil {
		glog.Error(clusterErr)
		return false, err
	}
	glog.Infof("Retrieving current task task list for cluster '%v'", cid)

	ok, err := a.rest.DeleteCluster(clusterId)
	if err != nil {
		return false, err
	}

	return ok, nil
}

// GetVimVMTemplates - return compute cluster attached to cloud provider.
// caller need indicate template type and version.
func (a *TcaApi) GetVimVMTemplates(cloudName string,
	templateType VmTemplateFilterType, ver string) (*models.VcInventory, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	tenant, err := tenants.FindCloudProvider(cloudName)
	if err != nil {
		return nil, err
	}

	glog.Infof("Retrieving network list for cloud provider %v", tenant.HcxUUID, tenant.VimURL)
	if len(tenant.HcxUUID) == 0 {
		return nil, fmt.Errorf("cloud provider is empty")
	}

	if !tenant.IsVMware() {
		return nil, &UnsupportedCloudProvider{errMsg: cloudName}
	}

	_filter := request.NewVMwareTemplateQuery(tenant.HcxUUID, string(templateType), ver)
	t, err := a.rest.GetVMwareTemplates(_filter)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// GetVimFolders - return folder in target VIM.
// caller need indicate template type and version.
func (a *TcaApi) GetVimFolders(cloudName string) (*models.Folders, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
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

	f := request.NewVmwareFolderQuery(tenant.HcxUUID)
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
func (a *TcaApi) GetVimResourcePool(cloudName string) (*models.ResourcePool, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
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

	f := request.NewVMwareResourcePoolQuery(tenant.HcxUUID)
	if f == nil {
		return nil, fmt.Errorf("failed create folder filter")
	}

	t, err := a.rest.GetVMwareResourcePool(f)
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
			return fmt.Errorf("cluster spec master and worker must contains name")
		}
		if len(masterSpec.PlacementParams) == 0 {
			return fmt.Errorf("cluster spec master and worker " +
				"must contains PlacementParams section and each type/name key value")
		}
		if len(masterSpec.Networks) == 0 {
			return fmt.Errorf("cluster spec master " +
				"and worker must contains Networks section")
		}
		for _, network := range masterSpec.Networks {
			if len(network.NetworkName) == 0 {
				return fmt.Errorf("cluster spec master and worker " +
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

func validateClusterSpec(spec *request.Cluster) error {

	if len(spec.Name) == 0 {
		return fmt.Errorf("cluster spec must contains name key:value")
	}
	if len(spec.ClusterPassword) == 0 {
		return fmt.Errorf("cluster spec must contains cluster password key:value")
	}
	if len(spec.ClusterTemplateId) == 0 {
		return fmt.Errorf("cluster spec must contains ClusterTemplateId key:value")
	}
	if len(spec.ClusterType) == 0 {
		return fmt.Errorf("cluster spec must contains clusterType and value management or workload key:value")
	}
	if len(spec.HcxCloudUrl) == 0 {
		return fmt.Errorf("cluster spec must contains hcxCloudUrl key:value")
	}
	if len(spec.EndpointIP) == 0 {
		return fmt.Errorf("cluster spec must contains EndpointIP key:value")
	}
	if len(spec.VmTemplate) == 0 {
		return fmt.Errorf("cluster spec must contains vmTemplate key:value")
	}
	if len(spec.VmTemplate) == 0 {
		return fmt.Errorf("cluster spec must contains vmTemplate key:value")
	}
	if len(spec.MasterNodes) == 0 {
		return fmt.Errorf("cluster spec must contains masterNodes section")
	}
	if len(spec.WorkerNodes) == 0 {
		return fmt.Errorf("cluster spec must contains WorkerNodes section")
	}
	if len(spec.PlacementParams) == 0 {
		return fmt.Errorf("cluster spec must contains PlacementParams section")
	}
	if err := validateKubeSpecCheck(spec.MasterNodes); err != nil {
		return err
	}
	if err := validateKubeSpecCheck(spec.WorkerNodes); err != nil {
		return err
	}

	if spec.ClusterType == string(request.ClusterWorkload) {
		if len(spec.ManagementClusterId) == 0 {
			return fmt.Errorf("workload cluster must contain ManagementClusterId key value")
		}
	}

	return nil
}

// ResolveTemplateId - resolve template name to id
// for a give template type. Both must match in order
// method return ture
func (a *TcaApi) ResolveTemplateId(templateId string, templateType string) (string, error) {

	// resolve template id, in case client used name instead id
	clusterTemplates, err := a.rest.GetClusterTemplates()
	if err != nil {
		return "", err
	}

	template, err := clusterTemplates.GetTemplate(templateId)
	if err != nil {
		return "", err
	}

	// check template type
	if template.ClusterType != templateType {
		return "", fmt.Errorf("found template by template type mistmatch")
	}

	return template.Id, nil
}

// doCheckCloudEndpoint
func (a *TcaApi) validateCloudEndpoint(cloudUrl string) (*response.TenantsDetails, error) {

	// resolve template id, in case client used name instead id
	vimTenants, err := a.rest.GetVimTenants()
	if err != nil {
		return nil, err
	}

	tenant, err := vimTenants.FindCloudLink(cloudUrl)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

// Validate cloud tenant state
func (a *TcaApi) validateTenant(tenant *response.TenantsDetails) error {

	glog.Infof("Validating cloud and tenant state.")

	if tenant.VimConn.Status != "ok" {
		return fmt.Errorf("cloud provider currently disconected")
	}

	if strings.ToLower(tenant.VimType) == response.VimTypeKubernetes {
		return fmt.Errorf("cloud provider already set to kubernetes")
	}
	return nil
}

func (a *TcaApi) validatePlacement(
	vmwareVim *models.VMwareClusters,
	folders *models.Folders,
	rps *models.ResourcePool,
	param models.PlacementParams) error {

	if models.VmwareTypeFolder(param.Type) == models.TypeFolder {
		if folders.IsValidFolder(param.Name) == false {
			return fmt.Errorf("failed find a target folder")
		} else {
			glog.Infof("Resolved remote datastore folder.")
		}
	}
	if models.VmwareDatastore(param.Type) == models.TypeDataStore {
		if vmwareVim.IsValidDatastore(param.Name) == false {
			return fmt.Errorf("failed find a target datastore")
		} else {
			glog.Infof("Resolved remote datastore name.")
		}
	}
	if models.VmwareResourcePool(param.Type) == models.TypeResourcePool {
		if rps.IsValidResource(param.Name) == false {
			return fmt.Errorf("failed find a target resource pool")
		} else {
			glog.Infof("Resolved remote resource pool.")
		}
	}
	if models.ClusterComputeResource(param.Type) == models.TypeClusterComputeResource {
		if vmwareVim.IsValidClusterCompute(param.Name) == false {
			return fmt.Errorf("failed find a cluster compute resource")
		} else {
			glog.Infof("Resolved remote cluster name.")
		}
	}

	return nil
}

// Validate cloud tenant state
func (a *TcaApi) validatePlacements(spec *request.Cluster, tenant *response.TenantsDetails) error {

	if tenant.VimConn.Status != "ok" {
		return fmt.Errorf("cloud provider currently disconected")
	}

	if tenant.IsVMware() {
		// vc compute
		vmwareVim, err := a.GetVimComputeClusters(tenant.VimName)
		if err != nil {
			return nil
		}
		// vc folders
		folders, err := a.GetVimFolders(tenant.VimName)
		if err != nil {
			return nil
		}
		// vc resource pools
		rps, err := a.GetVimResourcePool(tenant.VimName)
		if err != nil {
			return nil
		}

		for _, param := range spec.PlacementParams {
			if err := a.validatePlacement(vmwareVim, folders, rps, param); err != nil {
				return err
			}
		}
		for _, worker := range spec.WorkerNodes {
			for _, param := range worker.PlacementParams {
				if err := a.validatePlacement(vmwareVim, folders, rps, param); err != nil {
					return err
				}
			}
		}
		for _, master := range spec.MasterNodes {
			for _, param := range master.PlacementParams {
				if err := a.validatePlacement(vmwareVim, folders, rps, param); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// CreateClusters - create new cluster
// in Dry Run we only parse
func (a *TcaApi) CreateClusters(spec *request.Cluster, isDry bool) (bool, error) {

	if a.rest == nil {
		return false, fmt.Errorf("rest interface is nil")
	}

	if spec == nil {
		return false, fmt.Errorf("new cluster spec can't be nil")
	}

	spec.ClusterType = strings.ToUpper(spec.ClusterType)

	// validate cluster
	if err := validateClusterSpec(spec); err != nil {
		return false, err
	}

	// do all sanity check here.
	tenants, err := a.rest.GetClusters()
	if err != nil {
		return false, err
	}

	_, err = tenants.GetClusterId(spec.Name)
	// swap name
	if err == nil {
		spec.Name = spec.Name + "-" + uuid.New().String()
		spec.Name = string(spec.Name[0:25])
	}

	// resolve template id, and cluster type
	spec.ClusterTemplateId, err = a.ResolveTemplateId(spec.ClusterTemplateId, spec.ClusterType)
	if err != nil {
		return false, err
	}

	// get template and validate specs
	t, err := a.rest.GetClusterTemplate(spec.ClusterTemplateId)
	if err != nil {
		return false, err
	}

	glog.Infof("Validating node pool specs.")
	_, err = t.ValidateSpec(spec)
	if err != nil {
		return false, err
	}

	glog.Infof("Resolved template id %v", spec.ClusterTemplateId)
	tenant, err := a.validateCloudEndpoint(spec.HcxCloudUrl)
	if err != nil {
		return false, err
	}

	err = a.validateTenant(tenant)
	if err != nil {
		return false, err
	}

	if spec.ClusterType == string(request.ClusterWorkload) {
		// resolve template id, in case client used name instead id
		mgmtClusterId, err := tenants.GetClusterId(spec.ManagementClusterId)
		if err != nil {
			return false, err
		}
		mgmtSpec, err := tenants.GetClusterSpec(mgmtClusterId)
		if err != nil {
			return false, err
		}
		if mgmtSpec.ClusterType != string(request.ClusterManagement) {
			return false, fmt.Errorf("managementClusterId cluster id is not management cluster")
		}
		glog.Infof("Resolved cluster name %v", mgmtClusterId)
		spec.ManagementClusterId = mgmtClusterId

		err = a.validatePlacements(spec, tenant)
		if err != nil {
			return false, err
		}
	}

	if spec.ClusterType == string(request.ClusterManagement) {
		// ignoring mgmt cluster id
		spec.ManagementClusterId = ""
		err = a.validatePlacements(spec, tenant)
		if err != nil {
			return false, err
		}
	}

	spec.ClusterPassword = b64.StdEncoding.EncodeToString([]byte(spec.ClusterPassword))

	if isDry {
		return true, nil
	}

	return a.rest.CreateCluster(spec)
}

// DeleteTenantCluster - deletes tenant cluster
func (a *TcaApi) DeleteTenantCluster(tenantCluster string) (bool, error) {

	if a.rest == nil {
		return false, fmt.Errorf("rest interface is nil")
	}

	// TODO add validation and name resolution
	return a.rest.DeleteTenant(tenantCluster)
}

// DeleteTemplate delete cluster template from TCA
// template argument can be name or ID.
func (a *TcaApi) DeleteTemplate(template string) error {

	if a.rest == nil {
		return fmt.Errorf("rest interface is nil")
	}

	var templateId = ""

	if IsValidUUID(template) {
		tmpl, err := a.rest.GetClusterTemplate(template)
		if err != nil {
			return err
		}
		templateId = tmpl.Id
	} else {
		templates_, err := a.rest.GetClusterTemplates()
		if err != nil {
			return err
		}
		templateId, err = templates_.GetTemplateId(template)
		if err != nil {
			return err
		}
		glog.Infof("Resolved template id %s", templateId)
	}

	err := a.rest.DeleteClusterTemplate(templateId)
	if err != nil {
		return err
	}

	fmt.Printf("Template %v deleted.", templateId)

	return nil
}

// GetVdu retrieve Vdu
func (a *TcaApi) GetVdu(nfdName string) (*response.VduPackage, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	vnfCatalog, err := a.rest.GetVnfPkgm("", "")
	if err != nil || vnfCatalog == nil {
		glog.Errorf("Failed retrieve vnf package information, error %v", err)
		return nil, err
	}

	pkgCnf, err := vnfCatalog.GetVnfdID(nfdName)
	if err != nil || pkgCnf == nil {
		glog.Errorf("Failed retrieve vnfd information for %v.", nfdName)
		return nil, err
	}

	vnfd, err := a.rest.GetVnfPkgmVnfd(pkgCnf.PID)
	if err != nil || vnfd == nil {
		glog.Error("Failed acquire VDU information for %v.", pkgCnf.PID)
		return nil, err
	}

	return vnfd, nil
}

// GetRepos retrieve repos
func (a *TcaApi) GetRepos() (*response.ReposList, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.rest.GetVimTenants()
	if err != nil {
		glog.Error(err)
	}

	var allRepos response.ReposList
	for _, r := range tenants.TenantsList {
		repos, err := a.rest.RepositoriesQuery(&request.RepoQuery{
			QueryFilter: request.Filter{
				ExtraFilter: request.AdditionalFilters{
					VimID: r.TenantID,
				},
			},
		})

		if err != nil {
			return nil, err
		}

		allRepos.Items = append(allRepos.Items, repos.Items...)
	}

	return &allRepos, nil
}

// GetVim return vim
func (a *TcaApi) GetVim(vimId string) (*response.TenantSpecs, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	inputs := strings.Split(vimId, "_")
	if len(inputs) != 2 {
		return nil, &InvalidVimFormat{errMsg: vimId}
	}

	return a.rest.GetVim(vimId)
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
) *request.Cluster {

	var spec request.Cluster
	spec.Name = name
	spec.ClusterType = string(request.ClusterManagement)
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

	net := models.NewNetworks(string(request.ClusterManagement),
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

func (a *TcaApi) GetClusters() (*response.Clusters, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	glog.Infof("Retrieving cluster list.")
	return a.rest.GetClusters()
}

func (a *TcaApi) GetClusterNodePools(Id string) (*response.NodePool, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	glog.Infof("Retrieving cluster list.")

	return a.rest.GetClusterNodePools(Id)
}

// GetAllNodePools - return all node pool for clusterId
func (a *TcaApi) GetAllNodePools() ([]response.NodesSpecs, error) {

	glog.Infof("Retrieving node pools.")

	var allSpecs []response.NodesSpecs
	if a.rest == nil {
		return allSpecs, fmt.Errorf("rest interface is nil")
	}

	clusters, err := a.GetClusters()
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

// GetVnfPkgm - return packages
func (a *TcaApi) GetVnfPkgm(filter string, id string) (*response.VnfPackages, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetVnfPkgm(filter, id)
}

// GetCatalogId return vnf Package ID and VNFD ID
func (a *TcaApi) GetCatalogId(catalogId string) (string, string, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return "", "", fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetPackageCatalogId(catalogId)
}

// GetTenantsQuery query tenant information based on
// tenant id and package id
func (a *TcaApi) GetTenantsQuery(tenantId string, nfType string) (*response.Tenants, error) {

	glog.Infof("Retrieving tenants.")

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	reqFilter := request.TenantsNfFilter{
		Filter: request.TenantFilter{
			NfType: nfType,
		},
	}

	// attach tenant id if need
	if len(tenantId) > 0 {
		cid, vnfdId, err := a.GetCatalogId(tenantId)
		if err != nil {
			return nil, nil
		}
		glog.Infof("Acquired catalog id '%v', for vnfId '%v'", cid, vnfdId)
		reqFilter.Filter.NfdId = vnfdId
	}

	return a.rest.GetTenantsQuery(&reqFilter)
}

// GetClusterTemplates - return list of cluster templates
func (a *TcaApi) GetClusterTemplates() (*response.ClusterTemplates, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetClusterTemplates()
}

// GetNamedClusterTemplate - return template
// if name is string first resolve template
func (a *TcaApi) GetNamedClusterTemplate(name string) (*response.ClusterTemplate, error) {

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

// GetClusterTemplate return cluster template
func (a *TcaApi) GetClusterTemplate(templateId string) (*response.ClusterTemplate, error) {

	glog.Infof("Retrieving vnf packages.")

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetClusterTemplate(templateId)
}

// CreateNewPackage method create a new package
// it take file name that must compressed zip file
// package catalog name and a substitution map.
// substitution map used to replace CSAR values.
// a key of map is key in CSAR and value a new value
// that used to replace value in actual CSAR.
// i.e  existing CSAR used as template and substitution
// map applied a transformation.
func (a *TcaApi) CreateNewPackage(
	fileName string,
	catalogName string,
	substitution map[string]string) (bool, error) {

	glog.Infof("Create new package. Received substitution %v.", substitution)

	if a.rest == nil {
		return false, fmt.Errorf("rest interface is nil")
	}

	// Apply transformation to a CSAR file
	newCsarFile, err := csar.ApplyTransformation(
		fileName,
		csar.SpecNfd,                    // a file inside a CSAR that we need apply transformation
		csar.NfdYamlPropertyTransformer, // a callback that apply transformation
		substitution)
	if err != nil {
		glog.Errorf("Failed apply transformation %v", err)
		return false, err
	}

	file, err := os.Open(newCsarFile)
	if err != nil {
		glog.Errorf("Failed read , newly generated csar %v", err)
		return false, err
	}

	// Read new CSAR file, to buffer
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		glog.Errorf("Failed read generated csar %v", err)
		return false, err
	}

	newFileName := filepath.Base(newCsarFile)
	uploadReq := client.NewPackageUpload(catalogName)
	respond, err := a.rest.CreateVnfPkgmVnfd(uploadReq)
	if err != nil {
		glog.Errorf("Failed create cnf package entity generated csar %v", err)
		return false, err
	}

	if len(respond.Id) == 0 {
		glog.Errorf("Something is wrong, server must contain package id in respond", err)
		return false, fmt.Errorf("respond doesn't contain package id")

	}

	// upload csar to a catalog
	ok, err := a.rest.UploadVnfPkgmVnfd(respond.Id, fileBytes, newFileName)
	if err != nil {
		return false, err
	}

	// TODO do GET to cross check and respond with ok if package is created.
	return ok, nil
}
