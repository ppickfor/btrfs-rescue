package btrfs

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"syscall"

	"github.com/monnand/GoLLRB/llrb"
)

// extractMetadataRecord iterates items of the current block and proccess each blockgroup, chunk and dev extent adding them to caches
func ExtractMetadataRecord(rc *RecoverControl, generation uint64, items []BtrfsItem, itemBuf []byte) {

	for _, item := range items {
		switch item.Key.Type {
		case BTRFS_BLOCK_GROUP_ITEM_KEY:
			//			pthreadMutexLock(&rc->rcLock);
			//			fmt.Fprintf(os.Stderr,"BLOCK Group Item: %+v\n", &leaf.Items[i])
			processBlockGroupItem(&rc.Bg, generation, &item, itemBuf)
			//			pthreadMutexUnlock(&rc->rcLock);
			//			break;
		case BTRFS_CHUNK_ITEM_KEY:
			//			pthreadMutexLock(&rc->rcLock);
			processChunkItem(rc.Chunk, generation, &item, itemBuf)
			//			pthreadMutexUnlock(&rc->rcLock);
			//			break;
		case BTRFS_DEV_EXTENT_KEY:
			//			pthreadMutexLock(&rc->rcLock);
			processDeviceExtentItem(&rc.Devext, generation, &item, itemBuf)
			//			pthreadMutexUnlock(&rc->rcLock);
			//			break;
		}
	}
}

type InodeRec struct {
	Parent              uint64                             // from diritem or inoderef
	Name                string                             // from diritem or inoderef
	Type                uint8                              // from diritem
	InodeItem           BtrfsInodeItem                     // from inodeitem
	FileExtentItems     map[uint64]BtrfsFileExtentItem     // from fileextentsitem
	FileExtentItemsCont map[uint64]BtrfsFileExtentItemCont // from fileextentsitem
	DirItems            map[uint64]bool
	Data                []byte
}
type InodeKey struct {
	Owner uint64
	Inode uint64
}
type RootData struct {
	Name     string
	Treeid   uint64
	Dirid    uint64
	Sequence uint64
}

var Inodes = make(map[InodeKey]InodeRec)
var Roots = make(map[uint64]RootData)

func ProcessCsumItem(rootRefItemCache *llrb.LLRB, owner uint64, item *BtrfsItem, itemBuf []byte) {
	itemPtr := itemBuf[item.Offset : item.Offset+item.Size]
	nCsums := item.Size / 4
	csums := make([]uint32, nCsums)
	bytereader := bytes.NewReader(itemPtr)
	//	csumItem := new(BtrfsCsumItem)
	_ = binary.Read(bytereader, binary.LittleEndian, nCsums)
	fmt.Fprintf(os.Stderr, "csumItem owner: %d key: %+v\n", owner, item)
	fmt.Fprintf(os.Stderr, "csumItem value %+v\n", csums)
}

// ProcessRootRef extract root ref or backref info from item and bytebuffer
func ProcessRootRef(rootRefItemCache *llrb.LLRB, owner uint64, item *BtrfsItem, itemBuf []byte) {

	itemPtr := itemBuf[item.Offset : item.Offset+item.Size]
	bytereader := bytes.NewReader(itemPtr)
	rootRefItem := new(BtrfsRootRef)
	_ = binary.Read(bytereader, binary.LittleEndian, rootRefItem)
	namebytes := make([]byte, rootRefItem.Len)
	_ = binary.Read(bytereader, binary.LittleEndian, namebytes)
	switch item.Key.Type {
	case BTRFS_ROOT_REF_KEY:
		Roots[item.Key.Offset] = RootData{
			Name:     string(namebytes),
			Dirid:    rootRefItem.Dirid,
			Treeid:   item.Key.Objectid,
			Sequence: rootRefItem.Sequence,
		}
	case BTRFS_ROOT_BACKREF_KEY:
		Roots[item.Key.Objectid] = RootData{
			Name:     string(namebytes),
			Dirid:    rootRefItem.Dirid,
			Treeid:   item.Key.Offset,
			Sequence: rootRefItem.Sequence,
		}
	}
}

