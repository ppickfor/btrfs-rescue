package btrfs

import (
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"log"
	"os"
	"syscall"
	"unsafe"
)

func Is_super_block_address(offset uint64) bool {
	var i int

	for i = 0; i < BTRFS_SUPER_MIRROR_MAX; i++ {
		if offset == Btrfs_sb_offset(i) {
			return true
		}
	}
	return false
}
func Init_recover_control(rc *Recover_control, verbose bool,
	yes bool) {
	//	memset(rc, 0, sizeof(struct recover_control));
	//	cache_tree_init(&rc->chunk);
	//	cache_tree_init(&rc->eb_cache);
	//	block_group_tree_init(&rc->bg);
	//	device_extent_tree_init(&rc->devext);
	//
	//	INIT_LIST_HEAD(&rc->good_chunks);
	//	INIT_LIST_HEAD(&rc->bad_chunks);
	//	INIT_LIST_HEAD(&rc->unrepaired_chunks);
	//
	rc.Chunk = llrb.New()
	rc.Eb_cache = llrb.New()
	rc.Bg.Tree = llrb.New()
	rc.Bg.Block_Groups = List_head{Next: &rc.Bg.Block_Groups, Prev: &rc.Bg.Block_Groups}
	rc.Devext.Tree = llrb.New()
	rc.Devext.Chunk_orphans = List_head{Next: &rc.Devext.Chunk_orphans, Prev: &rc.Devext.Chunk_orphans}
	rc.Devext.Device_orphans = List_head{Next: &rc.Devext.Device_orphans, Prev: &rc.Devext.Device_orphans}

	rc.Verbose = verbose
	rc.Yes = yes
	//	pthread_mutex_init(&rc->rc_lock, NULL);
}

func Recover_prepare(rc *Recover_control, path string) bool {
	//	int ret;
	//	int fd;
	//	struct btrfs_super_block *sb;
	//	struct btrfs_fs_devices *fs_devices;
	//
	//	ret = 0;
	//	fd = open(path, O_RDONLY);
	//	if (fd < 0) {
	//		fprintf(stderr, "open %s\n error.\n", path);
	//		return -1;
	//	}
	//
	//	sb = malloc(sizeof(struct btrfs_super_block));
	//	if (!sb) {
	//		fprintf(stderr, "allocating memory for sb failed.\n");
	//		ret = -ENOMEM;
	//		goto fail_close_fd;
	//	}
	//
	//	ret = btrfs_read_dev_super(fd, sb, BTRFS_SUPER_INFO_OFFSET, 1);
	//	if (ret) {
	//		fprintf(stderr, "read super block error\n");
	//		goto fail_free_sb;
	//	}
	//
	//	rc->sectorsize = btrfs_super_sectorsize(sb);
	//	rc->leafsize = btrfs_super_leafsize(sb);
	//	rc->generation = btrfs_super_generation(sb);
	//	rc->chunk_root_generation = btrfs_super_chunk_root_generation(sb);
	//	rc->csum_size = btrfs_super_csum_size(sb);
	//
	//	/* if seed, the result of scanning below will be partial */
	//	if (btrfs_super_flags(sb) & BTRFS_SUPER_FLAG_SEEDING) {
	//		fprintf(stderr, "this device is seed device\n");
	//		ret = -1;
	//		goto fail_free_sb;
	//	}
	//
	//	ret = btrfs_scan_fs_devices(fd, path, &fs_devices, 0, 1, 1);
	//	if (ret)
	//		goto fail_free_sb;
	//
	//	rc->fs_devices = fs_devices;
	//
	//	if (rc->verbose)
	//		print_all_devices(&rc->fs_devices->devices);
	//
	//fail_free_sb:
	//	free(sb);
	//fail_close_fd:
	//	close(fd);
	//	return ret;
	var sb Btrfs_super_block
	fd, err := syscall.Open(path, 0, 0)
	if err != nil {
		log.Fatal(os.NewSyscallError("open", err))
	}
	rc.Fd = fd
	ret := btrfs_read_dev_super(fd, &sb, BTRFS_SUPER_INFO_OFFSET, true)
	if !ret {
		fmt.Errorf("read super block error\n")
	} else {
		fmt.Printf("\nSB: %+v\n", sb)
		rc.Sectorsize = sb.Sectorsize
		rc.Leafsize = sb.Leafsize
		rc.Generation = sb.Generation
		rc.Chunk_root_generation = sb.Chunk_root_generation
		rc.Csum_size = Btrfs_super_csum_size(&sb)
		rc.Fsid = sb.Fsid
		fmt.Printf("\nRC: %+v\n", rc)
		//		var buf []byte = make([]byte, rc.Leafsize)
		/* if seed, the result of scanning below will be partial */
		if (sb.Flags & BTRFS_SUPER_FLAG_SEEDING) != 0 {
			fmt.Errorf("this device is seed device\n")
			ret = false
		}
	}
	return ret
}

/*
 * Return 0 when succesful, < 0 on error and > 0 if aborted by user
 */
