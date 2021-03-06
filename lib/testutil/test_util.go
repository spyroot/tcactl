package testutil

import (
	ioutils "github.com/spyroot/tcactl/pkg/io"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"runtime"
	"time"
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

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//RunOnRootFolder Changes dir
func RunOnRootFolder() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		ioutils.CheckErr(err)
	}

	wd, err := os.Getwd()
	ioutils.CheckErr(err)

	return wd
}
