package kubernetes

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io"
)

//
type ReadWriteResetCloser interface {
	io.ReadWriteCloser
	Reset() error
}

//
type FileReader interface {
	Read() ([]ReadWriteResetCloser, error)
}

//
type KubeconfigFileReaderWriter struct {
	Kubeconfig KubeconfigStruct
	files      ReadWriteResetCloser
	reader     FileReader
}

func (k *KubeconfigFileReaderWriter) WithLoader(l FileReader) *KubeconfigFileReaderWriter {
	k.reader = l
	return k
}

// Close close file
func (k *KubeconfigFileReaderWriter) Close() error {
	if k.files == nil {
		return nil
	}
	return k.files.Close()
}

// Parse - parse kubeconfig
func (k *KubeconfigFileReaderWriter) Parse() error {

	if k == nil {
		return errors.New("nil instance")
	}

	if k.reader == nil {
		return errors.New("nil reader")
	}

	files, err := k.reader.Read()
	if err != nil {
		return errors.Wrap(err, "failed to load")
	}

	currentFile := files[0]
	if err := yaml.NewDecoder(currentFile).Decode(&k.Kubeconfig); err != nil {
		return errors.Wrap(err, "failed to decode")
	}

	return nil
}

// Bytes return bytes
func (k *KubeconfigFileReaderWriter) Bytes() ([]byte, error) {

	if k == nil {
		return nil, errors.New("nil instance")
	}

	return yaml.Marshal(k.Kubeconfig)
}

// Save kubeconfig
func (k *KubeconfigFileReaderWriter) Save() error {

	if k == nil {
		return errors.New("nil instance")
	}

	if err := k.files.Reset(); err != nil {
		return errors.Wrap(err, "failed to reset file")
	}

	enc := yaml.NewEncoder(k.files)
	enc.SetIndent(0)

	return enc.Encode(k.Kubeconfig)
}