// ProcessFileExtentItem extract file extent data from itembuf using the current item
func ProcessFileExtentItem(fileExtentItemCache *llrb.LLRB, owner uint64, item *BtrfsItem, itemBuf []byte) {

	//	BtrfsPrintKey(&item.Key)
	//	key := item.Key
	itemPtr := itemBuf[item.Offset : item.Offset+item.Size]
	bytereader := bytes.NewReader(itemPtr)
	fileExtentItem := new(BtrfsFileExtentItem)
	_ = binary.Read(bytereader, binary.LittleEndian, fileExtentItem)
	k := InodeKey{
		Owner: owner,
		Inode: item.Key.Objectid,
	}
	tmp := Inodes[k]
	if tmp.FileExtentItems == nil {
		tmp.FileExtentItems = make(map[uint64]BtrfsFileExtentItem)
	}
	if fileExtentItem.Generation < tmp.FileExtentItems[item.Key.Offset].Generation {
		return
	}
	tmp.FileExtentItems[item.Key.Offset] = *fileExtentItem
	switch fileExtentItem.Type {
	case BTRFS_FILE_EXTENT_INLINE:
		// data follows
		//		fmt.Fprintf(os.Stderr," Inline Size:%d %d", binary.Size(fileExtentItem), item.Size)

		//		fmt.Fprintf(os.Stderr,"Data: %s\n", string(itemBuf[item.Offset+21:]))
		tmp.Data = itemBuf[item.Offset+21:]
	case BTRFS_FILE_EXTENT_REG, BTRFS_FILE_EXTENT_PREALLOC:
		fileExtentItemCont := new(BtrfsFileExtentItemCont)
		_ = binary.Read(bytereader, binary.LittleEndian, fileExtentItemCont)
		if tmp.FileExtentItemsCont == nil {
			tmp.FileExtentItemsCont = make(map[uint64]BtrfsFileExtentItemCont)
		}
		tmp.FileExtentItems[item.Key.Offset] = *fileExtentItem
		tmp.FileExtentItemsCont[item.Key.Offset] = *fileExtentItemCont
		//		fmt.Fprintf(os.Stderr,"ProcessFileExtentItem ite %+v\n", item)
		//		fmt.Fprintf(os.Stderr,"BtrfsFileExtentItem %+v\n", fileExtentItem)
		//		fmt.Fprintf(os.Stderr,"BtrfsFileExtentItemCont %+v\n", fileExtentItemCont)
	}
	Inodes[k] = tmp
}

// processInodeItem extract inode ref from itembuf using the current item
func ProcessInodeItem(inodeItemCache *llrb.LLRB, owner uint64, item *BtrfsItem, itemBuf []byte) {

	itemPtr := itemBuf[item.Offset:]
	bytereader := bytes.NewReader(itemPtr)
	inodeItem := new(BtrfsInodeItem)
	_ = binary.Read(bytereader, binary.LittleEndian, inodeItem)
	k := InodeKey{
		Owner: owner,
		Inode: item.Key.Objectid,
	}
	tmp := Inodes[k]
	tmp.InodeItem = *inodeItem
	Inodes[k] = tmp
}

// processInodeRefItem extract inode ref from itembuf using the current item
func ProcessInodeRefItem(inodeRefItemCache *llrb.LLRB, owner uint64, item *BtrfsItem, itemBuf []byte) {
	itemPtr := itemBuf[item.Offset : item.Offset+item.Size]
	bytereader := bytes.NewReader(itemPtr)
	for bytereader.Len() != 0 {
		inodeRef := new(BtrfsInodeRef)
		_ = binary.Read(bytereader, binary.LittleEndian, inodeRef)
		namebytes := make([]byte, inodeRef.Len)
		_ = binary.Read(bytereader, binary.LittleEndian, namebytes)
		k := InodeKey{
			Owner: owner,
			Inode: item.Key.Objectid,
		}
		tmp := Inodes[k]
		if tmp.Name == "" {
			tmp.Name = string(namebytes)
		}
		if tmp.Parent == 0 {
			tmp.Parent = item.Key.Offset
		}
		Inodes[k] = tmp
		// dont point back to yourself
		if item.Key.Objectid != item.Key.Offset {
			// build dir items
			k = InodeKey{
				Owner: owner,
				Inode: item.Key.Offset,
			}
			tmp = Inodes[k]
			if tmp.DirItems == nil {
				tmp.DirItems = make(map[uint64]bool)
			}
			tmp.Type = BTRFS_FT_DIR
			tmp.DirItems[item.Key.Objectid] = true
			Inodes[k] = tmp
		} else {
			// remove name for root dir
			k = InodeKey{
				Owner: owner,
				Inode: item.Key.Offset,
			}
			tmp = Inodes[k]
			tmp.Name = ""
			Inodes[k] = tmp
		}
	}
}