func btrfs_recover_chunk_tree(path []byte, verbose bool, yes bool) int {
	//	int ret = 0;
	//	struct btrfs_root *root = NULL;
	//	struct btrfs_trans_handle *trans;
	var (
		rc = new(Recover_control)
	)
	//
	Init_recover_control(rc, verbose, yes)
	//
	//	ret = recover_prepare(&rc, path);
	//	if (ret) {
	//		fprintf(stderr, "recover prepare error\n");
	//		return ret;
	//	}
	//
	//	ret = scan_devices(&rc);
	//	if (ret) {
	//		fprintf(stderr, "scan chunk headers error\n");
	//		goto fail_rc;
	//	}
	//
	//	if (cache_tree_empty(&rc.chunk) &&
	//	    cache_tree_empty(&rc.bg.tree) &&
	//	    cache_tree_empty(&rc.devext.tree)) {
	//		fprintf(stderr, "no recoverable chunk\n");
	//		goto fail_rc;
	//	}
	//
	//	print_scan_result(&rc);
	//
	//	ret = check_chunks(&rc.chunk, &rc.bg, &rc.devext, &rc.good_chunks,
	//			   &rc.bad_chunks, 1);
	//	print_check_result(&rc);
	//	if (ret) {
	//		if (!list_empty(&rc.bg.block_groups) ||
	//		    !list_empty(&rc.devext.no_chunk_orphans)) {
	//			ret = btrfs_recover_chunks(&rc);
	//			if (ret)
	//				goto fail_rc;
	//		}
	//		/*
	//		 * If the chunk is healthy, its block group item and device
	//		 * extent item should be written on the disks. So, it is very
	//		 * likely that the bad chunk is a old one that has been
	//		 * droppped from the fs. Don't deal with them now, we will
	//		 * check it after the fs is opened.
	//		 */
	//	} else {
	//		fprintf(stderr, "Check chunks successfully with no orphans\n");
	//		goto fail_rc;
	//	}
	//
	//	root = open_ctree_with_broken_chunk(&rc);
	//	if (IS_ERR(root)) {
	//		fprintf(stderr, "open with broken chunk error\n");
	//		ret = PTR_ERR(root);
	//		goto fail_rc;
	//	}
	//
	//	ret = check_all_chunks_by_metadata(&rc, root);
	//	if (ret) {
	//		fprintf(stderr, "The chunks in memory can not match the metadata of the fs. Repair failed.\n");
	//		goto fail_close_ctree;
	//	}
	//
	//	ret = btrfs_rebuild_ordered_data_chunk_stripes(&rc, root);
	//	if (ret) {
	//		fprintf(stderr, "Failed to rebuild ordered chunk stripes.\n");
	//		goto fail_close_ctree;
	//	}
	//
	//	if (!rc.yes) {
	//		ret = ask_user("We are going to rebuild the chunk tree on disk, it might destroy the old metadata on the disk, Are you sure?");
	//		if (!ret) {
	//			ret = 1;
	//			goto fail_close_ctree;
	//		}
	//	}
	//
	//	trans = btrfs_start_transaction(root, 1);
	//	ret = remove_chunk_extent_item(trans, &rc, root);
	//	BUG_ON(ret);
	//
	//	ret = rebuild_chunk_tree(trans, &rc, root);
	//	BUG_ON(ret);
	//
	//	ret = rebuild_sys_array(&rc, root);
	//	BUG_ON(ret);
	//
	//	btrfs_commit_transaction(trans, root);
	//fail_close_ctree:
	//	close_ctree(root);
	//fail_rc:
	//	free_recover_control(&rc);
	//	return ret;
	return 0
}
func scan_one_device(dev_scan *Device_scan) bool {
	var (
		buf    *Extent_buffer
		bytenr uint64
		ret    bool = false
		//struct device_scan *dev_scan = (struct device_scan *)dev_scan_struct;
		rc *Recover_control = dev_scan.Rc
		//		device *Btrfs_device    = dev_scan.Dev
		fd int = dev_scan.Fd
		//		oldtype int
		h *Btrfs_header
	)
	buf = new(Extent_buffer)
	buf.Len = uint64(rc.Leafsize)
	buf.Data = make([]byte, rc.Leafsize)
loop:
	for bytenr = 0; ; bytenr += uint64(rc.Sectorsize) {
		ret = false
		if Is_super_block_address(bytenr) {
			bytenr += uint64(rc.Sectorsize)
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
		//		ret = process_extent_buffer(&rc.Eb_cache, buf, device, bytenr)
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

//func process_extent_buffer(eb_cache *Cache_tree,
//	eb *Extent_buffer,
//	device *Btrfs_device, offset uint64) bool {
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
//	return false
//}
func extract_metadata_record(rc *Recover_control,
	leaf *Extent_buffer) bool {
	var (
		//	struct btrfs_key key;
		ret = false
	//	int i;
	//	u32 nritems;
	//
	)
	nritems := uint64((*Btrfs_header)(unsafe.Pointer(&leaf.Data)).Nritems)
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
