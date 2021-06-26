package api

import (
	"github.com/spyroot/tcactl/lib/client"
	"github.com/spyroot/tcactl/pkg/io"
	"os"
)

func getClient() *client.RestClient {
	r, err := client.NewRestClient(os.Getenv("TCA_URL"),
		true,
		os.Getenv("TCA_USERNAME"),
		os.Getenv("TCA_PASSWORD"))
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

	return r
}

var (
	rest = getAuthenticatedClient()
)

//harbor = &client.RestClient{
//	BaseURL:               os.Getenv("HARBOR_URL"),
//	ApiKey:                "",
//	IsDebug:               true,
//	Username:              os.Getenv("HARBOR_USERNAME"),
//	Password:              os.Getenv("HARBOR_PASSWORD"),
//	SkipSsl:               true,
//}