// processDirItem extract dir item from itembuf using the current item
func ProcessDirItem(dirItemCache *llrb.LLRB, owner uint64, item *BtrfsItem, itemBuf []byte) {

	//	key := item.Key
	//		BtrfsPrintKey(&item.Key)
	itemPtr := itemBuf[item.Offset : item.Offset+item.Size]
	bytereader := bytes.NewReader(itemPtr)
	//		fmt.Fprintf(os.Stderr," ")
	for bytereader.Len() != 0 {
		dirItem := new(BtrfsDirItem)
		_ = binary.Read(bytereader, binary.LittleEndian, dirItem)
		//		printDirItemType(dirItem)
		//			BtrfsPrintKey(&dirItem.Location)
		length := dirItem.NameLen
		if length > BTRFS_NAME_LEN {
			length = BTRFS_NAME_LEN
		}
		//		fmt.Fprintf(os.Stderr," namelen: %d ",length)
		namebytes := make([]byte, length)
		_ = binary.Read(bytereader, binary.LittleEndian, namebytes)

		//				fmt.Fprintf(os.Stderr," %s ", string(namebytes))
		if dirItem.DataLen != 0 {
			databytes := make([]byte, dirItem.DataLen)
			_ = binary.Read(bytereader, binary.LittleEndian, databytes)
			//			fmt.Fprintf(os.Stderr,"datalen: %d", dirItem.DataLen)
		}
		k := InodeKey{
			Owner: owner,
			Inode: dirItem.Location.Objectid,
		}
		tmp := Inodes[k]
		tmp.Name = string(namebytes)
		tmp.Parent = item.Key.Objectid
		tmp.Type = dirItem.Type
		Inodes[k] = tmp
		if dirItem.Location.Objectid != item.Key.Objectid {
			// build dir items
			k = InodeKey{
				Owner: owner,
				Inode: item.Key.Objectid,
			}
			tmp = Inodes[k]
			if tmp.DirItems == nil {
				tmp.DirItems = make(map[uint64]bool)
			}
			tmp.Type = BTRFS_FT_DIR
			tmp.DirItems[dirItem.Location.Objectid] = true
			Inodes[k] = tmp
		} else {
			// remove name for root dir
			k = InodeKey{
				Owner: owner,
				Inode: item.Key.Offset,
			}
			tmp = Inodes[k]
			tmp.Name = ""
			Inodes[k] = tmp
		}
	}
	//		fmt.Fprintf(os.Stderr,"\n")
}

// processBlockGroupItem creates a new block group record and update the cache by latest generation for each blockgroup item
func processBlockGroupItem(bgCache *BlockGroupTree, generation uint64, item *BtrfsItem, itemBuf []byte) {

	rec := NewBlockGroupRecord(generation, item, itemBuf)
again:
	if i := bgCache.Tree.Get(rec); i != nil {
		exists := i.(*BlockGroupRecord)
		/*check the generation and replace if needed*/
		if exists.Generation > rec.Generation {
			return
		}
		if exists.Generation == rec.Generation {
			// TODO
			//			int offset = offsetof(struct blockGroupRecord,
			//					      generation);
			/*
			 * According to the current kernel code, the following
			 * case is impossble, or there is something wrong in
			 * the kernel code.
			 */
			//			if (memcmp(((void *)exist) + offset,
			//				   ((void *)rec) + offset,
			//				   sizeof(*rec) - offset))
			//				ret = -EEXIST;
			//			goto freeOut;
			fmt.Fprintf(os.Stderr, "processBlockGroupItem: same generation %+v\n", rec)
			return
		}
		bgCache.Block_Groups.Remove(exists.List)
		exists.List = nil
		//		listDelInit(&exists.List)
		bgCache.Tree.Delete(exists)
		/*
		 * We must do seach again to avoid the following cache.
		 * /--old bg 1--//--old bg 2--/
		 *        /--new bg--/
		 */
		goto again
	}
	bgCache.Tree.InsertNoReplace(rec)
	rec.List = bgCache.Block_Groups.PushBack(rec)
}

// MapLogical maps a logical to a physical address
func MapLogical(mapCache *llrb.LLRB, logical uint64) (error, uint64) {
	mapSearch := &MapLookup{
		CacheExtent: CacheExtent{
			Start: logical,
		},
	}
	var mapResult *MapLookup
	// find previous MapLookup
	mapCache.DescendLessOrEqual(mapSearch, func(i llrb.Item) bool {
		mapResult = i.(*MapLookup)
		return false
	})
	if mapResult != nil {
		fmt.Fprintf(os.Stderr, "MapLogical: searched %d, found %v\n", logical, mapResult)
		return nil, mapResult.Stripes[mapResult.NumStripes-1].Physical + logical - mapResult.CacheExtent.Start
	}
	return errors.New("No mapping"), 0
}

// processDeviceExtentItem creates a new device extent record and update the cache by latest generation
func processDeviceExtentItem(devextCache *DeviceExtentTree, generation uint64, item *BtrfsItem, itemBuf []byte) {

	rec := NewDeviceExtentRecord(generation, item, itemBuf)
again:
	if i := devextCache.Tree.Get(rec); i != nil {
		exists := i.(*DeviceExtentRecord)
		// check the generation and replace if needed
		if exists.Generation > rec.Generation {
			return
		}
		if exists.Generation == rec.Generation {
			// should not happen
			fmt.Fprintf(os.Stderr, "processDeviceExtentItem: same generation %+v\n", rec)
			return
		}
		devextCache.ChunkOrphans.Remove(exists.ChunkList)
		devextCache.DeviceOrphans.Remove(exists.DeviceList)
		exists.ChunkList = nil
		exists.DeviceList = nil
		//		listDelInit(&exists.ChunkList)
		//		listDelInit(&exists.DeviceList)
		devextCache.Tree.Delete(exists)
		/*
		 * We must do seach again to avoid the following cache.
		 * /--old bg 1--//--old bg 2--/
		 *        /--new bg--/
		 */
		goto again
	}
	rec.ChunkList = devextCache.ChunkOrphans.PushBack(rec)
	rec.DeviceList = devextCache.DeviceOrphans.PushBack(rec)
	devextCache.Tree.InsertNoReplace(rec)
}

