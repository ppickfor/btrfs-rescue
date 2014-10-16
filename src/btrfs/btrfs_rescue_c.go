package btrfs

import (

)

type Recover_control struct {
	Verbose int
	Yes     int

	Csum_size             U16
	Sectorsize            U32
	Leafsize              U32
	Generation            U64
	Chunk_root_generation U64

	Fs_devices *btrfs_fs_devices

	Chunk    cache_tree
	Bg       block_group_tree
	Devext   device_extent_tree
	Eb_cache cache_tree

	Good_chunks       list_head
	Bad_chunks        list_head
	Unrepaired_chunks list_head
	//pthread_mutex_t rc_lock;
}