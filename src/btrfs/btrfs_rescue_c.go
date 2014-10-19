package btrfs

import (
	"sync"
)

type Recover_control struct {
	Verbose int
	Yes     int

	Csum_size             uint16
	Sectorsize            uint64
	Leafsize              uint64
	Generation            uint64
	Chunk_root_generation uint64

	Fs_devices *Btrfs_fs_devices

	Chunk    Cache_tree
	Bg       Block_group_tree
	Devext   Device_extent_tree
	Eb_cache Cache_tree

	Good_chunks       List_head
	Bad_chunks        List_head
	Unrepaired_chunks List_head
	Rc_lock           sync.Mutex
}
