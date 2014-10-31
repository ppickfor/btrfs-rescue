package btrfs

import (
	"github.com/petar/GoLLRB/llrb"
	"sync"
)

type Recover_control struct {
	Verbose bool
	Yes     bool

	Csum_size             uint16
	Sectorsize            uint32
	Leafsize              uint32
	Generation            uint64
	Chunk_root_generation uint64

	Fs_devices *Btrfs_fs_devices

	Chunk    *llrb.LLRB
	Bg       Block_group_tree
	Devext   Device_extent_tree
	Eb_cache *llrb.LLRB

	Good_chunks       List_head
	Bad_chunks        List_head
	Unrepaired_chunks List_head
	Rc_lock           sync.Mutex
	Fd                int
	Fsid                  [16]uint8
}
