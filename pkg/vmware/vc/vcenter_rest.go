package vc

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

type Runner interface {
	Run(ctx context.Context, c *vim25.Client, args map[string]interface{}) error
}

// vc interface
type upload struct {
	*flags.OutputFlag
	*flags.DatastoreFlag
}

type VSphereRest struct {
	Ctl    *vim25.Client
	upload upload
}

func (rest *VSphereRest) Upload(ctx context.Context, datastoreName string, src string, dst string) error {

	cmdArgs := make(map[string]interface{})
	cmdArgs["datastore"] = datastoreName
	cmdArgs["src"] = src
	cmdArgs["dst"] = dst
	return rest.upload.Run(ctx, rest.Ctl, cmdArgs)
}

// stringify take key and
func stringify(k string, args map[string]interface{}) (string, error) {

	var strVal string
	var ok bool

	if x, found := args[k]; found {
		if strVal, ok = x.(string); !ok {
			return "", fmt.Errorf("type mistmatch for key %s", k)
		}
	} else {
		return "", fmt.Errorf("key %s not found ", k)
	}

	return strVal, nil
}

// Find key in interface and return value as interface
func Find(obj interface{}, key string) (interface{}, bool) {

	//if the argument is not a map, ignore it
	mobj, ok := obj.(map[string]interface{})
	if !ok {
		return nil, false
	}

	for k, v := range mobj {
		if k == key {
			return v, true
		}
		if m, ok := v.(map[string]interface{}); ok {
			if res, ok := Find(m, key); ok {
				return res, true
			}
		}
		if va, ok := v.([]interface{}); ok {
			for _, a := range va {
				if res, ok := Find(a, key); ok {
					return res, true
				}
			}
		}
	}

	return nil, false
}

func (cmd *upload) Run(ctx context.Context, c *vim25.Client, args map[string]interface{}) error {

	finder := find.NewFinder(c)

	var datastore, err = stringify("datastore", args)
	if err != nil {
		return err
	}

	src, err := stringify("src", args)
	if err != nil {
		return err
	}

	dst, err := stringify("dst", args)
	if err != nil {
		return err
	}

	ds, err := finder.Datastore(ctx, datastore)
	if err != nil {
		return err
	}

	p := soap.DefaultUpload
	return ds.UploadFile(ctx, src, dst, &p)
}
