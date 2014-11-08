package main

import (
	. "btrfs"
	"container/list"
	"testing"
)

func TestBtrfsFindDeviceByDevid(t *testing.T) {
	btrfsDevice1 := &BtrfsDevice{
		Devid: 1,
	}
	btrfsDevice2 := &BtrfsDevice{
		Devid: 2,
	}
	fsDevices := &BtrfsFsDevices{}
	fsDevices.Devices = list.New()
	fsDevices.Devices.PushBack(btrfsDevice2)
	fsDevices.Devices.PushBack(btrfsDevice1)
	devid := uint64(1)
	instance := 1
	btrfsDevice := btrfsFindDeviceByDevid(fsDevices, devid, instance)
	if btrfsDevice == nil {
		t.Errorf("btrfsDevice should not be nil")
	} else {
		if btrfsDevice.Devid != 1 {
			t.Errorf("btrfsDevice.Devid should be 1")
		}
	}

}
