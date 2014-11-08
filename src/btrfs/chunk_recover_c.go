package btrfs

import (
	"container/list"
	"fmt"
	"log"
	"os"
	"syscall"
	"unsafe"

	"github.com/monnand/GoLLRB/llrb"
)

func IsSuperBlockAddress(offset uint64) bool {
	var i int

	for i = 0; i < BTRFS_SUPER_MIRROR_MAX; i++ {
		if offset == BtrfsSbOffset(i) {
			return true
		}
	}
	return false
}
func NewRecoverControl(verbose bool, yes bool) *RecoverControl {

	rc := &RecoverControl{
		Chunk:   llrb.New(),
		EbCache: llrb.New(),
		Bg: BlockGroupTree{Tree: llrb.New(),
			Block_Groups: list.New(),
		},
		Devext: DeviceExtentTree{Tree: llrb.New(),
			ChunkOrphans:  list.New(),
			DeviceOrphans: list.New(),
		},
		Verbose:          verbose,
		Yes:              yes,
		GoodChunks:       list.New(),
		BadChunks:        list.New(),
		UnrepairedChunks: list.New(),
		FsDevices: &BtrfsFsDevices{
			Devices: list.New(),
		},
	}
	rc.FsDevices.Devices.PushBack(&BtrfsDevice{
		Devid: 1,
	},
	)
	return rc
	//	pthreadMutexInit(&rc->rcLock, NULL);
}

