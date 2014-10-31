package main

import (
	"btrfs"
	"bytes"
	"flag"
	"fmt"
	//	"log"
	"os"
	//	"syscall"
	"encoding/binary"
	"encoding/hex"
	//	"strings"
	"bufio"
	"github.com/petar/GoLLRB/llrb"
	"log"
	"regexp"
	"strconv"
)

const (
	APP_VERSION = "0.1"

	DEFAULT_BLOCKS = 1024
	DEFAULT_DEVICE = "/dev/sda"
)

// The flag package provides a default help printer via -h switch
var (
	versionFlag     *bool   = flag.Bool("v", false, "Print the version number.")
	deviceFlag      *string = flag.String("d", DEFAULT_DEVICE, "The device to scan")
	blocksFlag      *int64  = flag.Int64("n", DEFAULT_BLOCKS, "The number of BTRFS_SIZE blocks to read")
	startblocksFlag *int64  = flag.Int64("s", 0, "The number of BTRFS_SIZE blocks to start at")
	fsidFlag        *string = flag.String("f", "671e120db12f497da4c42e41b3d30fec", "The Hex fsid of the file system to match")
	bytenrFlag      *string = flag.String("b", "", "The previously scanned block bytenrs")
)

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func decodeObjectID(id int64) {
	// Object IDS
	/* holds pointers to all of the tree roots */
	switch id {
	case btrfs.BTRFS_ROOT_TREE_OBJECTID:
		fmt.Println("BTRFS_ROOT_TREE_OBJECTID")

	/* stores information about which extents are in use, and reference counts */
	case btrfs.BTRFS_EXTENT_TREE_OBJECTID:
		fmt.Println("BTRFS_EXTENT_TREE_OBJECTID")

	/*
	 * chunk tree stores translations from logical -> physical block numbering
	 * the super block points to the chunk tree
	 */
	case btrfs.BTRFS_CHUNK_TREE_OBJECTID:
		fmt.Println("BTRFS_CHUNK_TREE_OBJECTID")

	/*
	 * stores information about which areas of a given device are in use.
	 * one per device.  The tree of tree roots points to the device tree
	 */
	case btrfs.BTRFS_DEV_TREE_OBJECTID:
		fmt.Println("BTRFS_DEV_TREE_OBJECTID")

	/* one per subvolume, storing files and directories */
	case btrfs.BTRFS_FS_TREE_OBJECTID:
		fmt.Println("BTRFS_FS_TREE_OBJECTID")

	/* directory objectid inside the root tree */
	case btrfs.BTRFS_ROOT_TREE_DIR_OBJECTID:
		fmt.Println("BTRFS_ROOT_TREE_DIR_OBJECTID")
	/* holds checksums of all the data extents */
	case btrfs.BTRFS_CSUM_TREE_OBJECTID:
		fmt.Println("BTRFS_CSUM_TREE_OBJECTID")
	case btrfs.BTRFS_QUOTA_TREE_OBJECTID:
		fmt.Println("BTRFS_QUOTA_TREE_OBJECTID")

	/* for storing items that use the BTRFS_UUID_KEY* */
	case btrfs.BTRFS_UUID_TREE_OBJECTID:
		fmt.Println("BTRFS_UUID_TREE_OBJECTID")

	/* for storing balance parameters in the root tree */
	case btrfs.BTRFS_BALANCE_OBJECTID:
		fmt.Println("BTRFS_BALANCE_OBJECTID")

	/* oprhan objectid for tracking unlinked/truncated files */
	case btrfs.BTRFS_ORPHAN_OBJECTID:
		fmt.Println("BTRFS_ORPHAN_OBJECTID")

	/* does write ahead logging to speed up fsyncs */
	case btrfs.BTRFS_TREE_LOG_OBJECTID:
		fmt.Println("BTRFS_TREE_LOG_OBJECTID")
	case btrfs.BTRFS_TREE_LOG_FIXUP_OBJECTID:
		fmt.Println("BTRFS_TREE_LOG_FIXUP_OBJECTID")

	/* space balancing */
	case btrfs.BTRFS_TREE_RELOC_OBJECTID:
		fmt.Println("BTRFS_TREE_RELOC_OBJECTID")
	case btrfs.BTRFS_DATA_RELOC_TREE_OBJECTID:
		fmt.Println("BTRFS_DATA_RELOC_TREE_OBJECTID")

	/*
	 * extent checksums all have this objectid
	 * this allows them to share the logging tree
	 * for fsyncs
	 */
	case btrfs.BTRFS_EXTENT_CSUM_OBJECTID:
		fmt.Println("BTRFS_EXTENT_CSUM_OBJECTID")

	/* For storing free space cache */
	case btrfs.BTRFS_FREE_SPACE_OBJECTID:
		fmt.Println("BTRFS_FREE_SPACE_OBJECTID")

	/*
	 * The inode number assigned to the special inode for sotring
	 * free ino cache
	 */
	case btrfs.BTRFS_FREE_INO_OBJECTID:
		fmt.Println("BTRFS_FREE_INO_OBJECTID")

	/* dummy objectid represents multiple objectids */
	case btrfs.BTRFS_MULTIPLE_OBJECTIDS:
		fmt.Println("BTRFS_MULTIPLE_OBJECTIDS")

	/*
	 * All files have objectids in this range.
	 */
	//	case btrfs.BTRFS_FIRST_FREE_OBJECTID:
	//		fmt.Println("BTRFS_FIRST_FREE_OBJECTID")
	//	case btrfs.BTRFS_LAST_FREE_OBJECTID:
	//		fmt.Println("BTRFS_LAST_FREE_OBJECTID")
	//	case btrfs.BTRFS_FIRST_CHUNK_TREE_OBJECTID:
	//		fmt.Println("BTRFS_FIRST_CHUNK_TREE_OBJECTID")
	//
	//	/*
	//	 * the device items go into the chunk tree.  The key is in the form
	//	 * [ 1 BTRFS_DEV_ITEM_KEY device_id ]
	//	 */
	//	case btrfs.BTRFS_DEV_ITEMS_OBJECTID:
	//		fmt.Println("BTRFS_DEV_ITEMS_OBJECTID")
	default:
		if (id > 0 && id > btrfs.BTRFS_FIRST_FREE_OBJECTID) || (id < 0 && id < btrfs.BTRFS_LAST_FREE_OBJECTID) {
			fmt.Printf("Numbered Object %08x\n", uint64(id))
		} else {
			fmt.Println("UNKNOWN OBJECTID")

		}
	}
}

