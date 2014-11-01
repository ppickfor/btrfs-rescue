package btrfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"syscall"
)

const (
	BTRFS_MAX_MIRRORS       = 3
	BTRFS_SIZE              = 4096
	BTRFS_SUPER_INFO_OFFSET = (64 * 1024)
	BTRFS_SUPER_INFO_SIZE   = 4096

	BTRFS_SUPER_MIRROR_MAX   = 3
	BTRFS_SUPER_MIRROR_SHIFT = 12

	BTRFS_CSUM_SIZE = 32

	BTRFS_FSID_SIZE = 16
	BTRFS_UUID_SIZE = 16

	BTRFS_SYSTEM_CHUNK_ARRAY_SIZE = 2048
	BTRFS_LABEL_SIZE              = 256

	/*
	 * just in case we somehow lose the roots and are not able to mount,
	 * we store an array of the roots from previous transactions
	 * in the super.
	 */
	BTRFS_NUM_BACKUP_ROOTS    = 4
	BTRFS_MAGIC               = 0x4D5F53665248425F /* ascii _BHRfS_M, no null */
	BTRFS_HEADER_FLAG_WRITTEN = (1 << 0)
	BTRFS_HEADER_FLAG_RELOC   = (1 << 1)
	BTRFS_SUPER_FLAG_SEEDING  = (1 << 32)
	BTRFS_SUPER_FLAG_METADUMP = (1 << 33)
	// Object IDS
	/* holds pointers to all of the tree roots */
	BTRFS_ROOT_TREE_OBJECTID = 1

	/* stores information about which extents are in use, and reference counts */
	BTRFS_EXTENT_TREE_OBJECTID = 2

	/*
	 * chunk tree stores translations from logical -> physical block numbering
	 * the super block points to the chunk tree
	 */
	BTRFS_CHUNK_TREE_OBJECTID = 3

	/*
	 * stores information about which areas of a given device are in use.
	 * one per device.  The tree of tree roots points to the device tree
	 */
	BTRFS_DEV_TREE_OBJECTID = 4

	/* one per subvolume, storing files and directories */
	BTRFS_FS_TREE_OBJECTID = 5

	/* directory objectid inside the root tree */
	BTRFS_ROOT_TREE_DIR_OBJECTID = 6
	/* holds checksums of all the data extents */
	BTRFS_CSUM_TREE_OBJECTID  = 7
	BTRFS_QUOTA_TREE_OBJECTID = 8

	/* for storing items that use the BTRFS_UUID_KEY* */
	BTRFS_UUID_TREE_OBJECTID = 9

	/* for storing balance parameters in the root tree */
	BTRFS_BALANCE_OBJECTID = -4

	/* oprhan objectid for tracking unlinked/truncated files */
	BTRFS_ORPHAN_OBJECTID = -5

	/* does write ahead logging to speed up fsyncs */
	BTRFS_TREE_LOG_OBJECTID       = -6
	BTRFS_TREE_LOG_FIXUP_OBJECTID = -7

	/* space balancing */
	BTRFS_TREE_RELOC_OBJECTID      = -8
	BTRFS_DATA_RELOC_TREE_OBJECTID = -9

	/*
	 * extent checksums all have this objectid
	 * this allows them to share the logging tree
	 * for fsyncs
	 */
	BTRFS_EXTENT_CSUM_OBJECTID = -10

	/* For storing free space cache */
	BTRFS_FREE_SPACE_OBJECTID = -11

	/*
	 * The inode number assigned to the special inode for sotring
	 * free ino cache
	 */
	BTRFS_FREE_INO_OBJECTID = -12

	/* dummy objectid represents multiple objectids */
	BTRFS_MULTIPLE_OBJECTIDS = -255

	/*
	 * All files have objectids in this range.
	 */
	BTRFS_FIRST_FREE_OBJECTID       = 256
	BTRFS_LAST_FREE_OBJECTID        = -256
	BTRFS_FIRST_CHUNK_TREE_OBJECTID = 256

	/*
	 * the device items go into the chunk tree.  The key is in the form
	 * [ 1 BTRFS_DEV_ITEM_KEY deviceId ]
	 */
	BTRFS_DEV_ITEMS_OBJECTID = 1

	/*
	 * inode items have the data typically returned from stat and store other
	 * info about object characteristics.  There is one for every file and dir in
	 * the FS
	 */
	BTRFS_INODE_ITEM_KEY   = 1
	BTRFS_INODE_REF_KEY    = 12
	BTRFS_INODE_EXTREF_KEY = 13
	BTRFS_XATTR_ITEM_KEY   = 24
	BTRFS_ORPHAN_ITEM_KEY  = 48

	BTRFS_DIR_LOG_ITEM_KEY  = 60
	BTRFS_DIR_LOG_INDEX_KEY = 72
	/*
	 * dir items are the name -> inode pointers in a directory.  There is one
	 * for every name in a directory.
	 */
	BTRFS_DIR_ITEM_KEY  = 84
	BTRFS_DIR_INDEX_KEY = 96

	/*
	 * extent data is for file data
	 */
	BTRFS_EXTENT_DATA_KEY = 108

	/*
	 * csum items have the checksums for data in the extents
	 */
	BTRFS_CSUM_ITEM_KEY = 120
	/*
	 * extent csums are stored in a separate tree and hold csums for
	 * an entire extent on disk.
	 */
	BTRFS_EXTENT_CSUM_KEY = 128

	/*
	 * root items point to tree roots.  There are typically in the root
	 * tree used by the super block to find all the other trees
	 */
	BTRFS_ROOT_ITEM_KEY = 132

	/*
	 * root backrefs tie subvols and snapshots to the directory entries that
	 * reference them
	 */
	BTRFS_ROOT_BACKREF_KEY = 144

	/*
	 * root refs make a fast index for listing all of the snapshots and
	 * subvolumes referenced by a given root.  They point directly to the
	 * directory item in the root that references the subvol
	 */
	BTRFS_ROOT_REF_KEY = 156

	/*
	 * extent items are in the extent map tree.  These record which blocks
	 * are used, and how many references there are to each block
	 */
	BTRFS_EXTENT_ITEM_KEY = 168

	/*
	 * The same as the BTRFS_EXTENT_ITEM_KEY, except it's metadata we already know
	 * the length, so we save the level in key->offset instead of the length.
	 */
	BTRFS_METADATA_ITEM_KEY = 169

	BTRFS_TREE_BLOCK_REF_KEY = 176

	BTRFS_EXTENT_DATA_REF_KEY = 178

	/* old style extent backrefs */
	BTRFS_EXTENT_REF_V0_KEY = 180

	BTRFS_SHARED_BLOCK_REF_KEY = 182

	BTRFS_SHARED_DATA_REF_KEY = 184

	/*
	 * block groups give us hints into the extent allocation trees.  Which
	 * blocks are free etc etc
	 */
	BTRFS_BLOCK_GROUP_ITEM_KEY = 192

	BTRFS_DEV_EXTENT_KEY = 204
	BTRFS_DEV_ITEM_KEY   = 216
	BTRFS_CHUNK_ITEM_KEY = 228

	BTRFS_BALANCE_ITEM_KEY = 248

	/*
	 * quota groups
	 */
	BTRFS_QGROUP_STATUS_KEY   = 240
	BTRFS_QGROUP_INFO_KEY     = 242
	BTRFS_QGROUP_LIMIT_KEY    = 244
	BTRFS_QGROUP_RELATION_KEY = 246

	/*
	 * Persistently stores the io stats in the device tree.
	 * One key for all stats, (0, BTRFS_DEV_STATS_KEY, devid).
	 */
	BTRFS_DEV_STATS_KEY = 249

	/*
	 * Persistently stores the device replace state in the device tree.
	 * The key is built like this: (0, BTRFS_DEV_REPLACE_KEY, 0).
	 */
	BTRFS_DEV_REPLACE_KEY = 250

	/*
	 * Stores items that allow to quickly map UUIDs to something else.
	 * These items are part of the filesystem UUID tree.
	 * The key is built like this:
	 * (UUIDUpper_64Bits, BTRFS_UUID_KEY*, UUIDLower_64Bits).
	 */
	//#if BTRFS_UUID_SIZE != 16
	//#error "UUID items require BTRFS_UUID_SIZE == 16!"
	//#endif
	BTRFS_UUID_KEY_SUBVOL          = 251 /* for UUIDs assigned to subvols */
	BTRFS_UUID_KEY_RECEIVED_SUBVOL = 252 /* for UUIDs assigned to
	 * received subvols */

	/*
	 * string items are for debugging.  They just store a short string of
	 * data in the FS
	 */
	BTRFS_STRING_ITEM_KEY = 253
)