func RecoverPrepare(rc *RecoverControl, path string) bool {
	//	int ret;
	//	int fd;
	//	struct btrfsSuperBlock *sb;
	//	struct btrfsFsDevices *fsDevices;
	//
	//	ret = 0;
	//	fd = open(path, O_RDONLY);
	//	if (fd < 0) {
	//		fprintf(stderr, "open %s\n error.\n", path);
	//		return -1;
	//	}
	//
	//	sb = malloc(sizeof(struct btrfsSuperBlock));
	//	if (!sb) {
	//		fprintf(stderr, "allocating memory for sb failed.\n");
	//		ret = -ENOMEM;
	//		goto failCloseFd;
	//	}
	//
	//	ret = btrfsReadDevSuper(fd, sb, BTRFS_SUPER_INFO_OFFSET, 1);
	//	if (ret) {
	//		fprintf(stderr, "read super block error\n");
	//		goto failFreeSb;
	//	}
	//
	//	rc->sectorsize = btrfsSuperSectorsize(sb);
	//	rc->leafsize = btrfsSuperLeafsize(sb);
	//	rc->generation = btrfsSuperGeneration(sb);
	//	rc->chunkRootGeneration = btrfsSuperChunkRootGeneration(sb);
	//	rc->csumSize = btrfsSuperCsumSize(sb);
	//
	//	/* if seed, the result of scanning below will be partial */
	//	if (btrfsSuperFlags(sb) & BTRFS_SUPER_FLAG_SEEDING) {
	//		fprintf(stderr, "this device is seed device\n");
	//		ret = -1;
	//		goto failFreeSb;
	//	}
	//
	//	ret = btrfsScanFsDevices(fd, path, &fsDevices, 0, 1, 1);
	//	if (ret)
	//		goto failFreeSb;
	//
	//	rc->fsDevices = fsDevices;
	//
	//	if (rc->verbose)
	//		printAllDevices(&rc->fsDevices->devices);
	//
	//failFreeSb:
	//	free(sb);
	//failCloseFd:
	//	close(fd);
	//	return ret;
	var sb BtrfsSuperBlock
	fd, err := syscall.Open(path, 0, 0)
	if err != nil {
		log.Fatal(os.NewSyscallError("open", err))
	}
	rc.Fd = fd
	ret := btrfsReadDevSuper(fd, &sb, BTRFS_SUPER_INFO_OFFSET, true)
	if !ret {
		fmt.Errorf("read super block error\n")
	} else {
		fmt.Printf("\nSB: %+v\n", sb)
		rc.Sectorsize = sb.Sectorsize
		rc.Leafsize = sb.Leafsize
		rc.Generation = sb.Generation
		rc.ChunkRootGeneration = sb.ChunkRootGeneration
		rc.CsumSize = BtrfsSuperCsumSize(&sb)
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
func btrfsRecoverChunkTree(path []byte, verbose bool, yes bool) int {
	//	int ret = 0;
	//	struct btrfsRoot *root = NULL;
	//	struct btrfsTransHandle *trans;
	//
	//	rc := NewRecoverControl( verbose, yes)
	//
	//	ret = recoverPrepare(&rc, path);
	//	if (ret) {
	//		fprintf(stderr, "recover prepare error\n");
	//		return ret;
	//	}
	//
	//	ret = scanDevices(&rc);
	//	if (ret) {
	//		fprintf(stderr, "scan chunk headers error\n");
	//		goto failRc;
	//	}
	//
	//	if (cacheTreeEmpty(&rc.chunk) &&
	//	    cacheTreeEmpty(&rc.bg.tree) &&
	//	    cacheTreeEmpty(&rc.devext.tree)) {
	//		fprintf(stderr, "no recoverable chunk\n");
	//		goto failRc;
	//	}
	//
	//	printScanResult(&rc);
	//
	//	ret = checkChunks(&rc.chunk, &rc.bg, &rc.devext, &rc.goodChunks,
	//			   &rc.badChunks, 1);
	//	printCheckResult(&rc);
	//	if (ret) {
	//		if (!listEmpty(&rc.bg.blockGroups) ||
	//		    !listEmpty(&rc.devext.noChunkOrphans)) {
	//			ret = btrfsRecoverChunks(&rc);
	//			if (ret)
	//				goto failRc;
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
	//		goto failRc;
	//	}
	//
	//	root = openCtreeWithBrokenChunk(&rc);
	//	if (IS_ERR(root)) {
	//		fprintf(stderr, "open with broken chunk error\n");
	//		ret = PTR_ERR(root);
	//		goto failRc;
	//	}
	//
	//	ret = checkAllChunksByMetadata(&rc, root);
	//	if (ret) {
	//		fprintf(stderr, "The chunks in memory can not match the metadata of the fs. Repair failed.\n");
	//		goto failCloseCtree;
	//	}
	//
	//	ret = btrfsRebuildOrderedDataChunkStripes(&rc, root);
	//	if (ret) {
	//		fprintf(stderr, "Failed to rebuild ordered chunk stripes.\n");
	//		goto failCloseCtree;
	//	}
	//
	//	if (!rc.yes) {
	//		ret = askUser("We are going to rebuild the chunk tree on disk, it might destroy the old metadata on the disk, Are you sure?");
	//		if (!ret) {
	//			ret = 1;
	//			goto failCloseCtree;
	//		}
	//	}
	//
	//	trans = btrfsStartTransaction(root, 1);
	//	ret = removeChunkExtentItem(trans, &rc, root);
	//	BUG_ON(ret);
	//
	//	ret = rebuildChunkTree(trans, &rc, root);
	//	BUG_ON(ret);
	//
	//	ret = rebuildSysArray(&rc, root);
	//	BUG_ON(ret);
	//
	//	btrfsCommitTransaction(trans, root);
	//failCloseCtree:
	//	closeCtree(root);
	//failRc:
	//	freeRecoverControl(&rc);
	//	return ret;
	return 0
}
func scanOneDevice(devScan *DeviceScan) bool {
	var (
		buf    *ExtentBuffer
		bytenr uint64
		ret    bool = false
		//struct deviceScan *devScan = (struct deviceScan *)devScanStruct;
		rc *RecoverControl = devScan.Rc
		//		device *BtrfsDevice    = devScan.Dev
		fd int = devScan.Fd
		//		oldtype int
		h *BtrfsHeader
	)
	buf = new(ExtentBuffer)
	buf.Len = uint64(rc.Leafsize)
	buf.Data = make([]byte, rc.Leafsize)
loop:
	for bytenr = 0; ; bytenr += uint64(rc.Sectorsize) {
		ret = false
		if IsSuperBlockAddress(bytenr) {
			bytenr += uint64(rc.Sectorsize)
		}
		n, err := syscall.Pread(fd, buf.Data, int64(bytenr))
		if err != nil {
			log.Fatalln(os.NewSyscallError("pread64", err))
		}
		if n < int(rc.Leafsize) {
			break loop
		}
		h = (*BtrfsHeader)(unsafe.Pointer(&buf.Data))
		if rc.FsDevices.Fsid != h.Fsid || verifyTreeBlockCsumSilent(buf, uint16(rc.CsumSize)) {
			continue loop
		}
		rc.RcLock.Lock()
		//		ret = processExtentBuffer(&rc.EbCache, buf, device, bytenr)
		rc.RcLock.Unlock()
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
				if ret = extractMetadataRecord(rc, buf); !ret {
					break loop
				}
			case BTRFS_CHUNK_TREE_OBJECTID:
				if h.Generation > rc.ChunkRootGeneration {
					continue loop
				}
				if ret = extractMetadataRecord(rc, buf); !ret {
					break loop
				}
			}
		}
	}

	//	close(fd);
	//	free(buf);
	return ret
}

