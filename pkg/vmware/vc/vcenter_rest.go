package vc

import (
	"context"
	"fmt"
	ioutils "github.com/spyroot/tcactl/pkg/io"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/govc/flags"
	_ "github.com/vmware/govmomi/property"
	_ "github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	_ "github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
	_ "github.com/vmware/govmomi/vim25/types"
	"net"
)

type Runner interface {
	Run(ctx context.Context, c *vim25.Client, args map[string]interface{}) error
}

// vc interface
type upload struct {
	*flags.OutputFlag
	*flags.DatastoreFlag
}

//type HostDisk struct {
//	Capacity              HostDiskDimensionsLba `xml:"capacity"`
//	DevicePath            string                `xml:"devicePath"`
//	Ssd                   *bool                 `xml:"ssd"`
//	LocalDisk             *bool                 `xml:"localDisk"`
//	PhysicalLocation      []string              `xml:"physicalLocation,omitempty"`
//	EmulatedDIXDIFEnabled *bool                 `xml:"emulatedDIXDIFEnabled"`
//	VsanDiskInfo          *VsanHostVsanDiskInfo `xml:"vsanDiskInfo,omitempty"`
//	ScsiDiskType          string                `xml:"scsiDiskType,omitempty"`
//}

//type HostScsiDisk struct {
//	ScsiLun
//
//	Capacity              HostDiskDimensionsLba `xml:"capacity"`
//	DevicePath            string                `xml:"devicePath"`
//	Ssd                   *bool                 `xml:"ssd"`
//	LocalDisk             *bool                 `xml:"localDisk"`
//	PhysicalLocation      []string              `xml:"physicalLocation,omitempty"`
//	EmulatedDIXDIFEnabled *bool                 `xml:"emulatedDIXDIFEnabled"`
//	VsanDiskInfo          *VsanHostVsanDiskInfo `xml:"vsanDiskInfo,omitempty"`
//	ScsiDiskType          string                `xml:"scsiDiskType,omitempty"`
//}

// VsphereHost A datastore host
type VsphereHost struct {
	Reference     string               `json:"Reference" yaml:"Reference"`
	ManagementIPs []net.IP             `json:"ManagementAddress" yaml:"ManagementAddress"`
	Disks         []types.HostScsiDisk `json:"Disks" yaml:"Disks"`
}

// VsphereDatastore vsphere datastore
// Example /Datacenter is path
type VsphereDatastore struct {
	DatacenterPath string                                   `json:"DatacenterPath" yaml:"DatacenterPath"`
	InventoryPath  string                                   `json:"InventoryPath" yaml:"InventoryPath"`
	Name           string                                   `json:"Name" yaml:"Name"`
	Hosts          map[string]VsphereHost                   `json:"Hosts" yaml:"Hosts"`
	Type           types.HostFileSystemVolumeFileSystemType `json:"Type" yaml:"Type"`
}

type VsphereDatastores struct {
	Datastores map[string]VsphereDatastore
}

// VSphereRest main rest api struct
type VSphereRest struct {
	Ctl    *vim25.Client
	upload upload
}

// Upload upload a file to target datastore.
func (rest *VSphereRest) Upload(ctx context.Context, datastoreName string, src string, dst string) error {

	if len(datastoreName) == 0 {
		return fmt.Errorf("empty datastore name")
	}
	if len(src) == 0 {
		return fmt.Errorf("empty source file name")
	}
	if !ioutils.FileExists(src) {
		return fmt.Errorf("check path to a source file %s", src)
	}
	if len(dst) == 0 {
		return fmt.Errorf("empty destination path")
	}
	cmdArgs := make(map[string]interface{})
	cmdArgs["datastore"] = datastoreName
	cmdArgs["src"] = src
	cmdArgs["dst"] = dst
	return rest.upload.Run(ctx, rest.Ctl, cmdArgs)
}

// GetDatastores upload a file to target datastore.
func (rest *VSphereRest) GetDatastores(ctx context.Context, path string) (*VsphereDatastores, error) {

	finder := find.NewFinder(rest.Ctl)
	if finder == nil {
		return nil, fmt.Errorf("failed retrieve vc fidner")
	}

	_query := path
	if len(path) == 0 {
		_query = "*"
	}

	ds, err := finder.DatastoreList(ctx, _query)
	if err != nil {
		return nil, err
	}

	if ds == nil {
		return nil, fmt.Errorf("datastoreList nil")
	}

	vcdss := VsphereDatastores{}
	vcdss.Datastores = make(map[string]VsphereDatastore)

	for i := range ds {
		vcds := VsphereDatastore{}
		vcds.DatacenterPath = ds[i].DatacenterPath
		vcds.InventoryPath = ds[i].InventoryPath
		vcds.Name = ds[i].Name()
		h, err := ds[i].AttachedHosts(ctx)
		if err != nil {
			fmt.Println("Error can't get attached host list for ", i, ds[i].Name())
		}
		t, err := ds[i].Type(ctx)
		if err != nil {
			fmt.Println("Error can't get type for ", i, ds[i].Name())
		}
		if len(t) > 0 {
			vcds.Type = t
		}
		vcds.Hosts = make(map[string]VsphereHost)
		if h != nil && len(h) > 0 {
			for j := range h {
				_host := h[j].Reference().Value
				_addr, err := h[j].ManagementIPs(ctx)

				mg := h[j].ConfigManager()
				ds, _ := mg.DatastoreSystem(ctx)
				disk, _ := ds.QueryAvailableDisksForVmfs(ctx)

				if len(_host) > 0 && err == nil {
					vcds.Hosts[_host] = VsphereHost{
						Reference:     _host,
						ManagementIPs: _addr,
						Disks:         disk,
					}
				}
			}
			vcdss.Datastores[ds[i].Reference().Value] = vcds
		}
	}

	return &vcdss, nil
}

// Run run upload cmd
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
