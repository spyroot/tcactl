package api

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/pkg/io"
	"os"
)

func SetLoggingFlags() {
	fmt.Println("TestFlag testing:")

	err := flag.Set("alsologtostderr", "true")
	if err != nil {
		io.CheckErr(err)
	}
	//err = flag.Set("log_dir", ".")
	//if err != nil {
	//	io.CheckErr(err)
	//}
	err = flag.Set("v", "3")
	if err != nil {
		io.CheckErr(err)
	}

	flag.Parse()
	glog.Info("Logging configured")
}

// getClient() return tca client for unit testing
func getClient() *client.RestClient {
	tcaUrl := os.Getenv("TCA_URL")
	if len(tcaUrl) == 0 {
		io.PrintAndExit("TCA_URL not set")
	}
	tcaUsername := os.Getenv("TCA_USERNAME")
	if len(tcaUrl) == 0 {
		io.PrintAndExit("TCA_USERNAME not set")
	}
	tcaPassword := os.Getenv("TCA_PASSWORD")
	if len(tcaUrl) == 0 {
		io.PrintAndExit("TCA_PASSWORD not set")
	}
	r, err := client.NewRestClient(tcaUrl,
		true,
		tcaUsername,
		tcaPassword)

	fmt.Printf("Using tca endpoint %s and userame %s\n", tcaUrl, tcaPassword)
	if err != nil {
		io.CheckErr(err)
	}

	return r
}

func getAuthenticatedClient() *client.RestClient {
	r := getClient()

	ok, err := r.GetAuthorization()
	if err != nil {
		io.CheckErr(err)
	}

	if !ok {
		io.PrintAndExit("failed authenticated")
	}

	r.SetTrace(false)

	return r
}

var (
	rest = getAuthenticatedClient()
)

//harbor = &client.RestClient{
//	BaseURL:               os.Getenv("HARBOR_URL"),
//	apiKey:                "",
//	IsDebug:               true,
//	Username:              os.Getenv("HARBOR_USERNAME"),
//	Password:              os.Getenv("HARBOR_PASSWORD"),
//	SkipSsl:               true,
//}
