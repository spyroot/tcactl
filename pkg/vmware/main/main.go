package main

import (
	"context"
	"fmt"
	"github.com/spyroot/tcactl/pkg/vmware/vc"
	"os"
)

func main() {

	src := "core-13.0.iso"
	dst := "/ISO/core-13.0.iso"
	datastoreName := "vsanDatastore"

	ctx := context.TODO()
	c, err := vc.Connect(ctx, os.Getenv("VC_HOSTNAME"), os.Getenv("VC_USERNAME"), os.Getenv("VC_PASSWORD"))
	if err != nil {
		fmt.Println("error", err)
		return
	}

	vcRest := vc.VSphereRest{Ctl: c.Client}
	err = vcRest.Upload(ctx, datastoreName, src, dst)
	if err != nil {
		fmt.Print("Error", err)
		return
	}
}
