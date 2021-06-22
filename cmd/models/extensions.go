// Package models
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
package models

const (
	// ExtensionNokiaSubType - extinction types for nokia CBAM
	ExtensionNokiaSubType = "Nokia-CBAM"

	// ExtensionHarborSubType  - extinction types for nokia Harbor
	ExtensionHarborSubType = "Harbor"

	// ExtensionTypeRepository extinction type
	ExtensionTypeRepository = "Repository"

	// ExtensionTypeSVNFM extension type for SVNFM
	ExtensionTypeSVNFM = "SVNFM"
)

// ExtensionTypes - Extension abstraction used by TCA
// type is type and under sub-type
type ExtensionTypes struct {
	Types []struct {
		SubTypes []string `json:"subTypes" yaml:"sub_types"`
		Uri      string   `json:"uri" yaml:"uri"`
		Type     string   `json:"type" yaml:"type"`
	} `json:"types"`
}
