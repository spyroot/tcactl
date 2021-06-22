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
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spyroot/hestia/cmd/client/request"
	"github.com/spyroot/hestia/pkg/io"
)

func (ctl *TcaCtl) CmdCreateCnf() *cobra.Command {

	var DisableGrant = true
	var namespace = "mwcapp"
	var defaultVimType = "KUBERNETES"
	var defaultRepo = ""

	var cmdCreate = &cobra.Command{
		Use:   "cnf [cluster name, catalog name/id, instance name",
		Short: "Create a new cnf instance.",
		Long:  `Create a new cnf instance.`,
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {

			var cloudName = args[0]
			var nfdName = args[1]
			var instanceName = args[2]

			tenants, err := ctl.TcaClient.GetVimTenants()
			CheckErrLogError(err)

			cloud, err := tenants.GetTenantClouds(cloudName, defaultVimType)
			CheckErrLogError(err)

			vnfCatalog, err := ctl.TcaClient.GetVnfPkgm("", "")
			if err != nil || vnfCatalog == nil {
				glog.Errorf("Failed acquire vnf package information.")
				return
			}
			pkgCnf, err := vnfCatalog.GetVnfdID(nfdName)
			if err != nil || pkgCnf == nil {
				glog.Errorf("Failed acquire vnfd information for %v", nfdName)
				return
			}
			vnfd, err := ctl.TcaClient.GetVnfPkgmVnfd(pkgCnf.PID)
			if err != nil || vnfd == nil {
				glog.Error("Failed acquire VDU information for %v", pkgCnf.PID)
				return
			}
			// get linked repo, if caller provide repo that is not linked nothing to do
			reposId, err := ctl.TcaClient.LinkedRepositories(cloud.TenantID, defaultRepo)
			if err != nil {
				glog.Errorf("Failed acquire linked %v "+
					"repository to provider %v. Indicate a repo "+
					"linked to cloud provider.", defaultRepo, cloud.TenantID)
				return
			}
			ext, err := ctl.TcaClient.ExtensionQuery()
			if err != nil {
				glog.Error("Failed acquire extension information for %v", err)
				return
			}
			linkedRepos, err := ext.FindRepo(reposId)
			if err != nil || linkedRepos == nil {
				glog.Error("Failed acquire extension information for %v", pkgCnf.PID)
				return
			}
			nodePool, _, err := ctl.TcaClient.GetNamedClusterNodePools("edge")
			if err != nil || nodePool == nil {
				glog.Error("Failed acquire clusters node information for %v", pkgCnf.PID)
				return
			}
			pool, err := nodePool.GetPool("cellside01-pool")
			if err != nil {
				glog.Error("Failed acquire node pool information")
				return
			}

			vnfLcm, err := ctl.TcaClient.CnfVnfInstantiate(&request.CreateVnfLcm{
				VnfdId:                 pkgCnf.VnfdID,
				VnfInstanceName:        instanceName,
				VnfInstanceDescription: "",
			})
			if err != nil {
				glog.Errorf("Failed create instance information %v", err)
				return
			}

			glog.Infof("Creating vnf for vnfid id %v instance \n\t\t name %s cloud id %s", pkgCnf.VnfdID, instanceName, cloud.VimID)
			glog.Infof("Total number of VDUs %v", len(vnfd.Vdus))
			glog.Infof("Linked repo information %v", linkedRepos.InterfaceInfo.Url)

			var flavorName = "default"
			if len(vnfd.Vnf.Properties.FlavourId) > 0 {
				flavorName = vnfd.Vnf.Properties.FlavourId
			}

			for _, vdu := range vnfd.Vdus {
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
							Namespace: namespace,
							RepoURL:   defaultRepo,
							Username:  linkedRepos.AccessInfo.Username,
							Password:  linkedRepos.AccessInfo.Password,
							VduName:   vdu.VduId,
						}},
						DisableGrant:        true,
						IgnoreGrantFailure:  false,
						DisableAutoRollback: false,
					},
				}
				glog.Infof("Instantiating %v", vnfLcm.Id)
				err := ctl.TcaClient.CnfInstantiate(vnfLcm.Id, req)
				if err != nil {
					glog.Errorf("Failed create cnf instance information %v", err)
					return
				}
				io.PrettyPrint(req)
			}
		},
	}

	cmdCreate.Flags().BoolVar(&DisableGrant,
		"disable_grant", true,
		"disable grant validation")
	//
	cmdCreate.Flags().StringVarP(&namespace,
		"namespace", "n", "default",
		"cnf running instance name")
	//
	cmdCreate.Flags().StringVarP(&defaultRepo,
		"repo", "r", "https://repo.cnfdemo.io/chartrepo/library",
		"cnf running instance name")

	return cmdCreate
}

//{"flavourId":"default","vimConnectionInfo":[{"id":"vmware_8BF6253CE6E247018D909605A437B827",
//	"vimType":"",
//	"extra":{
//			"nodePoolId":"62fc287f-697f-4892-8ad3-847bad245115"
//			}
//	  }
//	],
//	"additionalParams":{"vduParams":[{"namespace":"test","repoUrl":"https://repo.cnfdemo.io/chartrepo/library",
//		"username":"admin","password":"31337Hax0rsMustDie","vduName":"mwc_app02"}],"disableGrant":true,"ignoreGrantFailure":false,"disableAutoRollback":false}}

func (ctl *TcaCtl) CmdCreate() *cobra.Command {
	// cnf instances

	var cmdCreate = &cobra.Command{
		Use:   "Create",
		Short: "Terminate CNF instance",
		Long:  `Terminate CNF instance, caller need to provide CNF Identifier.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	cmdCreate.AddCommand(ctl.CmdCreateCnf())
	return cmdCreate
}