func decodeKeyID(id uint8) {
	/*
	 * inode items have the data typically returned from stat and store other
	 * info about object characteristics.  There is one for every file and dir in
	 * the FS
	 */
	switch id {
	case btrfs.BTRFS_INODE_ITEM_KEY:
		fmt.Println("BTRFS_INODE_ITEM_KEY")
	case btrfs.BTRFS_INODE_REF_KEY:
		fmt.Println("BTRFS_INODE_REF_KEY")
	case btrfs.BTRFS_INODE_EXTREF_KEY:
		fmt.Println("BTRFS_INODE_EXTREF_KEY")
	case btrfs.BTRFS_XATTR_ITEM_KEY:
		fmt.Println("BTRFS_XATTR_ITEM_KEY")
	case btrfs.BTRFS_ORPHAN_ITEM_KEY:
		fmt.Println("BTRFS_ORPHAN_ITEM_KEY")

	case btrfs.BTRFS_DIR_LOG_ITEM_KEY:
		fmt.Println("BTRFS_DIR_LOG_ITEM_KEY")
	case btrfs.BTRFS_DIR_LOG_INDEX_KEY:
		fmt.Println("BTRFS_DIR_LOG_INDEX_KEY")
	/*
	 * dir items are the name -> inode pointers in a directory.  There is one
	 * for every name in a directory.
	 */
	case btrfs.BTRFS_DIR_ITEM_KEY:
		fmt.Println("BTRFS_DIR_ITEM_KEY")
	case btrfs.BTRFS_DIR_INDEX_KEY:
		fmt.Println("BTRFS_DIR_INDEX_KEY")

	/*
	 * extent data is for file data
	 */
	case btrfs.BTRFS_EXTENT_DATA_KEY:
		fmt.Println("BTRFS_EXTENT_DATA_KEY")

	/*
	 * csum items have the checksums for data in the extents
	 */
	case btrfs.BTRFS_CSUM_ITEM_KEY:
		fmt.Println("BTRFS_CSUM_ITEM_KEY")
	/*
	 * extent csums are stored in a separate tree and hold csums for
	 * an entire extent on disk.
	 */
	case btrfs.BTRFS_EXTENT_CSUM_KEY:
		fmt.Println("BTRFS_EXTENT_CSUM_KEY")

	/*
	 * root items point to tree roots.  There are typically in the root
	 * tree used by the super block to find all the other trees
	 */
	case btrfs.BTRFS_ROOT_ITEM_KEY:
		fmt.Println("BTRFS_ROOT_ITEM_KEY")

	/*
	 * root backrefs tie subvols and snapshots to the directory entries that
	 * reference them
	 */
	case btrfs.BTRFS_ROOT_BACKREF_KEY:
		fmt.Println("BTRFS_ROOT_BACKREF_KEY")

	/*
	 * root refs make a fast index for listing all of the snapshots and
	 * subvolumes referenced by a given root.  They point directly to the
	 * directory item in the root that references the subvol
	 */
	case btrfs.BTRFS_ROOT_REF_KEY:
		fmt.Println("BTRFS_ROOT_REF_KEY")

	/*
	 * extent items are in the extent map tree.  These record which blocks
	 * are used, and how many references there are to each block
	 */
	case btrfs.BTRFS_EXTENT_ITEM_KEY:
		fmt.Println("BTRFS_EXTENT_ITEM_KEY")

	/*
	 * The same as the BTRFS_EXTENT_ITEM_KEY, except it's metadata we already know
	 * the length, so we save the level in key->offset instead of the length.
	 */
	case btrfs.BTRFS_METADATA_ITEM_KEY:
		fmt.Println("BTRFS_METADATA_ITEM_KEY")

	case btrfs.BTRFS_TREE_BLOCK_REF_KEY:
		fmt.Println("BTRFS_TREE_BLOCK_REF_KEY")

	case btrfs.BTRFS_EXTENT_DATA_REF_KEY:
		fmt.Println("BTRFS_EXTENT_DATA_REF_KEY")

	/* old style extent backrefs */
	case btrfs.BTRFS_EXTENT_REF_V0_KEY:
		fmt.Println("BTRFS_EXTENT_REF_V0_KEY")

	case btrfs.BTRFS_SHARED_BLOCK_REF_KEY:
		fmt.Println("BTRFS_SHARED_BLOCK_REF_KEY")

	case btrfs.BTRFS_SHARED_DATA_REF_KEY:
		fmt.Println("BTRFS_SHARED_DATA_REF_KEY")

	/*
	 * block groups give us hints into the extent allocation trees.  Which
	 * blocks are free etc etc
	 */
	case btrfs.BTRFS_BLOCK_GROUP_ITEM_KEY:
		fmt.Println("BTRFS_BLOCK_GROUP_ITEM_KEY")

	case btrfs.BTRFS_DEV_EXTENT_KEY:
		fmt.Println("BTRFS_DEV_EXTENT_KEY")
	case btrfs.BTRFS_DEV_ITEM_KEY:
		fmt.Println("BTRFS_DEV_ITEM_KEY")
	case btrfs.BTRFS_CHUNK_ITEM_KEY:
		fmt.Println("BTRFS_CHUNK_ITEM_KEY")

	case btrfs.BTRFS_BALANCE_ITEM_KEY:
		fmt.Println("BTRFS_BALANCE_ITEM_KEY")

	/*
	 * quota groups
	 */
	case btrfs.BTRFS_QGROUP_STATUS_KEY:
		fmt.Println("BTRFS_QGROUP_STATUS_KEY")
	case btrfs.BTRFS_QGROUP_INFO_KEY:
		fmt.Println("BTRFS_QGROUP_INFO_KEY")
	case btrfs.BTRFS_QGROUP_LIMIT_KEY:
		fmt.Println("BTRFS_QGROUP_LIMIT_KEY")
	case btrfs.BTRFS_QGROUP_RELATION_KEY:
		fmt.Println("BTRFS_QGROUP_RELATION_KEY")

	/*
	 * Persistently stores the io stats in the device tree.
	 * One key for all stats, (0, BTRFS_DEV_STATS_KEY, devid).
	 */
	case btrfs.BTRFS_DEV_STATS_KEY:
		fmt.Println("BTRFS_DEV_STATS_KEY")

	/*
	 * Persistently stores the device replace state in the device tree.
	 * The key is built like this: (0, BTRFS_DEV_REPLACE_KEY, 0).
	 */
	case btrfs.BTRFS_DEV_REPLACE_KEY:
		fmt.Println("BTRFS_DEV_REPLACE_KEY")

		/*
		 * Stores items that allow to quickly map UUIDs to something else.
		 * These items are part of the filesystem UUID tree.
		 * The key is built like this:
		 * (UUID_upper_64_bits, BTRFS_UUID_KEY*, UUID_lower_64_bits).
		 */
		//#if BTRFS_UUID_SIZE != 16
		//#error "UUID items require case btrfs.BTRFS_UUID_SIZE:
		fmt.Println("BTRFS_UUID_SIZE")
	//#endif
	case btrfs.BTRFS_UUID_KEY_SUBVOL:
		fmt.Println("BTRFS_UUID_KEY_SUBVOL")
	case btrfs.BTRFS_UUID_KEY_RECEIVED_SUBVOL:
		fmt.Println("BTRFS_UUID_KEY_RECEIVED_SUBVOL")
		/* received subvols */

	/*
	 * string items are for debugging.  They just store a short string of
	 * data in the FS
	 */
	case btrfs.BTRFS_STRING_ITEM_KEY:
		fmt.Println("BTRFS_STRING_ITEM_KEY")
		//	default:
		//		fmt.Println("UNKNOWN ITEM KEY ID")

	}
}

