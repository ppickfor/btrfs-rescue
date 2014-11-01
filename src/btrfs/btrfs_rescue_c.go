package btrfs

import (
	"container/list"
	"github.com/petar/GoLLRB/llrb"
	"sync"
)

type RecoverControl struct {
	Verbose             bool
	Yes                 bool
	CsumSize            uint16
	Sectorsize          uint32
	Leafsize            uint32
	Generation          uint64
	ChunkRootGeneration uint64
	FsDevices           *BtrfsFsDevices
	Chunk               *llrb.LLRB
	Bg                  BlockGroupTree
	Devext              DeviceExtentTree
	EbCache             *llrb.LLRB
	GoodChunks          *list.List
	BadChunks           *list.List
	UnrepairedChunks    *list.List
	RcLock              sync.Mutex
	Fd                  int
	Fsid                [16]uint8
}
