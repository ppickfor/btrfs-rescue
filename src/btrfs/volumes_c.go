package btrfs

import (
	"bytes"
	"container/list"
	)

func btrfsScanOneDevice(
	fd int,
	path string,
	fsDevicesRet **BtrfsFsDevices,
	totalDevs *uint64,
	superOffset uint64,
	superRecover bool) int {
	//	struct btrfsSuperBlock *diskSuper;
	//	char *buf;
	//	int ret;
	//	u64 devid;
	//
	//	buf = malloc(4096);
	//	if (!buf) {
	//		ret = -ENOMEM;
	//		goto error;
	//	}
	//	diskSuper = (struct btrfsSuperBlock *)buf;
	diskSuper := new(BtrfsSuperBlock)
	if ret := btrfsReadDevSuper(fd, diskSuper, superOffset, superRecover); ret {
		//		devid := diskSuper.DevItem.Devid
		if diskSuper.Flags&BTRFS_SUPER_FLAG_METADUMP > 0 {
			*totalDevs = 1
		} else {
			*totalDevs = diskSuper.NumDevices
		}
	}
	//	if (ret < 0) {
	//		ret = -EIO;
	//		goto errorBrelse;
	//	}
	//	devid = btrfsStackDeviceId(&diskSuper->devItem);
	//	if (btrfsSuperFlags(diskSuper) & BTRFS_SUPER_FLAG_METADUMP)
	//		*totalDevs = 1;
	//	else
	//		*totalDevs = btrfsSuperNumDevices(diskSuper);
	//
	//	ret = deviceListAdd(path, diskSuper, devid, fsDevicesRet);
	//
	//errorBrelse:
	//	free(buf);
	//error:
	return 0
}
func deviceListAdd(path string,
	diskSuper *BtrfsSuperBlock,
	devid uint64,
	fsDevicesRet **BtrfsFsDevices) int {

	//	struct btrfsDevice *device;
	//	struct btrfsFsDevices *fsDevices;
	//	u64 foundTransid = btrfsSuperGeneration(diskSuper);
	//
	//	fsDevices = findFsid(diskSuper->fsid);
	//	if (!fsDevices) {
	//		fsDevices = kzalloc(sizeof(*fsDevices), GFP_NOFS);
	//		if (!fsDevices)
	//			return -ENOMEM;
	//		INIT_LIST_HEAD(&fsDevices->devices);
	//		listAdd(&fsDevices->list, &fsUuids);
	//		memcpy(fsDevices->fsid, diskSuper->fsid, BTRFS_FSID_SIZE);
	//		fsDevices->latestDevid = devid;
	//		fsDevices->latestTrans = foundTransid;
	//		fsDevices->lowestDevid = (u64)-1;
	//		device = NULL;
	//	} else {
	//		device = _FindDevice(&fsDevices->devices, devid,
	//				       diskSuper->devItem.uuid);
	//	}
	//	if (!device) {
	//		device = kzalloc(sizeof(*device), GFP_NOFS);
	//		if (!device) {
	//			/* we can safely leave the fsDevices entry around */
	//			return -ENOMEM;
	//		}
	//		device->fd = -1;
	//		device->devid = devid;
	//		memcpy(device->uuid, diskSuper->devItem.uuid,
	//		       BTRFS_UUID_SIZE);
	//		device->name = kstrdup(path, GFP_NOFS);
	//		if (!device->name) {
	//			kfree(device);
	//			return -ENOMEM;
	//		}
	//		device->label = kstrdup(diskSuper->label, GFP_NOFS);
	//		if (!device->label) {
	//			kfree(device->name);
	//			kfree(device);
	//			return -ENOMEM;
	//		}
	//		device->totalDevs = btrfsSuperNumDevices(diskSuper);
	//		device->superBytesUsed = btrfsSuperBytesUsed(diskSuper);
	//		device->totalBytes =
	//			btrfsStackDeviceTotalBytes(&diskSuper->devItem);
	//		device->bytesUsed =
	//			btrfsStackDeviceBytesUsed(&diskSuper->devItem);
	//		listAdd(&device->devList, &fsDevices->devices);
	//		device->fsDevices = fsDevices;
	//	} else if (!device->name || strcmp(device->name, path)) {
	//		char *name = strdup(path);
	//                if (!name)
	//                        return -ENOMEM;
	//                kfree(device->name);
	//                device->name = name;
	//        }
	//
	//
	//	if (foundTransid > fsDevices->latestTrans) {
	//		fsDevices->latestDevid = devid;
	//		fsDevices->latestTrans = foundTransid;
	//	}
	//	if (fsDevices->lowestDevid > devid) {
	//		fsDevices->lowestDevid = devid;
	//	}
	//	*fsDevicesRet = fsDevices;
	return 0
}

func _FindDevice(devices *list.List, devId uint64, uuid []byte) *BtrfsDevice {
	for i := devices.Front(); i != nil && i.Value != nil; i = i.Next() {
		device := i.Value.(*BtrfsDevice)
		if device.Devid == devId && bytes.Equal(device.Uuid[:], uuid) {
			return device
		}
	}
	return nil
}
func btrfsFindDevice(root *BtrfsRoot, devId uint64, uuid []byte, fsId *[]byte) *BtrfsDevice {
	curDevices := root.FsInfo.FsDevices
	for curDevices != nil {
		if fsId == nil || bytes.Equal(curDevices.Fsid[:], *fsId) {
			device := _FindDevice(curDevices.Devices, devId, uuid)
			return device

		}
		curDevices = curDevices.Seed
	}
	return nil
}
func buildDeviceMapByChunkRecord(root *BtrfsRoot, chunk *ChunkRecord) {

	mapTree := &root.FsInfo.MappingTree
	numStripes := chunk.NumStripes
	Map := &MapLookup{
		CacheExtent: CacheExtent{
			Start: chunk.Offset,
			Size:  chunk.Length,
		},
		NumStripes: int32(numStripes),
		IoWidth:    int32(chunk.IoWidth),
		IoAlign:    int32(chunk.IoAlign),
		SectorSize: int32(chunk.SectorSize),
		StripeLen:  int32(chunk.StripeLen),
		Type:       chunk.TypeFlags,
		SubStripes: int32(chunk.SubStripes),
	}
	for i, stripe := range chunk.Stripes {
		devId := stripe.Devid
		uuid := stripe.Uuid
		Map.Stripes[i].Physical = stripe.Offset
		Map.Stripes[i].Dev = btrfsFindDevice(root, devId, uuid[:], nil)
	}
	mapTree.Tree.InsertNoReplace(Map)
}

func BuildDeviceMapsByChunkRecords(rc *RecoverControl, root *BtrfsRoot) {
	for i := rc.GoodChunks.Front(); i != nil && i.Value != nil; i = i.Next() {
		chunk := i.Value.(*ChunkRecord)
		buildDeviceMapByChunkRecord(root, chunk)
	}
}