// processChunkItem: inserts new chunks into the chunk cache based on generation
func processChunkItem(chunkCache *llrb.LLRB, generation uint64, item *BtrfsItem, itemBuf []byte) {
	//	fmt.Fprintf(os.Stderr,"found chunk item\n")
	rec := NewChunkRecord(generation, item, itemBuf)
again:
	if i := chunkCache.Get(rec); i != nil {
		exists := i.(*ChunkRecord)
		if exists.Generation > rec.Generation {
			return
		}
		if exists.Generation == rec.Generation {
			numStripes := rec.NumStripes
			if exists.NumStripes != numStripes {
				fmt.Fprintf(os.Stderr, "processChunkItem: same generation but different details %+v\n", rec)
				// TODO
				// compare everything from generation to end of the record (varies by num stripes)
				//int num_stripes = rec->num_stripes;
				//			int rec_size = btrfs_chunk_record_size(num_stripes);
				//			int offset = offsetof(struct chunk_record, generation);
				//
				//			if (exist->num_stripes != rec->num_stripes ||
				//			    memcmp(((void *)exist) + offset,
				//				   ((void *)rec) + offset,
				//				   rec_size - offset))
				//				ret = -EEXIST;
				//			goto free_out;
				return
			}
		}
		chunkCache.Delete(exists)
		goto again
	}
	chunkCache.InsertNoReplace(rec)
}