type Cache_Extent_record struct {
	Objectid uint64
	Start    uint64
	Size     uint64
}

type Extent_record struct {
	Cache      Cache_Extent_record
	Generation uint64
	Csum       [btrfs.BTRFS_CSUM_SIZE]uint8
	devices    [btrfs.BTRFS_MAX_MIRRORS](*btrfs.Btrfs_device)
	offsets    [btrfs.BTRFS_MAX_MIRRORS]uint64
	nmirrors   uint32
}

// create a new extenet record from the block header and the Leaf size in recover control
func NewExtent_record(header *btrfs.Btrfs_header, rc *btrfs.Recover_control) *Extent_record {
	return &Extent_record{
		Cache:      Cache_Extent_record{Start: header.Bytenr, Size: uint64(rc.Leafsize)},
		Generation: header.Generation,
		Csum:       header.Csum,
	}
}
func (x *Extent_record) Less(y llrb.Item) bool {
	return (x.Cache.Start + x.Cache.Size) <= y.(*Extent_record).Cache.Start
}

type Block_group_record struct {
	Cache      Cache_Extent_record
	List       btrfs.List_head
	Generation uint64
	Objectid   uint64
	Type       uint8
	Offset     uint64
	Flags      uint64
}

// create a new block group record from the a block heade, a block group ietem key and the items byte buffer
func NewBlock_group_record(header *btrfs.Btrfs_header, item *btrfs.Btrfs_item, itemsBuf []byte) *Block_group_record {

	key := item.Key
	itemPtr := itemsBuf[item.Offset:]
	bytereader := bytes.NewReader(itemPtr)
	bgItem := new(btrfs.Btrfs_block_group_item)
	_ = binary.Read(bytereader, binary.LittleEndian, bgItem)
	//	fmt.Printf("bgItem: %+v\n", bgItem)
	this := &Block_group_record{
		Cache:      Cache_Extent_record{Start: key.Objectid, Size: key.Offset},
		Generation: header.Generation,
		Objectid:   key.Objectid,
		Type:       key.Type,
		Offset:     key.Offset,
		Flags:      bgItem.Flags,
		//		List:	btrfs.List_head{Next: &List, Prev: &List},
	}
	this.List = btrfs.List_head{Next: &this.List, Prev: &this.List}
	return this
}
func (x *Block_group_record) Less(y llrb.Item) bool {
	return (x.Cache.Start + x.Cache.Size) <= y.(*Block_group_record).Cache.Start
}

