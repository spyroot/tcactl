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
package app

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spyroot/hestia/cmd/api"
	"github.com/spyroot/hestia/cmd/client"
	"github.com/spyroot/hestia/cmd/client/main/app/ui"
	"github.com/spyroot/hestia/cmd/client/printers"
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/cmd/client/response"
	"github.com/spyroot/hestia/cmd/models"
	"github.com/spyroot/hestia/pkg/io"
	"os"
)

const (
	// ConfigDefaultPinter default printer.
	ConfigDefaultPinter = "default"

	// ConfigJsonPinter json printers
	ConfigJsonPinter = "json"

	// ConfigYamlPinter yaml printers
	ConfigYamlPinter = "yaml"

	//FilteredOutFilter - Filtered output printer
	FilteredOutFilter = "filtered"

	// ConfigFile default config file name.
	ConfigFile = "config"

	// ConfigFormat specifies default config format.
	ConfigFormat = "yaml"

	// ConfigTcaEndpoint URI endpoint.
	ConfigTcaEndpoint = "tca-endpoint"

	// ConfigTcaUsername TCA Username
	ConfigTcaUsername = "tca-username"

	// ConfigTcaPassword specifies password to authenticate
	ConfigTcaPassword = "tca-password"

	// ConfigDefaultCluster default cluster name used to do a placement
	ConfigDefaultCluster = "defaultCluster"

	// ConfigNodePool default node pool name used for placement
	ConfigNodePool = "defaultNodePool"

	// ConfigDefaultCloud default target cloud
	ConfigDefaultCloud = "defaultCloud"

	// ConfigStderrThreshold - default logging level
	ConfigStderrThreshold = "stderrthreshold"

	// ConfigHarborEndpoint Default harbor end point
	ConfigHarborEndpoint = "harbor-endpoint"

	// ConfigHarborUsername Default harbor username
	ConfigHarborUsername = "harbor-username"

	// ConfigHarborPassword default harbor password
	ConfigHarborPassword = "harbor-password"

	// FlagOutput - default logging level
	FlagOutput = "output"

	// FlagConfig config
	FlagConfig = "config"

	// FlagCliWide wide output
	FlagCliWide = "wide"

	//FlagCliTerm normal terminal mode no color.
	FlagCliTerm = "term"
)