func IsSuperBlockAddress(offset uint64) bool {
	var i int

	for i = 0; i < BTRFS_SUPER_MIRROR_MAX; i++ {
		if offset == BtrfsSbOffset(i) {
			return true
		}
	}
	return false
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
		fmt.Fprintf(os.Stderr, "\nSB: %+v\n", sb)
		rc.Sectorsize = sb.Sectorsize
		rc.Leafsize = sb.Leafsize
		rc.Generation = sb.Generation
		rc.ChunkRootGeneration = sb.ChunkRootGeneration
		rc.CsumSize = BtrfsSuperCsumSize(&sb)
		rc.Fsid = sb.Fsid
		fmt.Fprintf(os.Stderr, "\nRC: %+v\n", rc)
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

// calcStripeLength: calculates the length of a single stripe based on the type of block group a total length and the number of stripes.
func calcStripeLength(typeFlags, length uint64, numStripes uint16) uint64 {
	var stripeSize uint64
	switch {
	case typeFlags&BTRFS_BLOCK_GROUP_RAID0 != 0:
		stripeSize = length
		stripeSize /= uint64(numStripes)
	case typeFlags&BTRFS_BLOCK_GROUP_RAID10 != 0:
		stripeSize = length * 2
		stripeSize /= uint64(numStripes)
	case typeFlags&BTRFS_BLOCK_GROUP_RAID5 != 0:
		stripeSize = length
		stripeSize /= (uint64(numStripes) - 1)
	case typeFlags&BTRFS_BLOCK_GROUP_RAID6 != 0:
		stripeSize = length
		stripeSize /= (uint64(numStripes) - 2)
	default:
		stripeSize = length
	}
	return stripeSize
}

// checkChunkRefs: checks that the chunk has a block group and that it matches the chunk
// for each stripe it also checks that the stripe exist in the device extent and that it matches
func checkChunkRefs(chunkRec *ChunkRecord, blockGroupTree *BlockGroupTree, devExtentTree *DeviceExtentTree, silent bool) bool {

	ret := true
	if i := blockGroupTree.Tree.Get(&BlockGroupRecord{CacheExtent: CacheExtent{Start: chunkRec.Offset, Size: chunkRec.Length}}); i != nil {
		blockGroupRec := i.(*BlockGroupRecord)
		if chunkRec.Length != blockGroupRec.Offset || chunkRec.Offset != blockGroupRec.Objectid || chunkRec.TypeFlags != blockGroupRec.Flags {
			if !silent {
				fmt.Fprintf(os.Stderr, "Chunk[%d, %d, %d]: length(%d), offset(%d), type(%d) mismatch with block group[%d, %d, %d]: offset(%d), objectid(%d), flags(%d)\n",
					chunkRec.Objectid,
					chunkRec.Type,
					chunkRec.Offset,
					chunkRec.Length,
					chunkRec.Offset,
					chunkRec.TypeFlags,
					blockGroupRec.Objectid,
					blockGroupRec.Type,
					blockGroupRec.Offset,
					blockGroupRec.Offset,
					blockGroupRec.Objectid,
					blockGroupRec.Flags)
				ret = false
			} else {
				chunkRec.List = nil
				chunkRec.BgRec = blockGroupRec
				//					list_del_init(&block_group_rec->list);
				//			chunk_rec->bg_rec = block_group_rec;
			}
		} else {
			if !silent {
				fmt.Fprintf(os.Stderr, "Chunk[%d, %d, %d]: length(%d), offset(%d), type(%d) is not found in block group\n",
					chunkRec.Objectid,
					chunkRec.Type,
					chunkRec.Offset,
					chunkRec.Length,
					chunkRec.Offset,
					chunkRec.TypeFlags)
				ret = false

			}
		}
	}
	length := calcStripeLength(chunkRec.TypeFlags, chunkRec.Length, chunkRec.NumStripes)
	for i := uint16(0); i < chunkRec.NumStripes; i++ {
		devid := chunkRec.Stripes[i].Devid
		offset := chunkRec.Stripes[i].Offset
		if item := devExtentTree.Tree.Get(&DeviceExtentRecord{CacheExtent: CacheExtent{Objectid: devid, Start: offset, Size: length}}); item != nil {
			devExtentRec := item.(*DeviceExtentRecord)
			if devExtentRec.Objectid != devid ||
				devExtentRec.Offset != offset ||
				devExtentRec.ChunkOffset != chunkRec.Offset ||
				devExtentRec.Length != length {
				if !silent {
					fmt.Fprintf(os.Stderr,
						"Chunk[%d, %d, %d] stripe[%d, %d] dismatch dev extent[%d, %d, %d] local [%d, %d, %d]\n",
						chunkRec.Objectid,
						chunkRec.Type,
						chunkRec.Offset,
						chunkRec.Stripes[i].Devid,
						chunkRec.Stripes[i].Offset,
						devExtentRec.Objectid,
						devExtentRec.Offset,
						devExtentRec.Length,
						devid,
						offset,
						length)
					ret = false
				} else {
					devExtentRec.ChunkList = chunkRec.Dextents.PushBack(devExtentRec)
					//									list_move(&dev_extent_rec->chunk_list,
					//					  &chunk_rec->dextents);
				}
			}
		} else {
			if !silent {
				fmt.Fprintf(os.Stderr,
					"Chunk[%d, %d, %d] stripe[%d, %d] is not found in dev extent\n",
					chunkRec.Objectid,
					chunkRec.Type,
					chunkRec.Offset,
					chunkRec.Stripes[i].Devid,
					chunkRec.Stripes[i].Offset)
				ret = false
			}
		}
	}
	return ret
}

// checkChunks: walks the chunk cache and checks the chunks references incrementing the good and bad chunk lists accordingliy
// if silent is false interates the block groups and device extent chunk oprphans list printing a could not find message
func CheckChunks(chunkCache *llrb.LLRB, blockGroupTree *BlockGroupTree, deviceExtentTree *DeviceExtentTree, good, bad *list.List, silent bool) bool {

	var ret bool
	chunkCache.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		chunkRec := i.(*ChunkRecord)
		err := checkChunkRefs(chunkRec, blockGroupTree, deviceExtentTree, silent)
		if !err {
			ret = false
			if bad != nil {
				chunkRec.List = bad.PushBack(chunkRec)
			}
		} else {
			if good != nil {
				chunkRec.List = good.PushBack(chunkRec)
			}
		}
		return true
	})
	if !silent {
		for i := blockGroupTree.Block_Groups.Front(); i != nil && i.Value != nil; i = i.Next() {
			bgRec := i.Value.(*BlockGroupRecord)
			fmt.Fprintf(os.Stderr, "Block group[%d, %d] (flags = %d) didn't find the relative chunk.\n",
				bgRec.Objectid,
				bgRec.Offset,
				bgRec.Flags)
		}
		for i := deviceExtentTree.ChunkOrphans.Front(); i != nil && i.Value != nil; i = i.Next() {
			dextRec := i.Value.(*DeviceExtentRecord)
			fmt.Fprintf(os.Stderr, "Device extent[%d, %d, %d] didn't find the relative chunk.\n",
				dextRec.Objectid,
				dextRec.Offset,
				dextRec.Length)
		}
	}
	return ret
}

//btrfsGetDeviceExtents: counts the chunks in the orphan device extents list that match this chunk and returns a list of them
func btrfsGetDeviceExtents(chunkObject uint64, orphanDevexts *list.List, retList *list.List) uint16 {

	count := uint16(0)
	for i := orphanDevexts.Front(); i != nil && i.Value != nil; i = i.Next() {
		devExt := i.Value.(*DeviceExtentRecord)
		if devExt.ChunkOffset == chunkObject {
			retList.PushBack(devExt)
			//			list_move_tail(&devext->chunk_list, retList);
			count++
		}
	}
	return count
}

