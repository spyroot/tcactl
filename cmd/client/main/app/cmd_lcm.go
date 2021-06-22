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
	"github.com/spyroot/hestia/cmd/client/main/app/ui"
	"github.com/spyroot/hestia/cmd/client/response"
	"strings"
)

// CmdGetInstances Get CNF/VNF instances
func (ctl *TcaCtl) CmdGetInstances() *cobra.Command {

	var (
		_defaultPrinter = ctl.Printer
		_defaultStyler  = ctl.DefaultStyle

		_defaultFilter string
		_instanceID    string
		_outputFilter  string
	)

	var cmdCnfInstance = &cobra.Command{
		Use:     "cnfi",
		Short:   "Return cnf instance or all instances",
		Long:    `Returns list of cnf instances or instance if -i id provided.`,
		Example: "tcactl get cnfi -o json --filter \"{eq,id,5c11bd9c-085d-4913-a453-572457ddffe2}\"",
		Run: func(cmd *cobra.Command, args []string) {

			var (
				err            error
				genericRespond interface{}
			)

			// global output type
			_defaultPrinter = ctl.RootCmd.PersistentFlags().Lookup(FlagOutput).Value.String()

			// swap filter if output filter required
			if len(_outputFilter) > 0 {
				outputFields := strings.Split(_outputFilter, ",")
				_defaultPrinter = FilteredOutFilter
				_defaultStyler = ui.NewFilteredOutputStyler(outputFields)
			}

			_defaultStyler.SetColor(ctl.IsColorTerm)
			_defaultStyler.SetWide(ctl.IsWideTerm)

			// rest call
			if len(args) > 0 {
				genericRespond, err = ctl.TcaClient.GetVnflcm(_defaultFilter, args[0])
			} else {
				genericRespond, err = ctl.TcaClient.GetVnflcm(_defaultFilter)
			}
			if err != nil {
				glog.Error(err)
				return
			}

			// for extension request we route to correct printer
			cnfsExt, ok := genericRespond.(*response.CnfsExtended)
			if ok {
				if printer, ok := ctl.CnfInstanceExtendedPrinters[_defaultPrinter]; ok {
					printer(cnfsExt, _defaultStyler)
				}
				return
			}

			// for regular request we route to correct printer
			cnfsReg, ok := genericRespond.(*response.Cnfs)
			if ok {
				if printer, ok := ctl.CnfInstancePrinters[_defaultPrinter]; ok {
					printer(cnfsReg, _defaultStyler)
				}
			}
		},
	}

	//
	cmdCnfInstance.Flags().StringVarP(&_instanceID,
		"package_id", "i", "", "VNF package id")

	//
	cmdCnfInstance.Flags().StringVar(&_defaultFilter,
		"filter", "",
		"filter for query example, filter by id --filter \"{eq,id,5c11bd9c-085d-4913-a453-572457ddffe2}\"")

	return cmdCnfInstance
}
