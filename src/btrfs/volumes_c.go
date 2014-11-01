package btrfs

import ()

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
