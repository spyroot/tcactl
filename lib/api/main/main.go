package main

import (
	"fmt"
	"github.com/spyroot/tcactl/lib/api/kubernetes"
	_ "github.com/spyroot/tcactl/lib/api/kubernetes"
	"github.com/spyroot/tcactl/pkg/io"
)

func main() {
	kubeconfig := new(kubernetes.KubeconfigFileReaderWriter).WithLoader(kubernetes.DefaultLoader)

	err := kubeconfig.Parse()
	if err != nil {
		fmt.Println("error ", err)
		return
	}

	fmt.Println(kubeconfig.Kubeconfig)
	io.PrettyPrint(kubeconfig.Kubeconfig)
}
