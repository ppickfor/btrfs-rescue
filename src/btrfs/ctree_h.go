package btrfs

import ()

type btrfs_fs_devices struct {
	fsid [BTRFS_FSID_SIZE]U8 /* FS specific uuid */

	/* the device with this id has the most recent copy of the super */
	latest_devid U64
	latest_trans U64
	lowest_devid U64
	latest_bdev  int
	lowest_bdev  int
	devices      list_head
	list         list_head

	seeding int
	seed    *btrfs_fs_devices
}
type btrfs_root_backup struct {
	Tree_root     Le64
	Tree_root_gen Le64

	Chunk_root     Le64
	Chunk_root_gen Le64

	Extent_root     Le64
	Extent_root_gen Le64

	Fs_root     Le64
	Fs_root_gen Le64

	Dev_root     Le64
	Dev_root_gen Le64

	Csum_root     Le64
	Csum_root_gen Le64

	Total_bytes Le64
	Bytes_used  Le64
	Num_devices Le64
	/* future */
	Unsed_64 [4]Le64

	Tree_root_level   U8
	Chunk_root_level  U8
	Extent_root_level U8
	Fs_root_level     U8
	Dev_root_level    U8
	Csum_root_level   U8
	/* future and to align */
	Unused_8 [10]U8
}
type btrfs_dev_item struct {
	/* the internal btrfs device id */
	Devid Le64
	/* size of the device */
	Total_bytes Le64
	/* bytes used */
	Bytes_used Le64
	/* optimal io alignment for this device */
	Io_align Le32
	/* optimal io width for this device */
	Io_width Le32
	/* minimal io size for this device */
	Sector_size Le32
	/* type and info about this device */
	Type Le64
	/* expected generation for this device */
	Generation Le64
	/*
	 * starting byte of this partition on the device,
	 * to allowr for stripe alignment in the future
	 */
	Start_offset Le64
	/* grouping information for allocation decisions */
	Dev_group Le32
	/* seek speed 0-100 where 100 is fastest */
	Seek_speed U8
	/* bandwidth 0-100 where 100 is fastest */
	Bandwidth U8
	/* btrfs generated uuid for this device */
	Uuid [BTRFS_UUID_SIZE]U8
	/* uuid of FS who owns this device */
	Fsid [BTRFS_UUID_SIZE]U8
}

type Btrfs_super_block struct {
	Csum [BTRFS_CSUM_SIZE]U8
	/* the first 3 fields must match struct btrfs_header */
	Fsid   [BTRFS_FSID_SIZE]U8 /* FS specific uuid */
	Bytenr Le64                /* this block number */
	Flags  Le64

	/* allowed to be different from the btrfs_header from here own down */
	Magic      Le64
	Generation Le64
	Root       Le64
	Chunk_root Le64
	Log_root   Le64

	/* this will help find the new super based on the log root */
	Log_root_transid      Le64
	Total_bytes           Le64
	Bytes_used            Le64
	Root_dir_objectid     Le64
	Num_devices           Le64
	Sectorsize            Le32
	Nodesize              Le32
	Leafsize              Le32
	Stripesize            Le32
	Sys_chunk_array_size  Le32
	Chunk_root_generation Le64
	Compat_flags          Le64
	Compat_ro_flags       Le64
	Incompat_flags        Le64
	Csum_type             Le16
	Root_level            U8
	Chunk_root_level      U8
	Log_root_level        U8
	Dev_item              btrfs_dev_item

	Label [BTRFS_LABEL_SIZE]Char

	Cache_generation     Le64
	Uuid_tree_generation Le64

	/* future expansion */
	Reserved        [30]Le64
	Sys_chunk_array [BTRFS_SYSTEM_CHUNK_ARRAY_SIZE]U8
	Super_roots     [BTRFS_NUM_BACKUP_ROOTS]btrfs_root_backup
}
