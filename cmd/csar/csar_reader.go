package csar

import (
	"archive/zip"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/cmd/models"
	"github.com/spyroot/tcactl/pkg/io"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	SpecNfd = "NFD.yaml"
)

// YamlParser - parser callback
type YamlParser func(path string, substitution map[string]string) error

// NfdYamlPropertyTransformer - substitution callback
func NfdYamlPropertyTransformer(file string, substitution map[string]string) error {

	// read file to a buffer
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// parse to csar struct
	glog.Infof("Parsing target file. %v", file)
	var tosca models.CSAR
	err = yaml.Unmarshal(buffer, &tosca)
	if err != nil {
		return err
	}

	glog.Infof("Applying substitution map to a file %v", file)
	nodeTemplate := tosca.TopologyTemplate.NodeTemplates

	// each key in substitution is key in csar
	for sk, val := range substitution {
		for _, node := range nodeTemplate {
			node.Properties.UpdateField(sk, val)
		}
	}

	//
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(&tosca)
	if err != nil {
		return err
	}

	return nil
}

// ApplyTransformation adjusts yaml file based on substitution map
// It unzip csar file,
// find target yaml file and apply transformation function
// Compress csar back as new csar file.
func ApplyTransformation(zipFile string, fileName string,
	parser YamlParser, substitution map[string]string) (string, error) {

	dirName, err := ioutil.TempDir("", "tosca")
	if err != nil {
		return "", err
	}

	files, err := io.Unzip(zipFile, dirName)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if filepath.Base(file) == fileName {
			err := parser(file, substitution)
			if err != nil {
				return "", err
			}
			break
		}
	}

	newFileName := zipFile + ".new.csar"
	// compress to new csar
	err = io.ZipDir(dirName+"/", newFileName)
	if err != nil {
		return "", err
	}

	return newFileName, nil
}

// read zip reader
func read(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// Reader read file
//  Example:
// 		var topology models.CSAR
//  	find inside a zip NFD parse it and return
//		t, err := csar.Reader("/tests/smokeping-cnf.csar", "NFD.yaml", topology)
//		b, err := yaml.Marshal(&t)
//		fmt.Println(string(b))
func Reader(fileName string, targetFile string, topology interface{}) (interface{}, error) {

	zipReader, err := zip.OpenReader(fileName)
	if err != nil {
		return nil, err
	}

	// Read  the files from zip archive
	for _, zipFile := range zipReader.File {
		if targetFile == filepath.Base(zipFile.Name) {
			unzippedBytes, err := read(zipFile)
			if err != nil {
				log.Println(err)
				continue
			}

			err = yaml.Unmarshal(unzippedBytes, &topology)
			if err != nil {
				glog.Warningf("error during unmarshalling %v", err)
				continue
			}

			return &topology, nil
		}
	}

	return nil, nil
}
