package main

import (
	. "btrfs"
	"container/list"
	"testing"
)

func TestBtrfsVerifyDeviceExtents(t *testing.T) {
	bg := &BlockGroupRecord{
		Flags:  0,
		Offset: 1024,
	}
	devExt1 := &DeviceExtentRecord{
		Length: 1024,
	}
	devExt2 := &DeviceExtentRecord{
		Length: 1024,
	}
	devExts := list.New()
	devExts.PushBack(devExt1)
	devExts.PushBack(devExt2)
	nDevExts := uint16(1)
	x := btrfsVerifyDeviceExtents(bg, devExts, nDevExts)
	if x == false {
		t.Errorf("btrfsVerifyDeviceExtents should be true")
	}
	bg.Flags = BTRFS_BLOCK_GROUP_RAID1
	nDevExts = 2
	x = btrfsVerifyDeviceExtents(bg, devExts, nDevExts)
	if x == false {
		t.Errorf("btrfsVerifyDeviceExtents should be true")
	}
	bg.Flags = BTRFS_BLOCK_GROUP_RAID1
	nDevExts = 2

	x = btrfsVerifyDeviceExtents(bg, devExts, nDevExts)
	if x == false {
		t.Errorf("btrfsVerifyDeviceExtents should be true")
	}
}
