// Package api
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
	"strings"
)

const (

	// BlockMaxRetryTimer - default retry timer
	BlockMaxRetryTimer = 32

	// TaskPoolSeconds  default pool timer.
	TaskPoolSeconds = 10

	StateInstantiate  = "INSTANTIATE"
	StateCompleted    = "COMPLETED"
	StateStarting     = "STARTING"
	StateInstantiated = "INSTANTIATED"
	StateTerminated   = "TERMINATED"
	StateTerminate    = "TERMINATE"
	StateFailedTemp   = "FAILED_TEMP"

	StateNotInstantiated = "NOT_INSTANTIATED"

	DefaultMaxRetry = 32

	// TaskStateSuccess task state successes
	TaskStateSuccess = "SUCCESS"

	// TaskStateQueued task state queued by tca
	TaskStateQueued = "QUEUED"

	// TaskStateFailed task state failed
	TaskStateFailed = "FAILED"

	TaskTypeNodePoolCreation = "Node Pool Creation"
	TaskTypeInventoryUpdate  = "Inventory Update"
	TaskTypeUpdateNodePool   = "Update Node Pool"

	CLusterCsiNfs     = "nfs_client"
	ClusterCsiVsphere = "vsphere-csi"
)

type TcaTaskStateType int

// TODO move all state
const (
	// Instantiate CNF/VNF instantiated
	Instantiate TcaTaskStateType = 1 << iota
	// Terminated CNF/VNF terminated
	Terminated
	// Completed task competed
	Completed
)

// IsInState - abstract state change,  late if I move code to
// state as different representation we swap code here.
func IsInState(currentState string, predicate string) bool {
	return strings.Contains(currentState, predicate)
}
