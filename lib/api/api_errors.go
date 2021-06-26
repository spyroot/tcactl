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

// InvalidVimFormat error must returned if client supplied incorrect format for vim ID
type InvalidVimFormat struct {
	errMsg string
}

func (m *InvalidVimFormat) Error() string {
	return "vim id format " + m.errMsg + " invalid. Example vmware_FB40D3DE2967483FBF9033B451DC7571"
}

// InvalidTaskId error must returned if client supplied incorrect task id
type InvalidTaskId struct {
	errMsg string
}

func (m *InvalidTaskId) Error() string {
	return "task id" + m.errMsg + " invalid. Example 9411f70f-d24d-4842-ab56-b7214d39d1b1"
}

// InvalidSpec error must returned if client supplied incorrect yaml or json initialSpec
type InvalidSpec struct {
	errMsg string
}

func (m *InvalidSpec) Error() string {
	return "Spec is in invalid format. " + m.errMsg
}
