package btrfs

import ()

func btrfs_scan_one_device(
	fd int,
	path string,
	fs_devices_ret **Btrfs_fs_devices,
	total_devs *uint64,
	super_offset uint64,
	super_recover bool) int {
	//	struct btrfs_super_block *disk_super;
	//	char *buf;
	//	int ret;
	//	u64 devid;
	//
	//	buf = malloc(4096);
	//	if (!buf) {
	//		ret = -ENOMEM;
	//		goto error;
	//	}
	//	disk_super = (struct btrfs_super_block *)buf;
	disk_super := new(Btrfs_super_block)
	if ret := btrfs_read_dev_super(fd, disk_super, super_offset, super_recover); ret {
//		devid := disk_super.Dev_item.Devid
		if disk_super.Flags&BTRFS_SUPER_FLAG_METADUMP > 0 {
			*total_devs = 1
		} else {
			*total_devs = disk_super.Num_devices
		}
	}
	//	if (ret < 0) {
	//		ret = -EIO;
	//		goto error_brelse;
	//	}
	//	devid = btrfs_stack_device_id(&disk_super->dev_item);
	//	if (btrfs_super_flags(disk_super) & BTRFS_SUPER_FLAG_METADUMP)
	//		*total_devs = 1;
	//	else
	//		*total_devs = btrfs_super_num_devices(disk_super);
	//
	//	ret = device_list_add(path, disk_super, devid, fs_devices_ret);
	//
	//error_brelse:
	//	free(buf);
	//error:
	return 0
}
func device_list_add(path string,
	disk_super *Btrfs_super_block,
	devid uint64,
	fs_devices_ret **Btrfs_fs_devices) int {

	//	struct btrfs_device *device;
	//	struct btrfs_fs_devices *fs_devices;
	//	u64 found_transid = btrfs_super_generation(disk_super);
	//
	//	fs_devices = find_fsid(disk_super->fsid);
	//	if (!fs_devices) {
	//		fs_devices = kzalloc(sizeof(*fs_devices), GFP_NOFS);
	//		if (!fs_devices)
	//			return -ENOMEM;
	//		INIT_LIST_HEAD(&fs_devices->devices);
	//		list_add(&fs_devices->list, &fs_uuids);
	//		memcpy(fs_devices->fsid, disk_super->fsid, BTRFS_FSID_SIZE);
	//		fs_devices->latest_devid = devid;
	//		fs_devices->latest_trans = found_transid;
	//		fs_devices->lowest_devid = (u64)-1;
	//		device = NULL;
	//	} else {
	//		device = __find_device(&fs_devices->devices, devid,
	//				       disk_super->dev_item.uuid);
	//	}
	//	if (!device) {
	//		device = kzalloc(sizeof(*device), GFP_NOFS);
	//		if (!device) {
	//			/* we can safely leave the fs_devices entry around */
	//			return -ENOMEM;
	//		}
	//		device->fd = -1;
	//		device->devid = devid;
	//		memcpy(device->uuid, disk_super->dev_item.uuid,
	//		       BTRFS_UUID_SIZE);
	//		device->name = kstrdup(path, GFP_NOFS);
	//		if (!device->name) {
	//			kfree(device);
	//			return -ENOMEM;
	//		}
	//		device->label = kstrdup(disk_super->label, GFP_NOFS);
	//		if (!device->label) {
	//			kfree(device->name);
	//			kfree(device);
	//			return -ENOMEM;
	//		}
	//		device->total_devs = btrfs_super_num_devices(disk_super);
	//		device->super_bytes_used = btrfs_super_bytes_used(disk_super);
	//		device->total_bytes =
	//			btrfs_stack_device_total_bytes(&disk_super->dev_item);
	//		device->bytes_used =
	//			btrfs_stack_device_bytes_used(&disk_super->dev_item);
	//		list_add(&device->dev_list, &fs_devices->devices);
	//		device->fs_devices = fs_devices;
	//	} else if (!device->name || strcmp(device->name, path)) {
	//		char *name = strdup(path);
	//                if (!name)
	//                        return -ENOMEM;
	//                kfree(device->name);
	//                device->name = name;
	//        }
	//
	//
	//	if (found_transid > fs_devices->latest_trans) {
	//		fs_devices->latest_devid = devid;
	//		fs_devices->latest_trans = found_transid;
	//	}
	//	if (fs_devices->lowest_devid > devid) {
	//		fs_devices->lowest_devid = devid;
	//	}
	//	*fs_devices_ret = fs_devices;
	return 0
}