type Device_extent_record struct {
	Cache       Cache_Extent_record
	Chunk_list  btrfs.List_head
	Device_list btrfs.List_head
	Generation  uint64
	Objectid    uint64
	Type        uint8

	Offset         uint64
	Chunk_objectid uint64
	Chunk_offset   uint64
	Length         uint64
}

// create a new device extent record from the a block header, a device extent item key and the items byte buffer
func NewDevice_extent_record(header *btrfs.Btrfs_header, item *btrfs.Btrfs_item, itemsBuf []byte) *Device_extent_record {

	key := item.Key
	itemPtr := itemsBuf[item.Offset:]
	bytereader := bytes.NewReader(itemPtr)
	deItem := new(btrfs.Btrfs_dev_extent)
	_ = binary.Read(bytereader, binary.LittleEndian, deItem)
	//	fmt.Printf("bgItem: %+v\n", bgItem)
	this := &Device_extent_record{
		Cache:          Cache_Extent_record{Objectid: key.Objectid, Start: key.Offset, Size: deItem.Length},
		Generation:     header.Generation,
		Objectid:       key.Objectid,
		Type:           key.Type,
		Offset:         key.Offset,
		Chunk_objectid: deItem.Chunk_Objectid,
		Chunk_offset:   deItem.Chunk_Offset,
		Length:         deItem.Length,
	}
	this.Chunk_list = btrfs.List_head{Next: &this.Chunk_list, Prev: &this.Chunk_list}
	this.Device_list = btrfs.List_head{Next: &this.Device_list, Prev: &this.Device_list}
	return this
}

