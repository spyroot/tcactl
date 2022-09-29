// Package cmds
// Copyright 2020-2021 Author.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Mustafa mbayramo@vmware.com
package cmds

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spyroot/tcactl/app/main/cmds/ui"
	"github.com/spyroot/tcactl/lib/api"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/lib/client/printers"
	"github.com/spyroot/tcactl/lib/client/response"
	"github.com/spyroot/tcactl/lib/client/specs"
	"github.com/spyroot/tcactl/lib/models"
	"github.com/spyroot/tcactl/pkg/io"
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

	// ConfigRepoName default repo name
	ConfigRepoName = "defaultRepoName"

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

	// ConfigVcUrl default vc url
	ConfigVcUrl = "vc"

	// ConfigVcUsername default repo name
	ConfigVcUsername = "vc_username"

	// ConfigVcPassword default repo name
	ConfigVcPassword = "vc_password"

	// ConfigTrace dump server respond
	ConfigTrace = "trace"

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

	//
	HarborClient *client.RestClient

	// API interface
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
	TemplatePrinter map[string]func(*response.ClusterTemplateSpec, ui.PrinterStyle)

	// TenantQueryPrinter tenant query printer
	TemplatesPrinter map[string]func([]response.ClusterTemplateSpec, ui.PrinterStyle)

	// TenantQueryPrinter tenant query printer
	ClusterRequestPrinter map[string]func(*specs.SpecCluster, ui.PrinterStyle)

	// cloud tenant printer
	TenantsResponsePrinter map[string]func(*response.TenantSpecs, ui.PrinterStyle)

	// Vmware specific infra printers
	VMwareClusterPrinter    map[string]func(*models.VMwareClusters, ui.PrinterStyle)
	VMwareDatastorePrinter  map[string]func(*models.VMwareClusters, ui.PrinterStyle)
	VmwareNetworkPrinter    map[string]func(*models.CloudNetworks, ui.PrinterStyle)
	VmwareVmTemplatePrinter map[string]func(*models.VcInventory, ui.PrinterStyle)
	VmwareResourcePrinter   map[string]func(*models.ResourcePool, ui.PrinterStyle)

	// Tca consumption specific printers
	TcaConsumptionPrinter map[string]func(*models.ConsumptionResp, ui.PrinterStyle)

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

	// DefaultNodePoolName node pool tool will use.
	DefaultNodePoolName string

	// DefaultRepoName default repo name
	DefaultRepoName string

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

	IsTrace    bool
	VcUrl      string
	VcUsername string
	VcPassword string
}

