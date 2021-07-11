package specs

import (
	"github.com/spyroot/tcactl/lib/testutil"
	"github.com/spyroot/tcactl/pkg/io"
	"os"
)

//GetTestAssetsDir return dir where test assets are
//it call RunOnRootFolder that change dir
func GetTestAssetsDir() string {

	wd := testutil.RunOnRootFolder()
	wd = wd + "/test_assets"

	_, err := os.Stat(wd)
	if os.IsNotExist(err) {
		io.CheckErr(err)
	}

	return wd
}