func (x *Device_extent_record) Less(y llrb.Item) bool {
	switch {
	case x.Cache.Objectid < y.(*Device_extent_record).Cache.Objectid:
		return true
	case x.Cache.Objectid > y.(*Device_extent_record).Cache.Objectid:
		return false
	}
	return (x.Cache.Start + x.Cache.Size) <= y.(*Device_extent_record).Cache.Start
}

type treeBlock struct {
	bytent    uint64
	byteblock []byte
}

func list_add_tail(xnew, head *btrfs.List_head) {
	__list_add(xnew, head.Prev, head)
}
func __list_add(xnew,
	prev,
	next *btrfs.List_head) {
	next.Prev = xnew
	xnew.Next = next
	xnew.Prev = prev
	prev.Next = xnew
}
func __list_del(prev *btrfs.List_head, next *btrfs.List_head) {
	next.Prev = prev
	prev.Next = next
}
func list_del_init(entry *btrfs.List_head) {

	__list_del(entry.Prev, entry.Next)
	entry.Next = entry
	entry.Prev = entry
}

// reads super block into sb at offset sb_bytenr
// if super_recover is != 0 then read superblock backups and find latest generation
func main() {

	var rc btrfs.Recover_control

	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}
	btrfs.Init_recover_control(&rc, true, false)
	btrfs.Recover_prepare(&rc, *deviceFlag)
	//	header := new(btrfs.Btrfs_header)
	nonhex := regexp.MustCompile("[^0-9a-fA-F]*")
	fsidstr := nonhex.ReplaceAllString(*fsidFlag, "")
	fsid, err := hex.DecodeString(fsidstr)
	bytenrs := make([]uint64, 5874688)
	off0 := btrfs.Btrfs_sb_offset(0)
	off1 := btrfs.Btrfs_sb_offset(1)
	off2 := btrfs.Btrfs_sb_offset(2)

	treeBlockchan := make(chan treeBlock, 2)

	if err != nil {
		fmt.Errorf("fsidFlag:%v is not hex, %v\n", *fsidFlag, err)
		os.Exit(1)
	}
	//	var items *[]btrfs.Btrfs_item
	if len(*bytenrFlag) == 0 {
		// read treeblock from disk and pass to channel for processing
		go func() {
			treeBlock := new(treeBlock)

			size := uint64(rc.Leafsize)
			fsid := rc.Fsid[:]
		loop:
			for bytenr := uint64(*startblocksFlag) * size; bytenr < uint64(*blocksFlag)*size; bytenr += size {
				switch bytenr {
				case off0, off1, off2:
					continue loop
				default:
					byteblock := make([]byte, size)
					ok, err := btrfs.Btrfs_read_treeblock(rc.Fd, bytenr, size, fsid, &byteblock)
					if err != nil {
						fmt.Println(err)
						break loop
					}
					if ok {
						treeBlock.bytent = bytenr
						treeBlock.byteblock = byteblock
						treeBlockchan <- *treeBlock
					}
				}
			}
			close(treeBlockchan)
		}()

	} else {
		file, err := os.Open(*bytenrFlag)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for i := 0; scanner.Scan(); i++ {
			bytenr, err := strconv.ParseUint(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			bytenrs[i] = bytenr
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
		// read treeblock from disk and pass to channel for processing
		go func() {
			treeBlock := new(treeBlock)
			size := uint64(rc.Leafsize)
			fsid := rc.Fsid[:]
		loop:
			for _, bytenr := range bytenrs {
				switch bytenr {
				case off0, off1, off2:
					continue loop
				default:
					byteblock := make([]byte, size)
					ok, err := btrfs.Btrfs_read_treeblock(rc.Fd, bytenr, size, fsid, &byteblock)
					if err != nil {
						fmt.Println(err)
						break loop
					}
					if ok {
						treeBlock.bytent = bytenr
						treeBlock.byteblock = byteblock
						treeBlockchan <- *treeBlock
					}
				}
			}
			close(treeBlockchan)
		}()

	}
	// process treeblock from goroutine via channel
	for treeBlock := range treeBlockchan {
		//			treeBlock = treeBlock
		//		fmt.Printf("treeBlock: %+v\n", treeBlock)
		//		fmt.Printf("from chan treeblock: @%d, %v\n", treeBlock.bytent, treeBlock.byteblock[0:4])
		detailBlock(&treeBlock, &rc, fsid)
	}
	fmt.Printf("\nAll %d Extent buffers\n", rc.Eb_cache.Len())
	rc.Eb_cache.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		fmt.Printf("%+v\n", i)
		return true
	})
	fmt.Printf("\nAll %d Block Groups\n", rc.Bg.Tree.Len())
	rc.Bg.Tree.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		fmt.Printf("%+v\n", i)
		return true
	})
	fmt.Printf("\nAll %d Device Extents\n", rc.Devext.Tree.Len())
	rc.Devext.Tree.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		fmt.Printf("%+v\n", i)
		return true
	})
}

