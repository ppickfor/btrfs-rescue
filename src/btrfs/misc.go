package btrfs

import (
	"math"
)

const (
	BTRFS_MAX_MIRRORS       = 3
	BTRFS_SIZE              = 4096
	BTRFS_SUPER_INFO_OFFSET = (64 * 1024)
	BTRFS_SUPER_INFO_SIZE   = 4096
	BTRFS_STRIPE_LEN        = (64 * 1024)

	BTRFS_SUPER_MIRROR_MAX   = 3
	BTRFS_SUPER_MIRROR_SHIFT = 12

	BTRFS_CSUM_SIZE = 32

	BTRFS_FSID_SIZE = 16
	BTRFS_UUID_SIZE = 16

	BTRFS_SYSTEM_CHUNK_ARRAY_SIZE = 2048
	BTRFS_LABEL_SIZE              = 256

	/* tag for the radix tree of block groups in ram */
	BTRFS_BLOCK_GROUP_DATA     = (1 << 0)
	BTRFS_BLOCK_GROUP_SYSTEM   = (1 << 1)
	BTRFS_BLOCK_GROUP_METADATA = (1 << 2)
	BTRFS_BLOCK_GROUP_RAID0    = (1 << 3)
	BTRFS_BLOCK_GROUP_RAID1    = (1 << 4)
	BTRFS_BLOCK_GROUP_DUP      = (1 << 5)
	BTRFS_BLOCK_GROUP_RAID10   = (1 << 6)
	BTRFS_BLOCK_GROUP_RAID5    = (1 << 7)
	BTRFS_BLOCK_GROUP_RAID6    = (1 << 8)
	BTRFS_ORDERED_RAID         = (BTRFS_BLOCK_GROUP_RAID0 |
		BTRFS_BLOCK_GROUP_RAID10 |
		BTRFS_BLOCK_GROUP_RAID5 |
		BTRFS_BLOCK_GROUP_RAID6)
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
	BTRFS_BALANCE_OBJECTID uint64 = -4 & math.MaxUint64

	/* oprhan objectid for tracking unlinked/truncated files */
	BTRFS_ORPHAN_OBJECTID uint64 = -5 & math.MaxUint64

	/* does write ahead logging to speed up fsyncs */
	BTRFS_TREE_LOG_OBJECTID       uint64 = -6 & math.MaxUint64
	BTRFS_TREE_LOG_FIXUP_OBJECTID uint64 = -7 & math.MaxUint64

	/* space balancing */
	BTRFS_TREE_RELOC_OBJECTID      uint64 = -8 & math.MaxUint64
	BTRFS_DATA_RELOC_TREE_OBJECTID uint64 = -9 & math.MaxUint64

	/*
	 * extent checksums all have this objectid
	 * this allows them to share the logging tree
	 * for fsyncs
	 */
	BTRFS_EXTENT_CSUM_OBJECTID uint64 = -10 & math.MaxUint64

	/* For storing free space cache */
	BTRFS_FREE_SPACE_OBJECTID uint64 = -11 & math.MaxUint64

	/*
	 * The inode number assigned to the special inode for sotring
	 * free ino cache
	 */
	BTRFS_FREE_INO_OBJECTID uint64 = -12 & math.MaxUint64

	/* dummy objectid represents multiple objectids */
	BTRFS_MULTIPLE_OBJECTIDS uint64 = -255 & math.MaxUint64

	/*
	 * All files have objectids in this range.
	 */
	BTRFS_FIRST_FREE_OBJECTID              = 256
	BTRFS_LAST_FREE_OBJECTID        uint64 = -256 & math.MaxUint64
	BTRFS_FIRST_CHUNK_TREE_OBJECTID        = 256

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
