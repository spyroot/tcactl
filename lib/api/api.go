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
	"github.com/go-playground/validator/v10"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/request"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/csar"
	"github.com/spyroot/tcactl/lib/models"
	"io/ioutil"
	"os"
	"path/filepath"
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

// TcaApi - TCA Api interface
// Called need to use NewTcaApi to get instance before
type TcaApi struct {

	// rest client used to interact with tca
	rest *client.RestClient

	// spec validator.
	specValidator *validator.Validate
}

// NewTcaApi - return instance for API.
//
func NewTcaApi(rest *client.RestClient) (*TcaApi, error) {

	if rest == nil {
		return nil, fmt.Errorf("rest client is nil, initilize rest client first")
	}

	a := &TcaApi{
		rest: rest,
	}

	a.specValidator = validator.New()

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
	repo string, instanceName string, nodePoolName string) *InstanceRequestSpec {
	i := &InstanceRequestSpec{cloudName: cloudName, clusterName: clusterName,
		vimType: vimType, nfdName: nfdName, repo: repo,
		instanceName: instanceName, nodePoolName: nodePoolName}

	i.flavorName = DefaultNamespace
	i.description = ""
	i.namespace = DefaultFlavor
	i.disableAutoRollback = false
	i.disableGrant = false
	i.disableAutoRollback = false

	return i
}

// GetCnfs method return list of cnf instances in
// response.CnfsExtended that encapsulate in collection
func (a *TcaApi) GetCnfs() (*response.CnfsExtended, error) {

	genericRespond, err := a.rest.GetVnflcm()
	if err != nil {
		return nil, err
	}

	// for extension request we route to correct printer
	cnfs, ok := genericRespond.(*response.CnfsExtended)
	if ok {
		return cnfs, nil
	}

	return nil, err
}

// GetAllNodePool return a Node pool for particular cluster
// It generally useful to get list only if we need to display all
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
			cluster.ClusterName, cluster.Id, cluster.Status)

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