// print some info while scanning for metadata of the fsid
func processBlock(bytenr uint64, header *btrfs.Btrfs_header, rc *btrfs.Recover_control, fsid []byte) {

	if !btrfs.Is_super_block_address(uint64(bytenr)) {
		if btrfs.Btrfs_read_header(rc.Fd, header, bytenr) {
			if bytes.Equal(header.Fsid[:], fsid) {
				fmt.Printf("Btrfs header @%v:\n%+v\n\n", bytenr, header)
				decodeObjectID(int64(header.Owner))

				if ret, items := btrfs.Btrfs_read_items(rc.Fd, header.Nritems, bytenr); ret {
					for _, item := range items {
						decodeKeyID(item.Key.Type)
					}
				}
			}
		}
	}
}

// build cache structures for selected block records and items at physical bytenr
func detailBlock(treeBlock *treeBlock, rc *btrfs.Recover_control, fsid []byte) {

	byteblock := treeBlock.byteblock
	header := new(btrfs.Btrfs_header)
	bytereader := bytes.NewReader(byteblock)
	_ = binary.Read(bytereader, binary.LittleEndian, header)
	bytenr := treeBlock.bytent

	tree := rc.Eb_cache

	er := NewExtent_record(header, rc)

again:
	if exists := tree.Get(er); exists != nil {
		if exists.(*Extent_record).Generation > er.Generation {
			//				fmt.Printf("Exists:%d\r", tree.Len())
			return
		}
		if exists.(*Extent_record).Generation == er.Generation {
			if exists.(*Extent_record).Cache.Start != er.Cache.Start ||
				exists.(*Extent_record).Cache.Size != er.Cache.Size ||
				bytes.Compare(exists.(*Extent_record).Csum[:], er.Csum[:]) != 0 {
				//							exists but different
				fmt.Printf("detailBlock: Exists but dif %+v\n", er)
				return
			} else {
				//					fmt.Printf("Mirror:%d\r", tree.Len())
				//				fmt.Printf("mirror %+v\n", exists.(*Extent_record))
				if exists.(*Extent_record).nmirrors < 3 {
					exists.(*Extent_record).devices[exists.(*Extent_record).nmirrors] = nil
					exists.(*Extent_record).offsets[exists.(*Extent_record).nmirrors] = bytenr
				}
				exists.(*Extent_record).nmirrors++
				tree.ReplaceOrInsert(exists)
				return
			}
		}
		tree.Delete(er)
		//			fmt.Printf("Worse:%d\n", tree.Len())
		goto again
	}
	tree.InsertNoReplace(er)
	// process items in new treeblock
	if header.Level == 0 {
		// Leaf
		items := make([]btrfs.Btrfs_item, header.Nritems)
		infoByteBlock := byteblock[len(byteblock)-bytereader.Len():]
		_ = binary.Read(bytereader, binary.LittleEndian, items)
		//			fmt.Printf("Leaf @%08x: Items: %d\r", bytenr, leaf.Header.Nritems)
		switch header.Owner {
		case btrfs.BTRFS_EXTENT_TREE_OBJECTID, btrfs.BTRFS_DEV_TREE_OBJECTID:
			/* different tree use different generation */
			//			if header.Generation <= rc.Generation {
			extract_metadata_record(rc, header, items, infoByteBlock)
			//			}
		case btrfs.BTRFS_CHUNK_TREE_OBJECTID:
			//			if header.Generation <= rc.Chunk_root_generation {
			extract_metadata_record(rc, header, items, infoByteBlock)
			//			}
		}

	} else {
		// Node
		//			node := new(btrfs.Btrfs_node)
		//			node.Header = *header
		//			node.Ptrs = make([]btrfs.Btrfs_key_ptr, header.Nritems)
		//			_ = binary.Read(bytereader, binary.LittleEndian, node.Ptrs)
		//			fmt.Printf("Node @%08x: Pointers: %d\r", bytenr, node.Header.Nritems)
	}

}

