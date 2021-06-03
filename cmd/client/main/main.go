package main

import (
"flag"
"fmt"
"os"

"github.com/go-resty/resty/v2"
)




// print console usage
func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: example -stderrthreshold=[INFO|WARNING|FATAL] -log_dir=[string]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

//Inits logger
func init() {
	flag.Usage = usage
	_ = flag.Set("logtostderr", "true")
	_ = flag.Set("stderrthreshold", "WARNING")
	_ = flag.Set("v", "2")
	flag.Parse()

}

func Get_app() {

	// Create a resty client
	client := resty.New()

	resp, err := client.R().Get("https://tca-vip03.cnfdemo.io/admin/hybridity/api/authz/rbac/privileges")

	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Body: %v", resp)
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
}

// @title MWC proto API
// @version 1.0
// @description MWC proto API server.
// @termsOfService http://swagger.io/terms/
func main() {

	// default file
	var configFile = "config.yaml"

	if len(os.Args) > 1 {
		argList := os.Args[1:]
		configFile = argList[0]
	}

	fmt.Println(configFile)
	Get_app()

	//if callisto, err := server.NewCallisto(configFile); err == nil {
	//	done := make(chan bool, 1)
	//	callisto.Run(done)
	//}




}
