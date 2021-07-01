package testutil

import (
	ioutils "github.com/spyroot/tcactl/pkg/io"
	"io"
	"io/ioutil"
	"os"
)

// SpecTempReader  helper specs from string
// create temp file and write all spec to a file
// and return descriptor
func SpecTempReader(spec string) io.Reader {

	tmpFile, err := ioutil.TempFile("", "tcactltest")
	ioutils.CheckErr(err)

	// write to file,close it and read spec
	if _, err = tmpFile.Write([]byte(spec)); err != nil {
		ioutils.CheckErr(err)
	}

	if err := tmpFile.Close(); err != nil {
		ioutils.CheckErr(err)
	}

	fd, err := os.Open(tmpFile.Name())
	ioutils.CheckErr(err)

	return fd
}

// SpecTempFileName  helper specs from string
// create temp file and write all spec to a file
// and return temp file name
func SpecTempFileName(spec string) string {

	tmpFile, err := ioutil.TempFile("", "tcactltest")
	ioutils.CheckErr(err)

	// write to file,close it and read spec
	if _, err = tmpFile.Write([]byte(spec)); err != nil {
		ioutils.CheckErr(err)
	}

	if err := tmpFile.Close(); err != nil {
		ioutils.CheckErr(err)
	}

	return tmpFile.Name()
}
