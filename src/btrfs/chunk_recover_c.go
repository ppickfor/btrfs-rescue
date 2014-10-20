package btrfs

import (
	"log"
	"os"
	"syscall"
	"unsafe"
)

func is_super_block_address(offset uint64) bool {
	var i int

	for i = 0; i < BTRFS_SUPER_MIRROR_MAX; i++ {
		if offset == btrfs_sb_offset(i) {
			return true
		}
	}
	return false
}

func scan_one_device(dev_scan *Device_scan) bool {
	var (
		buf    *Extent_buffer
		bytenr uint64
		ret    bool = false
		//struct device_scan *dev_scan = (struct device_scan *)dev_scan_struct;
		rc     *Recover_control = dev_scan.Rc
		device *Btrfs_device    = dev_scan.Dev
		fd     int              = dev_scan.Fd
		//		oldtype int
		h *Btrfs_header
	)
	buf = new(Extent_buffer)
	buf.Len = uint64(rc.Leafsize)
	buf.Data = make([]byte, rc.Leafsize)
loop:
	for bytenr = 0; ; bytenr += rc.Sectorsize {
		ret = false
		if is_super_block_address(bytenr) {
			bytenr += rc.Sectorsize
		}
		n, err := syscall.Pread(fd, buf.Data, int64(bytenr))
		if err != nil {
			log.Fatalln(os.NewSyscallError("pread64", err))
		}
		if n < int(rc.Leafsize) {
			break loop
		}
		h = (*Btrfs_header)(unsafe.Pointer(&buf.Data))
		if rc.Fs_devices.Fsid != h.Fsid || verify_tree_block_csum_silent(buf, uint16(rc.Csum_size)) {
			continue loop
		}
		rc.Rc_lock.Lock()
		ret = process_extent_buffer(&rc.Eb_cache, buf, device, bytenr)
		rc.Rc_lock.Unlock()
		if !ret {
			break loop
		}
		if h.Level != 0 {
			switch h.Owner {
			case BTRFS_EXTENT_TREE_OBJECTID, BTRFS_DEV_TREE_OBJECTID:
				/* different tree use different generation */
				if h.Generation > rc.Generation {
					continue loop
				}
				if ret = extract_metadata_record(rc, buf); !ret {
					break loop
				}
			case BTRFS_CHUNK_TREE_OBJECTID:
				if h.Generation > rc.Chunk_root_generation {
					continue loop
				}
				if ret = extract_metadata_record(rc, buf); !ret {
					break loop
				}
			}
		}
	}

	//	close(fd);
	//	free(buf);
	return ret
}

func process_extent_buffer(eb_cache *Cache_tree,
	eb *Extent_buffer,
	device *Btrfs_device, offset uint64) bool {
	//	struct extent_record *rec;
	//	struct extent_record *exist;
	//	struct cache_extent *cache;
	//	int ret = 0;
	//
	//	rec = btrfs_new_extent_record(eb);
	//	if (!rec->cache.size)
	//		goto free_out;
	//again:
	//	cache = lookup_cache_extent(eb_cache,
	//				    rec->cache.start,
	//				    rec->cache.size);
	//	if (cache) {
	//		exist = container_of(cache, struct extent_record, cache);
	//
	//		if (exist->generation > rec->generation)
	//			goto free_out;
	//		if (exist->generation == rec->generation) {
	//			if (exist->cache.start != rec->cache.start ||
	//			    exist->cache.size != rec->cache.size ||
	//			    memcmp(exist->csum, rec->csum, BTRFS_CSUM_SIZE)) {
	//				ret = -EEXIST;
	//			} else {
	//				BUG_ON(exist->nmirrors >= BTRFS_MAX_MIRRORS);
	//				exist->devices[exist->nmirrors] = device;
	//				exist->offsets[exist->nmirrors] = offset;
	//				exist->nmirrors++;
	//			}
	//			goto free_out;
	//		}
	//		remove_cache_extent(eb_cache, cache);
	//		free(exist);
	//		goto again;
	//	}
	//
	//	rec->devices[0] = device;
	//	rec->offsets[0] = offset;
	//	rec->nmirrors++;
	//	ret = insert_cache_extent(eb_cache, &rec->cache);
	//	BUG_ON(ret);
	//out:
	//	return ret;
	//free_out:
	//	free(rec);
	//	goto out;
	return false
}
func extract_metadata_record(rc *Recover_control,
	leaf *Extent_buffer) bool {
	var (
		//	struct btrfs_key key;
		ret = false
	//	int i;
	//	u32 nritems;
	//
	)
	nritems := (*Btrfs_header)(unsafe.Pointer(&leaf.Data)).Nritems
	for i := uint64(0); i < nritems; i++ {
		//		btrfs_item_key_to_cpu(leaf, &key, i);
		//		switch (key.type) {
		//		case BTRFS_BLOCK_GROUP_ITEM_KEY:
		//			pthread_mutex_lock(&rc->rc_lock);
		//			ret = process_block_group_item(&rc->bg, leaf, &key, i);
		//			pthread_mutex_unlock(&rc->rc_lock);
		//			break;
		//		case BTRFS_CHUNK_ITEM_KEY:
		//			pthread_mutex_lock(&rc->rc_lock);
		//			ret = process_chunk_item(&rc->chunk, leaf, &key, i);
		//			pthread_mutex_unlock(&rc->rc_lock);
		//			break;
		//		case BTRFS_DEV_EXTENT_KEY:
		//			pthread_mutex_lock(&rc->rc_lock);
		//			ret = process_device_extent_item(&rc->devext, leaf,
		//							 &key, i);
		//			pthread_mutex_unlock(&rc->rc_lock);
		//			break;
		//		}
		//		if (ret)
		//			break;
	}
	return ret
}
