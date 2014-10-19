package btrfs

import ()

const (
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
	 * [ 1 BTRFS_DEV_ITEM_KEY device_id ]
	 */
	BTRFS_DEV_ITEMS_OBJECTID = 1
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