type TcaCtl struct {

	// rest client (TODO this need to be slowly refactored and moved to TcaAPI abstraction)
	TcaClient *client.RestClient

	//
	HarborClient *client.RestClient

	// api interface
	tca *api.TcaApi

	// CnfInstancePrinters cnf instance printer
	CnfInstancePrinters map[string]func(*response.Cnfs, ui.PrinterStyle)

	//
	CnfInstanceExtendedPrinters map[string]func(*response.CnfsExtended, ui.PrinterStyle)

	// CnfPackagePrinters cnf catalog packages printer
	CnfPackagePrinters map[string]func(*response.VnfPackages, ui.PrinterStyle)

	// RepoPrinter repositories printer
	RepoPrinter map[string]func(*response.ReposList, ui.PrinterStyle)

	// TenantsPrinter cloud tenant printer
	TenantsPrinter map[string]func(*response.Tenants, ui.PrinterStyle)

	// NodePoolPrinter k8s node pool printer
	NodePoolPrinter map[string]func(*response.NodePool, ui.PrinterStyle)

	// PoolSpecTablePrinter k8s single pool printer
	PoolSpecPrinter map[string]func(*response.NodesSpecs, ui.PrinterStyle)

	// ClustersPrinter k8s cluster printer
	ClustersPrinter map[string]func(*response.Clusters, ui.PrinterStyle)

	// ClusterPrinter k8s cluster printer
	ClusterPrinter map[string]func(*response.ClusterSpec, ui.PrinterStyle)

	// VduPrinter vdu printer
	VduPrinter map[string]func(*response.VduPackage, ui.PrinterStyle)

	// TenantQueryPrinter tenant query printer
	TenantQueryPrinter map[string]func(*response.Tenants, ui.PrinterStyle)

	// TenantQueryPrinter tenant query printer
	NodesPrinter map[string]func(*response.NodePool, ui.PrinterStyle)

	// TenantQueryPrinter tenant query printer
	TemplatePrinter map[string]func(*response.ClusterTemplate, ui.PrinterStyle)

	// TenantQueryPrinter tenant query printer
	TemplatesPrinter map[string]func([]response.ClusterTemplate, ui.PrinterStyle)

	// TenantQueryPrinter tenant query printer
	ClusterRequestPrinter map[string]func(*request.Cluster, ui.PrinterStyle)

	// cloud tenant printer
	TenantsResponsePrinter map[string]func(*response.TenantSpecs, ui.PrinterStyle)

	// Vmware specific infra printers
	VMwareClusterPrinter    map[string]func(*models.VMwareClusters, ui.PrinterStyle)
	VMwareDatastorePrinter  map[string]func(*models.VMwareClusters, ui.PrinterStyle)
	VmwareNetworkPrinter    map[string]func(*models.CloudNetworks, ui.PrinterStyle)
	VmwareVmTemplatePrinter map[string]func(*models.VcInventory, ui.PrinterStyle)
	VmwareResourcePrinter   map[string]func(*models.ResourcePool, ui.PrinterStyle)

	// cluster task list printer.  cluster task is take current executing or
	// already executed.
	TaskClusterPrinter map[string]func(*models.ClusterTask, ui.PrinterStyle)

	// global flag what output printer to use
	Printer string

	// global debug flag for a tool
	IsDebug bool

	// config file
	CfgFile string

	// root entry for cli
	RootCmd *cobra.Command

	// default style for table Printer
	DefaultStyle ui.PrinterStyle

	// DefaultClusterName cluster name from config or flag
	DefaultClusterName string

	// DefaultCloudName default cloud name tool will use
	DefaultCloudName string

	// DefaultNodeName node pool tool will use.
	DefaultNodeName string

	// IsColorTerm color or not term
	IsColorTerm bool

	// IsWideTerm is wide or not output
	IsWideTerm bool

	// Harbor harbor api end-point
	Harbor string

	// HarborUsername harbor username
	HarborUsername string

	// HarborPassword harbor password
	HarborPassword string
}

