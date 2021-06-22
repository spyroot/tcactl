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
	"github.com/google/uuid"
	"regexp"
	"strconv"
)

const (
	// numericRegexString match numeric value
	numericRegexString = "^[-+]?[0-9]+(?:\\.[0-9]+)?$"
	// numberRegexString match number
	numberRegexString = "^[0-9]+$"
)

var (
	// numericRegex must used to match numeric
	numericRegex = regexp.MustCompile(numericRegexString)
	// numberRegex must be used to match number
	numberRegex = regexp.MustCompile(numberRegexString)
)

// IsValidUUID check value in UUID format
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// panicIfNeeded if value incorrect.
func panicIfNeeded(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// asBool returns the parameter as a bool
// or panics if it can't convert
func asBool(param string) bool {

	i, err := strconv.ParseBool(param)
	panicIfNeeded(err)
	return i
}