// NewTcaCtl - main abstraction for a tool
func NewTcaCtl() *TcaCtl {

	ctl := TcaCtl{
		//TcaClient: nil,
		CnfInstancePrinters: map[string]func(*response.Cnfs, ui.PrinterStyle){
			ConfigDefaultPinter: printer.CnfInstanceTablePrinter,
			ConfigJsonPinter:    printer.CnfInstanceJsonPrinter,
			ConfigYamlPinter:    printer.CnfInstanceYamlPrinter,
		},
		CnfInstanceExtendedPrinters: map[string]func(*response.CnfsExtended, ui.PrinterStyle){
			ConfigDefaultPinter: printer.CnfInstanceExtendedTablePrinter,
			ConfigJsonPinter:    printer.CnfInstanceExtendedJsonPrinter,
			ConfigYamlPinter:    printer.CnfInstanceExtendedYamlPrinter,
			FilteredOutFilter:   printer.CnfsExtendedFilteredOutput,
		},
		CnfPackagePrinters: map[string]func(*response.VnfPackages, ui.PrinterStyle){
			ConfigDefaultPinter: printer.CnfPackageTablePrinter,
			ConfigJsonPinter:    printer.CnfPackageJsonPrinter,
			ConfigYamlPinter:    printer.CnfPackageYamlPrinter,
			FilteredOutFilter:   printer.VnfPackageFilteredOutput,
		},
		RepoPrinter: map[string]func(*response.ReposList, ui.PrinterStyle){
			ConfigDefaultPinter: printer.RepoTablePrinter,
			ConfigJsonPinter:    printer.RepoJsonPrinter,
			ConfigYamlPinter:    printer.RepoYamlPrinter,
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
			ConfigYamlPinter:    printer.NodePoolYamlPrinter,
		},
		ClustersPrinter: map[string]func(*response.Clusters, ui.PrinterStyle){
			ConfigDefaultPinter: printer.ClusterTablePrinter,
			ConfigJsonPinter:    printer.ClusterJsonPrinter,
			ConfigYamlPinter:    printer.ClusterYamlPrinter,
		},
		ClusterPrinter: map[string]func(*response.ClusterSpec, ui.PrinterStyle){
			ConfigDefaultPinter: printer.ClusterSpecTablePrinter,
			ConfigJsonPinter:    printer.ClusterSpecJsonPrinter,
			ConfigYamlPinter:    printer.ClusterSpecYamlPrinter,
		},
		VduPrinter: map[string]func(*response.VduPackage, ui.PrinterStyle){
			ConfigDefaultPinter: printer.VduTablePrinter,
			ConfigJsonPinter:    printer.VduJsonPrinter,
			ConfigYamlPinter:    printer.VduYamlPrinter,
		},
		TenantQueryPrinter: map[string]func(*response.Tenants, ui.PrinterStyle){
			ConfigDefaultPinter: printer.TenantTabularPinter,
			ConfigJsonPinter:    printer.TenantJsonPrinter,
			ConfigYamlPinter:    printer.TenantYamlPrinter,
		},
		NodesPrinter: map[string]func(*response.NodePool, ui.PrinterStyle){
			ConfigDefaultPinter: printer.NodesTablePrinter,
			ConfigJsonPinter:    printer.NodesJsonPrinter,
			ConfigYamlPinter:    printer.NodesYamlPrinter,
		},
		PoolSpecPrinter: map[string]func(*response.NodesSpecs, ui.PrinterStyle){
			ConfigDefaultPinter: printer.PoolSpecTablePrinter,
			ConfigJsonPinter:    printer.PoolSpecJsonPrinter,
		},
		// printer for single template
		TemplatePrinter: map[string]func(*response.ClusterTemplateSpec, ui.PrinterStyle){
			ConfigDefaultPinter: printer.TemplateSpecTablePrinter,
			ConfigJsonPinter:    printer.TemplateSpecJsonPrinter,
			ConfigYamlPinter:    printer.TemplateSpecYamlPrinter,
		},
		// printer for array of templates
		TemplatesPrinter: map[string]func([]response.ClusterTemplateSpec, ui.PrinterStyle){
			ConfigDefaultPinter: printer.TemplatesSpecTablePrinter,
			ConfigJsonPinter:    printer.TemplatesJsonPrinter,
			ConfigYamlPinter:    printer.TemplatesYamlPrinter,
		},

		ClusterRequestPrinter: map[string]func(*specs.SpecCluster, ui.PrinterStyle){
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

		TcaConsumptionPrinter: map[string]func(*models.ConsumptionResp, ui.PrinterStyle){
			ConfigDefaultPinter: printer.ConsumptionTablePrinter,
			ConfigJsonPinter:    printer.ConsumptionJsonPrinter,
			ConfigYamlPinter:    printer.ConsumptionSpecYamlPrinter,
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

	ctl.HarborClient = &client.RestClient{
		BaseURL:  "",
		SkipSsl:  true,
		Client:   nil,
		IsDebug:  true,
		Username: "",
		Password: "",
	}

	tcaApi, err := api.NewTcaApi(&client.RestClient{
		BaseURL:  "",
		SkipSsl:  true,
		Client:   nil,
		IsDebug:  true,
		Username: "",
		Password: "",
	})

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

// Authorize authenticate and obtain a session.
// TODO this method will go away
func (ctl *TcaCtl) Authorize() error {
	ok, err := ctl.tca.GetAuthorization()
	if err != nil {
		return err
	}
	if ok {
		glog.Infof("Received TCA Authorization Key %v", ctl.tca.GetApiKey())
	}
	return nil
}

// BasicAuthentication TODO this method will go away
func (ctl *TcaCtl) BasicAuthentication() {
	ok, err := ctl.tca.GetAuthorization()
	io.CheckErr(err)
	if ok {
		glog.Infof("Received TCA Authorization Key %v", ctl.tca.GetApiKey())
	}
}

// ResolvePoolName - resolve pool name to id in given cluster
// TODO this method will go away
func (ctl *TcaCtl) ResolvePoolName(poolName string, clusterName string) (string, string, error) {
	return ctl.tca.ResolvePoolName(context.Background(), poolName, clusterName)
}

// ResolveClusterName - resolve cluster name to cluster id
// TODO this method will go away
func (ctl *TcaCtl) ResolveClusterName(q string) (string, error) {
	return ctl.tca.ResolveClusterName(context.Background(), q)
}

// SetTcaBase sets tca base url
func (ctl *TcaCtl) SetTcaBase(url string) {
	if ctl.tca != nil {
		ctl.tca.SetBaseUrl(url)
	}
}

// SetTcaUsername sets tca username
func (ctl *TcaCtl) SetTcaUsername(username string) {
	if ctl.tca != nil {
		ctl.tca.SetUsername(username)
	}
}

// SetPassword sets tca password
func (ctl *TcaCtl) SetPassword(password string) {
	if ctl.tca != nil {
		ctl.tca.SetPassword(password)
	}
}

// GetApi returns TcaApi api.TcaApi
func (ctl *TcaCtl) GetApi() *api.TcaApi {
	return ctl.tca
}

// CheckErrLogError , print error and log error
func CheckErrLogError(msg interface{}) {
	if msg != nil {
		glog.Error(msg)
		_, err := fmt.Fprintln(os.Stderr, "Error:", msg)
		if err != nil {
			fmt.Printf("Failed to write %v", err)
			return
		}
		os.Exit(1)
	}
}

func CheckNotOkLogError(predicate bool, msg interface{}) {
	if predicate != true {
		_, err := fmt.Fprintln(os.Stderr, "Error:", msg)
		if err != nil {
			fmt.Printf("Failed to write %v", err)
			return
		}
		os.Exit(1)
	}
}

func CheckNilLogError(predicate interface{}, msg interface{}) {
	if predicate == nil {
		glog.Error(msg)
		_, err := fmt.Fprintln(os.Stderr, "Error:", msg)
		if err != nil {
			fmt.Printf("Failed to write %v", err)
			return
		}
		os.Exit(1)
	}
}

func CheckErrLogInfoMsg(msg interface{}) {
	if msg != nil {
		_, err := fmt.Fprintln(os.Stderr, "Error:", msg)
		if err != nil {
			fmt.Printf("Failed to write to stderr %v", err)
			return
		}
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
	if templateType == models.TemplateMgmt || templateType == models.TemplateWorkload {
		return true
	}
	return false
}