type (
	U8  byte
	U16 uint16
	U32 uint32
	U64 uint64

	Char byte
	Le16 uint16
	Le32 uint32
	Le64 uint64
)

// read header struct from fd at bytenr
func BtrfsReadHeader(fd int, header *BtrfsHeader, bytenr uint64) bool {

	var size = binary.Size(header)
	var byteheader = make([]byte, size)
	var bytebr = bytes.NewReader(byteheader)

	ret, _ := syscall.Pread(fd, byteheader, int64(bytenr))
	if ret < size {
		return false
	}
	_ = binary.Read(bytebr, binary.LittleEndian, header)

	return true

}

// read items structs from fd at bytenr
func BtrfsReadItems(fd int, n uint32, bytenr uint64) (bool, []BtrfsItem) {
	header := new(BtrfsHeader)
	item := new(BtrfsItem)
	myitems := make([]BtrfsItem, n)
	var hsize = binary.Size(header)
	var isize = binary.Size(item)
	var byteitems = make([]byte, uint32(isize)*n)
	var bytebr = bytes.NewReader(byteitems)

	ret, _ := syscall.Pread(fd, byteitems, int64(bytenr+uint64(hsize)))
	if ret < isize {
		return false, nil
	}
	err := binary.Read(bytebr, binary.LittleEndian, myitems)
	if err != nil {
		return false, nil
	} else {
		return true, myitems
	}

}

func BtrfsReadTreeblock(fd int, bytenr uint64, size uint64, fsid []byte, byteblock *[]byte) (bool, error) {
	//	csum := uint32(0)

	ret, err := syscall.Pread(fd, *byteblock, int64(bytenr))
	if err != nil {
		return false, os.NewSyscallError("pread64", err)
	}
	if ret < int(size) {
		err = errors.New("Pread: not all bytes read")
		return false, err
	}
	//	fmt.Printf("byteblock: %v\n", byteblock)
	if bytes.Equal((*byteblock)[32:32+16], fsid) {
		return true, nil
	}
	return false, nil
}
