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

package io

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"gopkg.in/yaml.v3"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

/**

 */
func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// CheckSocket
// Return if given IP:PORT bind open or not, network must
// must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
// No validation, caller must do that/*
func CheckSocket(hostPort string, proto string) bool {

	l, err := net.Listen(proto, hostPort)
	if l != nil {
		defer l.Close()
	}
	if err != nil {
		return false
	}
	return true
}

// ReadUint64 /**
func ReadUint64(data []byte) (ret uint64) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

// IsDir /**
func IsDir(dir string) (bool, error) {

	if len(dir) == 0 {
		return false, nil
	}

	src, err := os.Stat(dir)
	if err != nil {
		return false, err
	}

	if !src.IsDir() {
		return false, err
	}
	return true, nil
}

// FileExists /*
func FileExists(filename string) bool {
	stat, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !stat.IsDir()
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}

// PrettyPrint pretty printer
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// YamlPrinter Default Json printer
func YamlPrinter(t interface{}, isColor ...bool) error {
	b, err := yaml.Marshal(t)
	if err != nil {
		return err
	}

	if len(isColor) > 0 {
		fmt.Println(string(pretty.Color(b, nil)))
	} else {
		fmt.Println(string(b))
	}

	return nil
}

// PrettyString print string as pretty json
func PrettyString(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		return string(b)
	}
	return ""
}

func CheckErr(msg interface{}) {
	if msg != nil {
		fmt.Fprintln(os.Stderr, "Error:", msg)
		os.Exit(1)
	}
}

func PrintAndExit(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}