type ChunkOffset struct {
	ChunkOffset uint64
	devExt      *DeviceExtentRecord
}

// ByChunkOffset sort interface
type ByChunkOffset []ChunkOffset

func (a ByChunkOffset) Len() int           { return len(a) }
func (a ByChunkOffset) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByChunkOffset) Less(i, j int) bool { return a[i].ChunkOffset < a[j].ChunkOffset }

//fixupDevExt attempt to create missing device extents using the current block group and devExt
func fixupDevExt(bg *BlockGroupRecord, devextCache *DeviceExtentTree, devExts *list.List) uint16 {
	// create array of ChunkOffsets from device extents
	chunkOffsets := make([]ChunkOffset, devextCache.Tree.Len())
	j := 0
	devextCache.Tree.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		devExt := i.(*DeviceExtentRecord)
		chunkOffsets[j] = ChunkOffset{ChunkOffset: devExt.ChunkOffset, devExt: devExt}
		j++
		return true
	})

	// find prev devExt by ChunkOffset
	chunkObject := bg.Offset
	sort.Sort(ByChunkOffset(chunkOffsets))
	var prevDevExt *DeviceExtentRecord
	for _, chunkOffset := range chunkOffsets {
		if chunkOffset.ChunkOffset < chunkObject {
			prevDevExt = chunkOffset.devExt
		}
		if chunkOffset.ChunkOffset > chunkObject {
			break
		}
	}
	// create a new DevExt after the prevDevExt at the new ChunkOffset
	tmp := *prevDevExt
	newDevExt := &tmp
	offset := chunkObject - newDevExt.ChunkOffset
	newDevExt.ChunkOffset = chunkObject
	newDevExt.CacheExtent.Start += offset
	newDevExt.Offset += offset
	newDevExt.ChunkList = devextCache.ChunkOrphans.PushBack(newDevExt)
	newDevExt.DeviceList = devextCache.DeviceOrphans.PushBack(newDevExt)
	devextCache.Tree.InsertNoReplace(newDevExt)
	devExts.PushBack(newDevExt)

	return 1
}

// btrfsRecoverChunks create the chunks by block group
func BtrfsRecoverChunks(rc *RecoverControl) bool {
	ret := true

	for i := rc.Bg.Block_Groups.Front(); i != nil && i.Value != nil; i = i.Next() {
		bg := i.Value.(*BlockGroupRecord)
		devExts := list.New()
		nstripes := btrfsGetDeviceExtents(bg.Objectid, rc.Devext.ChunkOrphans, devExts)
		if nstripes == 0 {
			// this seems to break mapping for some reason
			//			nstripes = fixupDevExt(bg, &rc.Devext, devExts)
		}
		fmt.Fprintf(os.Stderr, "BG offset %d found count %d #devExt %d\n", bg.Objectid, nstripes, devExts.Len())
		chunk := &ChunkRecord{
			Dextents:    list.New(),
			BgRec:       bg,
			CacheExtent: CacheExtent{Start: bg.Objectid, Size: bg.Offset},
			Objectid:    BTRFS_FIRST_CHUNK_TREE_OBJECTID,
			Type:        BTRFS_CHUNK_ITEM_KEY,
			Offset:      bg.Objectid,
			Generation:  bg.Generation,
			Length:      bg.Offset,
			Owner:       BTRFS_CHUNK_TREE_OBJECTID,
			StripeLen:   BTRFS_STRIPE_LEN,
			TypeFlags:   bg.Flags,
			IoWidth:     BTRFS_STRIPE_LEN,
			IoAlign:     BTRFS_STRIPE_LEN,
			SectorSize:  rc.Sectorsize,
			SubStripes:  calcSubNstripes(bg.Flags),
		}
		rc.Chunk.InsertNoReplace(chunk)
		if nstripes == 0 {
			fmt.Fprintf(os.Stderr, "No stripes %+v\n", chunk)
			chunk.List = rc.BadChunks.PushBack(chunk)
			continue
		}
		chunk.Dextents.PushBackList(devExts)

		if !btrfsVerifyDeviceExtents(bg, devExts, nstripes) {
			fmt.Fprintf(os.Stderr, "No verify %+v\n", chunk)
			continue
			chunk.List = rc.BadChunks.PushBack(chunk)
		}
		chunk.NumStripes = nstripes
		chunk.Stripes = make([]Stripe, nstripes)
		var err error
		err, ret = btrfsRebuildChunkStripes(rc, chunk)
		switch {
		case err != nil:
			chunk.List = rc.UnrepairedChunks.PushBack(chunk)
		case ret:
			chunk.List = rc.GoodChunks.PushBack(chunk)
		case !ret:
			fmt.Fprintf(os.Stderr, "Bad rebuild %+v\n", chunk)
			chunk.List = rc.BadChunks.PushBack(chunk)
		}
	}
	/*
	 * Don't worry about the lost orphan device extents, they don't
	 * have its chunk and block group, they must be the old ones that
	 * we have dropped.
	 */
	return ret
}

