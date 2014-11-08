package main

import (
	. "btrfs"
	"testing"
)

func TestIsExtentRecordInDeviceExtent(t *testing.T) {
	mirror := uint32(0)
	//	er.Devices[i].Devid == dext.Objectid &&
	//			er.Offsets[i] >= dext.Offset &&
	//			er.Offsets[i] < dext.Offset+dext.Length {

	er := &ExtentRecord{
		Nmirrors: 2,
		Devices: [BTRFS_MAX_MIRRORS]*BtrfsDevice{
			&BtrfsDevice{
				Devid: 2,
			},
			&BtrfsDevice{
				Devid: 1,
			},
			&BtrfsDevice{
				Devid: 2,
			}},
		Offsets: [BTRFS_MAX_MIRRORS]uint64{
			0,
			2048,
			0,
		},
	}
	dExt := &DeviceExtentRecord{
		Objectid: 1,
		Offset:   1024,
		Length:   2048,
	}
	x := isExtentRecordInDeviceExtent(er, dExt, &mirror)
	if x == false {
		t.Errorf("isExtentRecordInDeviceExtent should be true is %v", x)
	}
	if mirror != 1 {
		t.Errorf("isExtentRecordInDeviceExtent mirror should be 1, is %v", mirror)
	}

}