func NewTcaCtl() *TcaCtl {

	ctl := TcaCtl{
		TcaClient: nil,
		CnfInstancePrinters: map[string]func(*response.Cnfs, ui.PrinterStyle){
			ConfigDefaultPinter: printer.CnfInstanceTablePrinter,
			ConfigJsonPinter:    printer.CnfInstanceJsonPrinter,
		},
		CnfInstanceExtendedPrinters: map[string]func(*response.CnfsExtended, ui.PrinterStyle){
			ConfigDefaultPinter: printer.CnfInstanceExtendedTablePrinter,
			ConfigJsonPinter:    printer.CnfInstanceExtendedJsonPrinter,
		},
		CnfPackagePrinters: map[string]func(*response.VnfPackages, ui.PrinterStyle){
			ConfigDefaultPinter: printer.CnfPackageTablePrinter,
			ConfigJsonPinter:    printer.CnfPackageJsonPrinter,
			FilteredOutFilter:   printer.VnfPackageFilteredOutput,
		},
		RepoPrinter: map[string]func(*response.ReposList, ui.PrinterStyle){
			ConfigDefaultPinter: printer.RepoTablePrinter,
			ConfigJsonPinter:    printer.RepoJsonPrinter,
		},
		TenantsPrinter: map[string]func(*response.Tenants, ui.PrinterStyle){
			ConfigDefaultPinter: printer.TenantsTablePrinter,
			ConfigJsonPinter:    printer.TenantsJsonPrinter,
			ConfigYamlPinter:    printer.TenantsYamlPrinter,
			FilteredOutFilter:   printer.TenantsFilteredOutput,
		},
		NodePoolPrinter: map[string]func(*response.NodePool, ui.PrinterStyle){
			ConfigDefaultPinter: printer.NodePoolTablePrinter,
			ConfigJsonPinter:    printer.NodePoolJsonPrinter,
		},
		ClustersPrinter: map[string]func(*response.Clusters, ui.PrinterStyle){
			ConfigDefaultPinter: printer.ClusterTablePrinter,
			ConfigJsonPinter:    printer.ClusterJsonPrinter,
		},
		ClusterPrinter: map[string]func(*response.ClusterSpec, ui.PrinterStyle){
			ConfigDefaultPinter: printer.ClusterSpecTablePrinter,
			ConfigJsonPinter:    printer.ClusterSpecJsonPrinter,
			ConfigYamlPinter:    printer.ClusterSpecYamlPrinter,
		},
		VduPrinter: map[string]func(*response.VduPackage, ui.PrinterStyle){
			ConfigDefaultPinter: printer.VduTablePrinter,
			ConfigJsonPinter:    printer.VduJsonPrinter,
		},
		TenantQueryPrinter: map[string]func(*response.Tenants, ui.PrinterStyle){
			ConfigDefaultPinter: printer.TenantTabularPinter,
			ConfigJsonPinter:    printer.TenantJsonPrinter,
		},
		NodesPrinter: map[string]func(*response.NodePool, ui.PrinterStyle){
			ConfigDefaultPinter: printer.NodesTablePrinter,
			ConfigJsonPinter:    printer.NodesJsonPrinter,
		},
		PoolSpecPrinter: map[string]func(*response.NodesSpecs, ui.PrinterStyle){
			ConfigDefaultPinter: printer.PoolSpecTablePrinter,
			ConfigJsonPinter:    printer.PoolSpecJsonPrinter,
		},
		// printer for single template
		TemplatePrinter: map[string]func(*response.ClusterTemplate, ui.PrinterStyle){
			ConfigDefaultPinter: printer.TemplateSpecTablePrinter,
			ConfigJsonPinter:    printer.TemplateSpecJsonPrinter,
			ConfigYamlPinter:    printer.TemplateSpecYamlPrinter,
		},
		// printer for array of templates
		TemplatesPrinter: map[string]func([]response.ClusterTemplate, ui.PrinterStyle){
			ConfigDefaultPinter: printer.TemplatesSpecTablePrinter,
			ConfigJsonPinter:    printer.TemplatesJsonPrinter,
			ConfigYamlPinter:    printer.TemplatesYamlPrinter,
		},

		ClusterRequestPrinter: map[string]func(*request.Cluster, ui.PrinterStyle){
			ConfigDefaultPinter: printer.ClusterRequestJsonPrinter,
			ConfigJsonPinter:    printer.ClusterRequestJsonPrinter,
			ConfigYamlPinter:    printer.ClusterRequestYamlPrinter,
		},

		TenantsResponsePrinter: map[string]func(*response.TenantSpecs, ui.PrinterStyle){
			ConfigDefaultPinter: printer.VimTablePrinter,
			ConfigJsonPinter:    printer.TenantsResponseYamlPrinter,
			ConfigYamlPinter:    printer.TenantsResponseYamlPrinter,
		},

		TaskClusterPrinter: map[string]func(*models.ClusterTask, ui.PrinterStyle){
			ConfigDefaultPinter: printer.ClusterTaskTablePrinter,
			ConfigJsonPinter:    printer.ClusterTaskJsonPrinter,
			ConfigYamlPinter:    printer.ClusterTaskYamlPrinter,
		},

		VMwareClusterPrinter: map[string]func(*models.VMwareClusters, ui.PrinterStyle){
			ConfigDefaultPinter: printer.VmwareInventoryTablePrinter,
			ConfigJsonPinter:    printer.VmwareInventoryJsonPrinter,
			ConfigYamlPinter:    printer.VmwareInventoryYamlPrinter,
		},

		VMwareDatastorePrinter: map[string]func(*models.VMwareClusters, ui.PrinterStyle){
			ConfigDefaultPinter: printer.VmwareDatastoreTablePrinter,
			ConfigJsonPinter:    printer.VmwareInventoryJsonPrinter,
			ConfigYamlPinter:    printer.VmwareInventoryYamlPrinter,
		},

		VmwareNetworkPrinter: map[string]func(*models.CloudNetworks, ui.PrinterStyle){
			ConfigDefaultPinter: printer.VmwareNetworkTablePrinter,
			ConfigJsonPinter:    printer.VmwareNetworkJsonPrinter,
			ConfigYamlPinter:    printer.VmwareNetworkYamlPrinter,
		},

		VmwareVmTemplatePrinter: map[string]func(*models.VcInventory, ui.PrinterStyle){
			ConfigDefaultPinter: printer.VmwareTemplateTablePrinter,
			ConfigJsonPinter:    printer.VmwareTemplateJsonPrinter,
			ConfigYamlPinter:    printer.VmwareTemplateYamlPrinter,
		},
		VmwareResourcePrinter: map[string]func(*models.ResourcePool, ui.PrinterStyle){
			ConfigDefaultPinter: printer.VmwareResourcePoolTablePrinter,
			ConfigJsonPinter:    printer.VmwareResourcePoolJsonPrinter,
			ConfigYamlPinter:    printer.VmwareResourcePoolYamlPrinter,
		},

		Printer:      ConfigDefaultPinter,
		IsDebug:      false,
		CfgFile:      "",
		DefaultStyle: ui.NewTableColorStyler(),
	}

	ctl.TcaClient = &client.RestClient{
		BaseURL:  "",
		ApiKey:   "",
		SkipSsl:  true,
		Client:   nil,
		IsDebug:  true,
		Username: "",
		Password: "",
	}

	ctl.HarborClient = &client.RestClient{
		BaseURL:  "",
		ApiKey:   "",
		SkipSsl:  true,
		Client:   nil,
		IsDebug:  true,
		Username: "",
		Password: "",
	}

	tcaApi, err := api.NewTcaApi(ctl.TcaClient)
	CheckErrLogError(err)
	ctl.tca = tcaApi

	// Init root cmd callback
	ctl.RootCmd = &cobra.Command{
		Use:  "tcactl",
		Long: "tcactl ctl tool for VMware Telco Cloud Automation",
		Args: cobra.MinimumNArgs(1),
		//		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	return &ctl
}

func (ctl *TcaCtl) Authorize() error {
	ok, err := ctl.TcaClient.GetAuthorization()
	if err != nil {
		return err
	}
	if ok {
		glog.Infof("Received TCA Authorization Key %v", ctl.TcaClient.ApiKey)
	}
	return nil
}

func (ctl *TcaCtl) Authorize2() {
	ok, err := ctl.TcaClient.GetAuthorization()
	io.CheckErr(err)
	if ok {
		glog.Infof("Received TCA Authorization Key %v", ctl.TcaClient.ApiKey)
	}
}

// ResolvePoolName - resolve pool name to id in given cluster
func (ctl *TcaCtl) ResolvePoolName(s string, clusterName string) (string, string, error) {

	// empty name no ops
	if len(s) == 0 {
		return s, "", nil
	}

	if len(clusterName) == 0 {
		return "", "", fmt.Errorf("provide cluster name to resolve pool name")
	}

	nodePool, clusterId, err := ctl.TcaClient.GetNamedClusterNodePools(clusterName)
	if err != nil || nodePool == nil {
		glog.Errorf("Failed acquire clusters node information %v", err)
		return s, "", err
	}

	pool, err := nodePool.GetPool(s)
	if err != nil {
		glog.Errorf("Failed acquire node pool information %v", err)
		return s, "", err
	}
	return pool.Id, clusterId, nil
}

// ResolveClusterName - resolve cluster name to cluster id
func (ctl *TcaCtl) ResolveClusterName(q string) (string, error) {
	clusters, err := ctl.TcaClient.GetClusters()
	if err != nil {
		return "", nil
	}
	return clusters.GetClusterId(q)
}

// CheckErrLogError , print error and log error
func CheckErrLogError(msg interface{}) {
	if msg != nil {
		glog.Error(msg)
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}

func CheckNotOkLogError(predicate bool, msg interface{}) {
	if predicate != true {
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}

func CheckNilLogError(predicate interface{}, msg interface{}) {
	if predicate == nil {
		glog.Error(msg)
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}

func CheckErrLogInfoMsg(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}

// IsValidUUID check value in UUID format
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// isValidTemplateType check template type value
func isValidTemplateType(templateType string) bool {
	if templateType == response.TemplateMgmt || templateType == response.TemplateWorkload {
		return true
	}
	return false
}