// calcSubNstripes caculates the number of sub stripes as 2 if this is a rad10 block 1 otherwise
func calcSubNstripes(Type uint64) uint16 {
	if Type&BTRFS_BLOCK_GROUP_RAID10 != 0 {
		return 2
	} else {
		return 1
	}
}

// btrfsRebuildOrderedMetaChunkStripes: rebuild stripes for raid block groups with existing devextents
func btrfsRebuildOrderedMetaChunkStripes(rc *RecoverControl, chunk *ChunkRecord) (error, bool) {
	var (
		start  = chunk.Offset
		end    = chunk.Offset + chunk.Length
		mirror uint32
	)
	cache := rc.EbCache.Get(&ExtentRecord{CacheExtent: CacheExtent{Start: start, Size: chunk.Length}})
	if cache == nil {
		return btrfsRebuildUnorderedChunkStripes(rc, chunk)
	}
	devExts := chunk.Dextents
	chunk.Dextents = list.New()
again:
	er := cache.(*ExtentRecord)
	index := btrfsCalcStripeIndex(chunk, er.CacheExtent.Start)
	if chunk.Stripes[index].Devid != 0 {
		goto next
	}
	for i := devExts.Front(); i != nil && i.Value != nil; i = i.Next() {
		devExt := i.Value.(*DeviceExtentRecord)
		if isExtentRecordInDeviceExtent(er, devExt, &mirror) {
			chunk.Stripes[index].Devid = devExt.Objectid
			chunk.Stripes[index].Offset = devExt.Offset
			chunk.Stripes[index].Uuid = er.Devices[mirror].Uuid
			index++
			//			list_move(&devext->chunk_list, &chunk->dextents);
			devExt.ChunkList = chunk.Dextents.PushBack(devExt)
		}
	}
next:
	start = btrfsNextStripeLogicalOffset(chunk, er.CacheExtent.Start)
	if start > end {
		goto noExtentRecord
	}
	cache = rc.EbCache.Get(&ExtentRecord{CacheExtent: CacheExtent{Start: start, Size: end - start}})
	if cache != nil {
		goto again
	}
noExtentRecord:
	if devExts.Len() == 0 {
		return nil, true
	}
	if chunk.TypeFlags&(BTRFS_BLOCK_GROUP_RAID5|BTRFS_BLOCK_GROUP_RAID6) == 0 {
		//		/* Fixme: try to recover the order by the parity block. */
		//		list_splice_tail(&devexts, &chunk->dextents);
		return errors.New("-EINVAL"), false
	}
	/* There is no data on the lost stripes, we can reorder them freely. */
	for index := uint16(0); index < chunk.NumStripes; index++ {
		if chunk.Stripes[index].Devid != 0 {
			continue
		}
		//		devext = list_first_entry(&devexts,
		//					  struct device_extent_record,
		//					   chunk_list);
		//		list_move(&devext->chunk_list, &chunk->dextents);
		if i := devExts.Front(); i != nil {
			devExt := i.Value.(*DeviceExtentRecord)
			devExt.ChunkList = chunk.Dextents.PushBack(devExt)
			chunk.Stripes[index].Devid = devExt.Objectid
			chunk.Stripes[index].Offset = devExt.Offset
			device := btrfsFindDeviceByDevid(rc.FsDevices, devExt.Objectid, 0)
			if device == nil {
				chunk.Dextents.PushBackList(devExts)
				return errors.New("-EINVAL"), false
			}
			chunk.Stripes[index].Uuid = device.Uuid
		}
	}
	return nil, true
}

// btrfsRebuildUnorderedChunkStripes rebuild stimple stripes
func btrfsRebuildUnorderedChunkStripes(rc *RecoverControl, chunk *ChunkRecord) (error, bool) {
	for i, item := uint16(0), chunk.Dextents.Front(); item != nil && item.Value != nil && i < chunk.NumStripes; i, item = i+1, item.Next() {
		devExt := item.Value.(*DeviceExtentRecord)
		chunk.Stripes[i].Devid = devExt.Objectid
		chunk.Stripes[i].Offset = devExt.Offset
		device := btrfsFindDeviceByDevid(rc.FsDevices, devExt.Objectid, 0)
		if device == nil {
			return errors.New("-EINVAL"), false
		}
		chunk.Stripes[i].Uuid = device.Uuid
	}
	return nil, true
}