// iterate items of the current block and proccess each blockgroup, chunk and dev extent adding them to caches
func extract_metadata_record(rc *btrfs.Recover_control, header *btrfs.Btrfs_header, items []btrfs.Btrfs_item, itemBuf []byte) {

	for _, item := range items {
		//		fmt.Printf("Leaf Item: %+v\n", &leaf.Items[i])

		switch item.Key.Type {
		case btrfs.BTRFS_BLOCK_GROUP_ITEM_KEY:
			//			pthread_mutex_lock(&rc->rc_lock);
			//			fmt.Printf("BLOCK Group Item: %+v\n", &leaf.Items[i])
			process_block_group_item(&rc.Bg, header, &item, itemBuf)
			//			pthread_mutex_unlock(&rc->rc_lock);
			//			break;
		case btrfs.BTRFS_CHUNK_ITEM_KEY:
			//			pthread_mutex_lock(&rc->rc_lock);
			//			ret = process_chunk_item(&rc->chunk, leaf, &key, i);
			//			pthread_mutex_unlock(&rc->rc_lock);
			//			break;
		case btrfs.BTRFS_DEV_EXTENT_KEY:
			//			pthread_mutex_lock(&rc->rc_lock);
			process_device_extent_item(&rc.Devext, header, &item, itemBuf)
			//			pthread_mutex_unlock(&rc->rc_lock);
			//			break;
		}
		//		if (ret)
		//			break;
	}

}

