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

	Fs_devices *Btrfs_fs_devices

	Chunk    Cache_tree
	Bg       Block_group_tree
	Devext   Device_extent_tree
	Eb_cache Cache_tree

	Good_chunks       List_head
	Bad_chunks        List_head
	Unrepaired_chunks List_head
	//pthread_mutex_t rc_lock;
}