// btrfsRebuildChunkStripes rebuild the rtripes fro the chunk
func btrfsRebuildChunkStripes(rc *RecoverControl, chunk *ChunkRecord) (error, bool) {
	/*
	 * All the data in the system metadata chunk will be dropped,
	 * so we need not guarantee that the data is right or not, that
	 * is we can reorder the stripes in the system metadata chunk.
	 */
	switch {
	case (chunk.TypeFlags&BTRFS_BLOCK_GROUP_METADATA) != 0 && (chunk.TypeFlags&BTRFS_ORDERED_RAID) != 0:
		return btrfsRebuildOrderedMetaChunkStripes(rc, chunk)
	case (chunk.TypeFlags&BTRFS_BLOCK_GROUP_DATA) != 0 && (chunk.TypeFlags&BTRFS_ORDERED_RAID) != 0:
		return nil, true /* Be handled after the fs is opened. */
	}
	return btrfsRebuildUnorderedChunkStripes(rc, chunk)
}

// btrfsVerifyDeviceExtents: verifies there is a device extent of the corrrect lenghth for each stripe of the block group based on block group Flags
func btrfsVerifyDeviceExtents(bg *BlockGroupRecord, devExts *list.List, nDevExts uint16) bool {
	var (
		stripeLength       uint64
		expectedNumStripes uint16
	)
	expectedNumStripes = calcNumStripes(bg.Flags)
	if expectedNumStripes != 0 && expectedNumStripes != nDevExts {
		fmt.Fprintf(os.Stderr, "Expected stripes %d Extents %d\n", expectedNumStripes, nDevExts)
		return false
	}
	stripeLength = calcStripeLength(bg.Flags, bg.Offset, nDevExts)
	for i := devExts.Front(); i != nil && i.Value != nil; i = i.Next() {
		devExt := i.Value.(*DeviceExtentRecord)
		if devExt.Length != stripeLength {
			fmt.Fprintf(os.Stderr, "Expected stripelen %d found len %d\n", stripeLength, devExt.Length)
			return false
		}
	}
	return true
}

func btrfsCalcStripeIndex(chunk *ChunkRecord, logical uint64) uint16 {
	var (
		offset                                     = logical - chunk.Offset
		numStripes, ntDataStripes, stripeNr, index uint64
	)

	stripeNr = uint64(offset / chunk.StripeLen)
	switch {
	case chunk.TypeFlags&BTRFS_BLOCK_GROUP_RAID0 != 0:
		numStripes = uint64(chunk.NumStripes)
		index = stripeNr % numStripes
	case chunk.TypeFlags&BTRFS_BLOCK_GROUP_RAID10 != 0:
		subStripes := uint64(chunk.SubStripes)
		index = stripeNr % numStripes / subStripes
		index *= subStripes
	case chunk.TypeFlags&BTRFS_BLOCK_GROUP_RAID5 != 0:
		numStripes = uint64(chunk.NumStripes)
		ntDataStripes = numStripes - 1
		index = stripeNr % ntDataStripes
		stripeNr /= ntDataStripes
		index = (index + stripeNr) % numStripes
	case chunk.TypeFlags&BTRFS_BLOCK_GROUP_RAID6 != 0:
		numStripes = uint64(chunk.NumStripes)
		ntDataStripes = numStripes - 2
		index = stripeNr % ntDataStripes
		stripeNr /= ntDataStripes
		index = (index + stripeNr) % numStripes
	}
	return uint16(index)
}

func btrfsFindDeviceByDevid(fsDevices *BtrfsFsDevices, devid uint64, instance int) *BtrfsDevice {
	numFound := 0
	for i := fsDevices.Devices.Front(); i != nil && i.Value != nil; i = i.Next() {
		dev := i.Value.(*BtrfsDevice)
		if dev.Devid == devid && numFound == instance {
			return dev
		}
		numFound++
	}
	return nil
}

func calcNumStripes(flags uint64) uint16 {

	switch {
	case flags&(BTRFS_BLOCK_GROUP_RAID0|
		BTRFS_BLOCK_GROUP_RAID10|
		BTRFS_BLOCK_GROUP_RAID5|
		BTRFS_BLOCK_GROUP_RAID6) != 0:
		return 0

	case flags&(BTRFS_BLOCK_GROUP_RAID1|
		BTRFS_BLOCK_GROUP_DUP) != 0:
		return 2
	default:
		return 1
	}

}

func isExtentRecordInDeviceExtent(er *ExtentRecord, dext *DeviceExtentRecord, mirror *uint32) bool {

	for i := uint32(0); i < er.Nmirrors; i++ {
		if er.Devices[i].Devid == dext.Objectid &&
			er.Offsets[i] >= dext.Offset &&
			er.Offsets[i] < dext.Offset+dext.Length {
			*mirror = i
			return true
		}
	}
	return false
}

// btrfsNextStripeLogicalOffset calc the logical offset which is the start of the next stripe
func btrfsNextStripeLogicalOffset(chunk *ChunkRecord, logical uint64) uint64 {

	var offset = logical - chunk.Offset
	offset /= chunk.StripeLen
	offset *= chunk.StripeLen
	offset += chunk.StripeLen
	return offset + chunk.Offset
}

// TODO:
// Proper support for multiple devices and device scan
