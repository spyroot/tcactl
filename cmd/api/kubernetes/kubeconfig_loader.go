package kubernetes

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	osutil "github.com/spyroot/tcactl/pkg/os"
)

const (
	KUBECONFIG = "KUBECONFIG"
	KUBEDIR    = ".kube"
	KUBEFILE   = "config"
)

var (
	DefaultLoader FileReader = new(KubeconfigLoader)
)

type KubeconfigLoader struct {
}

type ConfigFile struct {
	*os.File
}

// Reset reset file position
func (k *ConfigFile) Reset() error {

	if err := k.Truncate(0); err != nil {
		return errors.Wrap(err, "failed to truncate file")
	}

	_, err := k.Seek(0, 0)
	return errors.Wrap(err, "failed to seek in file")
}

// getFiles() return list of kubeconfig files
func getFiles() ([]string, error) {

	var files []string

	// check first env
	if v := os.Getenv(KUBECONFIG); v != "" {
		list := filepath.SplitList(v)
		for _, s := range list {
			files = append(files, filepath.Join(s, KUBEDIR, KUBEFILE))
		}
		return files, nil
	}

	home := osutil.HomeDir()
	if home == "" {
		return files, errors.New("can't determine user home dir")
	}

	files = append(files, filepath.Join(home, KUBEDIR, KUBEFILE))
	return files, nil
}

// Load load config files
func (k *KubeconfigLoader) Read() ([]ReadWriteResetCloser, error) {

	files, err := getFiles()
	if err != nil {
		return nil, errors.Wrap(err, "cannot determine kubeconfig path")
	}

	var rwc []ReadWriteResetCloser

	for _, file := range files {
		f, err := os.OpenFile(file, os.O_RDWR, 0)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, errors.Wrap(err, "kubeconfig file not found")
			}
			return nil, errors.Wrap(err, "failed to open file")
		}

		rwc = append(rwc, &ConfigFile{f})
	}

	return rwc, nil
}