// create a new block group record and update the cache by latest generation for each blockgroup item
func process_block_group_item(
	bg_cache *btrfs.Block_group_tree,
	header *btrfs.Btrfs_header,
	item *btrfs.Btrfs_item,
	itemBuf []byte) {

	rec := NewBlock_group_record(header, item, itemBuf)
again:
	if exists := bg_cache.Tree.Get(rec); exists != nil {
		/*check the generation and replace if needed*/
		if exists.(*Block_group_record).Generation > rec.Generation {
			return
		}
		//		if (exist->generation > rec->generation)
		//			goto free_out;

		if exists.(*Block_group_record).Generation == rec.Generation {
			//			int offset = offsetof(struct block_group_record,
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
			//			goto free_out;
			fmt.Printf("process_block_group_item: same generation %+v\n", rec)
			return

		}
		list_del_init(&exists.(*Block_group_record).List)
		bg_cache.Tree.Delete(exists)
		/*
		 * We must do seach again to avoid the following cache.
		 * /--old bg 1--//--old bg 2--/
		 *        /--new bg--/
		 */
		goto again

	}
	bg_cache.Tree.InsertNoReplace(rec)
	list_add_tail(&rec.List, &bg_cache.Block_Groups)
	//	bg_cache.InsertNoReplace(rec)
}

// create a new device extent record and update the cache by latest generation
func process_device_extent_item(devext_cache *btrfs.Device_extent_tree,
	header *btrfs.Btrfs_header,
	item *btrfs.Btrfs_item,
	itemBuf []byte) {
	rec := NewDevice_extent_record(header, item, itemBuf)
again:
	if exists := devext_cache.Tree.Get(rec); exists != nil {
		/*check the generation and replace if needed*/
		if exists.(*Device_extent_record).Generation > rec.Generation {
			return
		}
		//		if (exist->generation > rec->generation)
		//			goto free_out;

		if exists.(*Device_extent_record).Generation == rec.Generation {
			// should not happen
			fmt.Printf("process_device_extent_item: same generation %+v\n", rec)
			return
		}
		list_del_init(&exists.(*Device_extent_record).Chunk_list)
		list_del_init(&exists.(*Device_extent_record).Device_list)
		devext_cache.Tree.Delete(exists)
		/*
		 * We must do seach again to avoid the following cache.
		 * /--old bg 1--//--old bg 2--/
		 *        /--new bg--/
		 */
		goto again

	}
	devext_cache.Tree.ReplaceOrInsert(rec)
	//	devext_cache.InsertNoReplace(rec)
	list_add_tail(&rec.Chunk_list, &devext_cache.Chunk_orphans)
	list_add_tail(&rec.Device_list, &devext_cache.Device_orphans)
}
