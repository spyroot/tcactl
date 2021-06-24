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
	"strings"
	"time"
)

const (
	CliBlock       = "block"
	CliPool        = "pool"
	CliDisableGran = "grant"
	CliForce       = "force"
)

// BlockWaitStateChange blocs main thread of execution and wait specific state to change a state,
// it busy waiting.
func (ctl *TcaCtl) BlockWaitStateChange(instanceId string, waitFor string, maxRetry int) error {

	for i := 1; i < maxRetry; i++ {
		instance, err := ctl.TcaClient.GetRunningVnflcm(instanceId)
		if err != nil {
			glog.Error(err)
			return err
		}
		fmt.Printf("Current state %s waiting for %s\n", instance.InstantiationState, waitFor)
		fmt.Printf("Current LCM Operation status %s target state %s\n\n",
			instance.Metadata.LcmOperationState, instance.Metadata.LcmOperation)

		if strings.HasPrefix(instance.InstantiationState, waitFor) {
			break
		}

		time.Sleep(30 * time.Second)
	}

	return nil
}