//func processExtentBuffer(ebCache *CacheTree,
//	eb *ExtentBuffer,
//	device *BtrfsDevice, offset uint64) bool {
//	struct extentRecord *rec;
//	struct extentRecord *exist;
//	struct cacheExtent *cache;
//	int ret = 0;
//
//	rec = btrfsNewExtentRecord(eb);
//	if (!rec->cache.size)
//		goto freeOut;
//again:
//	cache = lookupCacheExtent(ebCache,
//				    rec->cache.start,
//				    rec->cache.size);
//	if (cache) {
//		exist = containerOf(cache, struct extentRecord, cache);
//
//		if (exist->generation > rec->generation)
//			goto freeOut;
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
//			goto freeOut;
//		}
//		removeCacheExtent(ebCache, cache);
//		free(exist);
//		goto again;
//	}
//
//	rec->devices[0] = device;
//	rec->offsets[0] = offset;
//	rec->nmirrors++;
//	ret = insertCacheExtent(ebCache, &rec->cache);
//	BUG_ON(ret);
//out:
//	return ret;
//freeOut:
//	free(rec);
//	goto out;
//	return false
//}
func extractMetadataRecord(rc *RecoverControl, leaf *ExtentBuffer) bool {
	var (
		//	struct btrfsKey key;
		ret = false
	//	int i;
	//	u32 nritems;
	//
	)
	nritems := uint64((*BtrfsHeader)(unsafe.Pointer(&leaf.Data)).Nritems)
	for i := uint64(0); i < nritems; i++ {
		//		btrfsItemKeyToCpu(leaf, &key, i);
		//		switch (key.type) {
		//		case BTRFS_BLOCK_GROUP_ITEM_KEY:
		//			pthreadMutexLock(&rc->rcLock);
		//			ret = processBlockGroupItem(&rc->bg, leaf, &key, i);
		//			pthreadMutexUnlock(&rc->rcLock);
		//			break;
		//		case BTRFS_CHUNK_ITEM_KEY:
		//			pthreadMutexLock(&rc->rcLock);
		//			ret = processChunkItem(&rc->chunk, leaf, &key, i);
		//			pthreadMutexUnlock(&rc->rcLock);
		//			break;
		//		case BTRFS_DEV_EXTENT_KEY:
		//			pthreadMutexLock(&rc->rcLock);
		//			ret = processDeviceExtentItem(&rc->devext, leaf,
		//							 &key, i);
		//			pthreadMutexUnlock(&rc->rcLock);
		//			break;
		//		}
		//		if (ret)
		//			break;
	}
	return ret
}