// GetVimTenants method return vim tenants as response.Tenants
// collection.
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
	task, err := a.rest.GetClustersTask(request.NewClusterTaskQuery(taskId))
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetVimVMTemplates - return compute cluster attached to cloud provider.
// caller need indicate template type and version.
func (a *TcaApi) GetVimVMTemplates(cloudName string,
	templateType VmTemplateFilterType, ver string) (*models.VcInventory, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(cloudName) == 0 {
		return nil, fmt.Errorf("empty cloud provider name")
	}

	tenants, err := a.rest.GetVimTenants()
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

	_filter := request.NewVMwareTemplateQuery(tenant.HcxUUID, string(templateType), ver)
	t, err := a.rest.GetVMwareTemplates(_filter)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// GetVimFolders - return folders in target cloud provider.
// for VMware VC it list of VM Folders, models.Folders
func (a *TcaApi) GetVimFolders(cloudName string) (*models.Folders, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(cloudName) == 0 {
		return nil, fmt.Errorf("empty cloud provider name")
	}

	tenants, err := a.rest.GetVimTenants()
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

func validateClusterSpec(spec *request.Cluster) error {

	if len(spec.Name) == 0 {
		return fmt.Errorf("cluster initialSpec must contains name key:value")
	}
	if len(spec.ClusterPassword) == 0 {
		return fmt.Errorf("cluster initialSpec must contains cluster password key:value")
	}
	if len(spec.ClusterTemplateId) == 0 {
		return fmt.Errorf("cluster initialSpec must contains ClusterTemplateId key:value")
	}
	if len(spec.ClusterType) == 0 {
		return fmt.Errorf("cluster initialSpec must contains clusterType and value management or workload key:value")
	}
	if len(spec.HcxCloudUrl) == 0 {
		return fmt.Errorf("cluster initialSpec must contains hcxCloudUrl key:value")
	}
	if len(spec.EndpointIP) == 0 {
		return fmt.Errorf("cluster initialSpec must contains EndpointIP key:value")
	}
	if len(spec.VmTemplate) == 0 {
		return fmt.Errorf("cluster initialSpec must contains vmTemplate key:value")
	}
	if len(spec.VmTemplate) == 0 {
		return fmt.Errorf("cluster initialSpec must contains vmTemplate key:value")
	}
	if len(spec.MasterNodes) == 0 {
		return fmt.Errorf("cluster initialSpec must contains masterNodes section")
	}
	if len(spec.WorkerNodes) == 0 {
		return fmt.Errorf("cluster initialSpec must contains WorkerNodes section")
	}
	if len(spec.PlacementParams) == 0 {
		return fmt.Errorf("cluster initialSpec must contains PlacementParams section")
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

// ResolveTypedTemplateId - resolve template name to id
// for a give template type. Both must match in order
// method return ture
func (a *TcaApi) ResolveTypedTemplateId(templateId string, templateType string) (string, error) {

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
	if strings.ToLower(template.ClusterType) != strings.ToLower(templateType) {
		return "", fmt.Errorf("found template %s but template type mistmatch, "+
			"template type %v %v", templateId, template.ClusterType, templateType)
	}

	glog.Infof("Resolved template to template id %v", template.Id)

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

	if strings.ToLower(tenant.VimType) == models.VimTypeKubernetes {
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

// DeleteTenantCluster - deletes tenant cluster
func (a *TcaApi) DeleteTenantCluster(tenantCluster string) (*models.TcaTask, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	// TODO add validation and name resolution
	return a.rest.DeleteTenant(tenantCluster)
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
		glog.Errorf("Failed acquire VDU information for %v.", pkgCnf.PID)
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

// GetTenant method return tenant as response.Tenants
// if tenant is name, method will lookup by name.
// if tenant is UUID it will lookup by id
// if it has prefix vmware it will lookup by VIM id.
func (a *TcaApi) GetTenant(tenant string) (*response.Tenants, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	if len(tenant) == 0 {
		return nil, fmt.Errorf("empty tenant")
	}

	tenants, err := a.rest.GetVimTenants()
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

	reqFilter := request.TenantsNfFilter{
		Filter: request.TenantFilter{
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
		glog.Error("Something is wrong, server must contain package id in respond")
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

// GetCatalogAndVdu return catalog entity and vdu package.
func (a *TcaApi) GetCatalogAndVdu(nfdName string) (*response.VnfPackage, *response.VduPackage, error) {

	vnfCatalog, err := a.rest.GetVnfPkgm("", "")
	if err != nil || vnfCatalog == nil {
		glog.Errorf("Failed acquire vnf package information. Error %v", err)
		return nil, nil, err
	}

	catalogEntity, err := vnfCatalog.GetVnfdID(nfdName)
	if err != nil || catalogEntity == nil {
		glog.Errorf("Failed acquire catalog information for catalog name %v", nfdName)
		return nil, nil, err
	}

	v, err := a.rest.GetVnfPkgmVnfd(catalogEntity.PID)

	return catalogEntity, v, err
}

// CreateCnfNewInstance create a new instance of VNF or CNF.
// Dry run will validate request but will not create any CNF.
func (a *TcaApi) CreateCnfNewInstance(ctx context.Context, n *InstanceRequestSpec, isDry bool, isBlocked bool) (*response.LcmInfo, error) {

	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	tenants, err := a.GetVimTenants()
	if err != nil {
		glog.Errorf("Failed acquire cloud tenant, error: %v", err)
		return nil, err
	}

	glog.Infof("Acquiring cloud provider %s details, type %s", n.cloudName, n.vimType)
	cloud, err := tenants.GetTenantClouds(n.cloudName, n.vimType)
	if err != nil {
		glog.Errorf("Failed acquire cloud provider details, error: %v", err)
		return nil, err
	}

	glog.Infof("Acquiring catalog information for entity %s", n.nfdName)
	pkg, vnfd, err := a.GetCatalogAndVdu(n.nfdName)
	if err != nil || vnfd == nil {
		glog.Errorf("Failed acquire VDU information for %v", n.nfdName)
		return nil, err
	}

	// get linked repo, if caller provide repo that is not
	// linked nothing to do.
	reposUuid, err := a.rest.LinkedRepositories(cloud.TenantID, n.repo)
	if err != nil {
		glog.Errorf("Failed acquire linked %v "+
			"repository to cloud provider %v. Indicate a repo "+
			"linked to cloud provider.", n.repo, cloud.TenantID)
		return nil, err
	}

	ext, err := a.rest.ExtensionQuery()
	if err != nil {
		glog.Errorf("Failed acquire extension information for %v", err)
		return nil, err
	}

	linkedRepos, err := ext.FindRepo(reposUuid)
	if err != nil || linkedRepos == nil {
		glog.Errorf("Failed acquire extension information for %v", reposUuid)
		return nil, err
	}

	if linkedRepos.IsEnabled() == false {
		glog.Errorf("Repository %v is disabled", linkedRepos.Name)
		return nil, fmt.Errorf("repository %v is disabled", linkedRepos.Name)
	}

	glog.Infof("Found attached repo %v and status %v", n.repo, linkedRepos.State)

	// resolve nodePools
	nodePool, _, err := a.rest.GetNamedClusterNodePools(n.clusterName)
	if err != nil || nodePool == nil {
		glog.Errorf("Failed acquire clusters node information for cluster %v, error %v", n.clusterName, err)
		return nil, err
	}
	pool, err := nodePool.GetPool(n.nodePoolName)
	if err != nil {
		glog.Errorf("Failed acquire node pool information for node pool %v, error %v", n.nodePoolName, err)
		return nil, err
	}

	if isDry == true {
		return nil, nil
	}

	vnfLcm, err := a.rest.CreateInstance(ctx, &request.CreateVnfLcm{
		VnfdId:                 pkg.VnfdID,
		VnfInstanceName:        n.instanceName,
		VnfInstanceDescription: n.description,
	})

	if err != nil {
		glog.Errorf("Failed create instance information %v", err)
		return nil, err
	}

	var flavorName = n.flavorName
	if len(vnfd.Vnf.Properties.FlavourId) > 0 {
		flavorName = vnfd.Vnf.Properties.FlavourId
	}

	for _, vdu := range vnfd.Vdus {
		glog.Infof("Instantiating vdu %s %s", vdu.VduId, linkedRepos.Name)
		var req = request.InstantiateVnfRequest{
			FlavourID: flavorName,
			VimConnectionInfo: []request.VimConInfo{
				{
					ID:      cloud.VimID,
					VimType: "",
					Extra: request.PoolExtra{
						NodePoolId: pool.Id,
					},
				},
			},
			AdditionalVduParams: request.AdditionalParams{
				VduParams: []request.VduParam{{
					Namespace: n.namespace,
					RepoURL:   n.repo,
					Username:  linkedRepos.AccessInfo.Username,
					Password:  linkedRepos.AccessInfo.Password,
					VduName:   vdu.VduId,
				}},
				DisableGrant:        n.disableGrant,
				IgnoreGrantFailure:  n.ignoreGrantFailure,
				DisableAutoRollback: n.disableAutoRollback,
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

				if strings.HasPrefix(instance.InstantiationState, waitFor) {
					break
				}

				time.Sleep(TaskWaitSeconds * time.Second)
			}
		}
	}

	return nil
}

// BlockWaitTaskFinish - simple block and pull status
// instanceId is instance that method will pull and check
// waitFor is target status method waits.
// maxRetry a limit.
func (a *TcaApi) BlockWaitTaskFinish(ctx context.Context, task *models.TcaTask, waitFor string, maxRetry int, verbose bool) error {

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	req := request.NewClusterTaskQuery(task.Id)

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

				task, err := a.rest.GetClustersTask(req)

				if err != nil {
					return err
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
func (a *TcaApi) ResolvePoolId(poolName string, clusterScope ...string) (string, error) {

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
	clusters, err := a.rest.GetClusters()
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
func (a *TcaApi) ResolvePoolName(poolName string, clusterName string) (string, string, error) {

	// empty name no ops
	if len(poolName) == 0 {
		return poolName, "", nil
	}

	if len(clusterName) == 0 {
		return "", "", fmt.Errorf("provide cluster name to resolve pool name")
	}

	nodePool, clusterId, err := a.rest.GetNamedClusterNodePools(clusterName)
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
// and return a string version.
func (a *TcaApi) ResolveClusterName(q string) (string, error) {

	clusters, err := a.rest.GetClusters()
	if err != nil {
		return "", nil
	}

	return clusters.GetClusterId(q)
}

// SetBaseUrl set base url
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

func (a *TcaApi) GetApiKey() string {

	if a != nil && a.rest != nil {
		return a.rest.GetApiKey()
	}

	return ""
}

func (a *TcaApi) SetTrace(trace bool) {
	a.rest.SetTrace(trace)
}

func (a *TcaApi) GetVims() (*response.Tenants, error) {
	if a.rest == nil {
		return nil, fmt.Errorf("rest interface is nil")
	}

	return a.rest.GetVimTenants()
}
