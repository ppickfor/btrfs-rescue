// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs=true -debug-define=true btrfsH.go
package btrfs

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"sync"

	"github.com/monnand/GoLLRB/llrb"
)

type Mutex struct {
	Lock uint64
}
type RbNode struct {
	Next *RbNode
}
type RbRoot struct {
	Next *RbNode
}
type CacheExtent struct {
	Objectid uint64
	Start    uint64
	Size     uint64
}

//type *list.List struct {
//	Next **list.List
//	Prev **list.List
//}
type ExtentIoTree struct {
	State *llrb.LLRB
	Cache *llrb.LLRB
	Lru   *list.List
	Size  uint64
}
type BlockGroupRecord struct {
	CacheExtent
	List       *list.Element
	Generation uint64
	Objectid   uint64
	Type       uint8
	Offset     uint64
	Flags      uint64
}

// NewBlockGroupRecord create a new block group record from the a block heade, a block group ietem key and the items byte buffer
func NewBlockGroupRecord(generation uint64, item *BtrfsItem, itemsBuf []byte) *BlockGroupRecord {

	key := item.Key
	itemPtr := itemsBuf[item.Offset:]
	bytereader := bytes.NewReader(itemPtr)
	bgItem := new(BtrfsBlockGroupItem)
	_ = binary.Read(bytereader, binary.LittleEndian, bgItem)
	//	fmt.Printf("bgItem: %+v\n", bgItem)
	this := &BlockGroupRecord{
		CacheExtent: CacheExtent{Start: key.Objectid, Size: key.Offset},
		Generation:  generation,
		Objectid:    key.Objectid,
		Type:        key.Type,
		Offset:      key.Offset,
		Flags:       bgItem.Flags,
	}
	return this
}

func (x *BlockGroupRecord) Less(y llrb.Item) bool {
	return (x.Start + x.Size) <= y.(*BlockGroupRecord).Start
}

type BlockGroupTree struct {
	Tree         *llrb.LLRB
	Block_Groups *list.List
}
type DeviceRecord struct {
	Node       RbNode
	Devid      uint64
	Generation uint64
	Objectid   uint64
	Type       uint8
	Offset     uint64
	TotalByte  uint64
	ByteUsed   uint64
	RealUsed   uint64
}
type Stripe struct {
	Devid  uint64
	Offset uint64
	Uuid   [16]uint8
}
type ChunkRecord struct {
	CacheExtent
	List       *list.Element
	Dextents   *list.List
	BgRec      *BlockGroupRecord
	Generation uint64
	Objectid   uint64
	Type       uint8
	Offset     uint64
	Owner      uint64
	Length     uint64
	TypeFlags  uint64
	StripeLen  uint64
	NumStripes uint16
	SubStripes uint16
	IoAlign    uint32
	IoWidth    uint32
	SectorSize uint32
	Stripes    []Stripe
}

func NewChunkRecord(generation uint64, item *BtrfsItem, itemsBuf []byte) *ChunkRecord {
	key := item.Key
	itemPtr := itemsBuf[item.Offset:]
	bytereader := bytes.NewReader(itemPtr)
	chunkItem := new(BtrfsChunk)
	_ = binary.Read(bytereader, binary.LittleEndian, chunkItem)
	//	fmt.Printf("chunkItem: %+v\n", chunkItem)
	this := &ChunkRecord{
		CacheExtent: CacheExtent{Start: key.Offset, Size: chunkItem.Length},
		//		List:        list.Element(),
		Dextents: list.New(),
		//		BgRec:      nil,
		Generation: generation,
		Objectid:   key.Objectid,
		Type:       key.Type,
		Offset:     key.Offset,
		Length:     chunkItem.Length,
		Owner:      chunkItem.Owner,
		StripeLen:  chunkItem.StripeLen,
		TypeFlags:  chunkItem.Type,
		IoWidth:    chunkItem.IoWidth,
		IoAlign:    chunkItem.IoAlign,
		SectorSize: chunkItem.SectorSize,
		NumStripes: chunkItem.NumStripes,
		SubStripes: chunkItem.SubStripes,
		Stripes:    make([]Stripe, chunkItem.NumStripes),
	}
	_ = binary.Read(bytereader, binary.LittleEndian, &this.Stripes)
	return this
}
func (x *ChunkRecord) Less(y llrb.Item) bool {
	return (x.Start + x.Size) <= y.(*ChunkRecord).Start
}

type DeviceExtentRecord struct {
	CacheExtent
	ChunkList     *list.Element
	DeviceList    *list.Element
	Generation    uint64
	Objectid      uint64
	Type          uint8
	Offset        uint64
	ChunkObjectid uint64
	ChunkOffset   uint64
	Length        uint64
}

// NewDeviceExtentRecord create a new device extent record from the a block header, a device extent item key and the items byte buffer
func NewDeviceExtentRecord(generation uint64, item *BtrfsItem, itemsBuf []byte) *DeviceExtentRecord {

	key := item.Key
	itemPtr := itemsBuf[item.Offset:]
	bytereader := bytes.NewReader(itemPtr)
	deItem := new(BtrfsDevExtent)
	_ = binary.Read(bytereader, binary.LittleEndian, deItem)
	//	fmt.Printf("bgItem: %+v\n", bgItem)
	this := &DeviceExtentRecord{
		CacheExtent:   CacheExtent{Objectid: key.Objectid, Start: key.Offset, Size: deItem.Length},
		Generation:    generation,
		Objectid:      key.Objectid,
		Type:          key.Type,
		Offset:        key.Offset,
		ChunkObjectid: deItem.Chunk_Objectid,
		ChunkOffset:   deItem.Chunk_Offset,
		Length:        deItem.Length,
	}
	return this
}

// Less compares device extents first by Objectid (device) the start range in thelook left red blakc tree
func (x *DeviceExtentRecord) Less(y llrb.Item) bool {
	cmp := y.(*DeviceExtentRecord)
	switch {
	case x.CacheExtent.Objectid < cmp.CacheExtent.Objectid:
		return true
	case x.CacheExtent.Objectid > cmp.CacheExtent.Objectid:
		return false
	}
	return (x.CacheExtent.Start + x.CacheExtent.Size) <= cmp.CacheExtent.Start
}

type DeviceExtentTree struct {
	Tree          *llrb.LLRB
	ChunkOrphans  *list.List
	DeviceOrphans *list.List
}
type RootInfo struct {
	RbNode     RbNode
	SortNode   RbNode
	RootId     uint64
	RootOffset uint64
	Flags      uint64
	RefTree    uint64
	DirId      uint64
	TopId      uint64
	Gen        uint64
	Ogen       uint64
	Otime      uint64
	Uuid       [16]uint8
	Puuid      [16]uint8
	Ruuid      [16]uint8
	Path       *int8
	Name       *int8
	FullPath   *int8
	Deleted    int32
}
type CmdStruct struct {
	Token    *int8
	Fn       *[0]byte
	Usagestr **int8
	Next     *CmdGroup
	Hidden   int32
}
type CmdGroup struct {
	Usagestr **int8
	Infostr  *int8
	Commands [0]byte
}
type BtrfsDiskKey struct {
	Objectid uint64
	Type     uint8
	Offset   uint64
}
type BtrfsKey struct {
	Objectid uint64
	Type     uint8
	Offset   uint64
}
type BtrfsMappingTree struct {
	Tree *llrb.LLRB
}
type BtrfsDevItem struct {
	Devid       uint64
	TotalBytes  uint64
	BytesUsed   uint64
	IoAlign     uint32
	IoWidth     uint32
	SectorSize  uint32
	Type        uint64
	Generation  uint64
	StartOffset uint64
	DevGroup    uint32
	SeekSpeed   uint8
	Bandwidth   uint8
	Uuid        [16]uint8
	Fsid        [16]uint8
}
type BtrfsStripe struct {
	Devid  uint64
	Offset uint64
	Uuid   [16]uint8
}
type BtrfsChunk struct {
	Length     uint64
	Owner      uint64
	StripeLen  uint64
	Type       uint64
	IoAlign    uint32
	IoWidth    uint32
	SectorSize uint32
	NumStripes uint16
	SubStripes uint16
}
type BtrfsFreeSpaceEntry struct {
	Offset uint64
	Bytes  uint64
	Type   uint8
}
type BtrfsFreeSpaceHeader struct {
	Location   BtrfsDiskKey
	Generation uint64
	Entries    uint64
	Bitmaps    uint64
}
type BtrfsHeader struct {
	Csum       [32]uint8
	Fsid       [16]uint8
	Bytenr     uint64
	Flags      uint64
	TreeUuid   [16]uint8
	Generation uint64
	Owner      uint64
	Nritems    uint32
	Level      uint8
}
type BtrfsRootBackup struct {
	TreeRoot        uint64
	TreeRootGen     uint64
	ChunkRoot       uint64
	ChunkRootGen    uint64
	ExtentRoot      uint64
	ExtentRootGen   uint64
	FsRoot          uint64
	FsRootGen       uint64
	DevRoot         uint64
	DevRootGen      uint64
	CsumRoot        uint64
	CsumRootGen     uint64
	TotalBytes      uint64
	BytesUsed       uint64
	NumDevices      uint64
	Unsed_64        [4]uint64
	TreeRootLevel   uint8
	ChunkRootLevel  uint8
	ExtentRootLevel uint8
	FsRootLevel     uint8
	DevRootLevel    uint8
	CsumRootLevel   uint8
	Unused_8        [10]uint8
}
type BtrfsSuperBlock struct {
	Csum                [32]uint8
	Fsid                [16]uint8
	Bytenr              uint64
	Flags               uint64
	Magic               uint64
	Generation          uint64
	Root                uint64
	ChunkRoot           uint64
	LogRoot             uint64
	LogRootTransid      uint64
	TotalBytes          uint64
	BytesUsed           uint64
	RootDirObjectid     uint64
	NumDevices          uint64
	Sectorsize          uint32
	Nodesize            uint32
	Leafsize            uint32
	Stripesize          uint32
	SysChunkArraySize   uint32
	ChunkRootGeneration uint64
	CompatFlags         uint64
	CompatRoFlags       uint64
	IncompatFlags       uint64
	CsumType            uint16
	RootLevel           uint8
	ChunkRootLevel      uint8
	LogRootLevel        uint8
	DevItem             BtrfsDevItem
	Label               [256]int8
	CacheGeneration     uint64
	UuidTreeGeneration  uint64
	Reserved            [30]uint64
	SysChunkArray       [2048]uint8
	SuperRoots          [4]BtrfsRootBackup
}
type BtrfsItem struct {
	Key    BtrfsDiskKey
	Offset uint32
	Size   uint32
}
type BtrfsLeaf struct {
	Header BtrfsHeader
	Items  []BtrfsItem
}
type BtrfsKeyPtr struct {
	Key        BtrfsDiskKey
	Blockptr   uint64
	Generation uint64
}
type BtrfsNode struct {
	Header BtrfsHeader
	Ptrs   []BtrfsKeyPtr
}
type BtrfsPath struct {
	Nodes [8]*ExtentBuffer
	Slots [8]int32
	Locks [8]int32
	Reada int32
	Level int32
}
type BtrfsExtentItem struct {
	Refs       uint64
	Generation uint64
	Flags      uint64
}
type BtrfsExtentItemV0 struct {
	Refs uint64
}
type BtrfsTreeBlockInfo struct {
	Key   BtrfsDiskKey
	Level uint8
}
type BtrfsExtentDataRef struct {
	Root     uint64
	Objectid uint64
	Offset   uint64
	Count    uint64
}
type BtrfsSharedDataRef struct {
	Count uint64
}
type BtrfsExtentInlineRef struct {
	Type   uint8
	Offset uint64
}
type BtrfsExtentRefV0 struct {
	Root       uint64
	Generation uint64
	Objectid   uint64
	Count      uint64
}

/* dev extents record free space on individual devices.  The owner
 * field points back to the chunk allocation mapping tree that allocated
 * the extent.  The chunk tree uuid field is a way to double check the owner
 */
type BtrfsDevExtent struct {
	Chunk_Tree     uint64
	Chunk_Objectid uint64
	Chunk_Offset   uint64
	Length         uint64
	Chunk_TreeUuid [16]uint8
}
type BtrfsInodeRef struct {
	Index uint64
	Len   uint16
}
type BtrfsInodeExtref struct {
	ParentObjectid uint64
	Index          uint64
	NameLen        uint16
	Name           [0]byte
}
type BtrfsTimespec struct {
	Sec  uint64
	Nsec uint64
}
type BtrfsInodeItem struct {
	Generation uint64
	Transid    uint64
	Size       uint64
	Nbytes     uint64
	Group      uint64
	Nlink      uint64
	Uid        uint64
	Gid        uint64
	Mode       uint64
	Rdev       uint64
	Flags      uint64
	Sequence   uint64
	Reserved   [4]uint64
	Atime      BtrfsTimespec
	Ctime      BtrfsTimespec
	Mtime      BtrfsTimespec
	Otime      BtrfsTimespec
}
type BtrfsDirLogItem struct {
	End uint64
}
type BtrfsDirItem struct {
	Location BtrfsDiskKey
	Transid  uint64
	DataLen  uint16
	NameLen  uint16
	Type     uint8
}
type BtrfsRootItemV0 struct {
	Inode        BtrfsInodeItem
	Generation   uint64
	RootDirid    uint64
	Bytenr       uint64
	ByteLimit    uint64
	BytesUsed    uint64
	LastSnapshot uint64
	Flags        uint64
	Refs         uint64
	DropProgress BtrfsDiskKey
	DropLevel    uint8
	Level        uint8
}
type BtrfsRootItem struct {
	Inode        BtrfsInodeItem
	Generation   uint64
	RootDirid    uint64
	Bytenr       uint64
	ByteLimit    uint64
	BytesUsed    uint64
	LastSnapshot uint64
	Flags        uint64
	Refs         uint64
	DropProgress BtrfsDiskKey
	DropLevel    uint8
	Level        uint8
	GenerationV2 uint64
	Uuid         [16]uint8
	ParentUuid   [16]uint8
	ReceivedUuid [16]uint8
	Ctransid     uint64
	Otransid     uint64
	Stransid     uint64
	Rtransid     uint64
	Ctime        BtrfsTimespec
	Otime        BtrfsTimespec
	Stime        BtrfsTimespec
	Rtime        BtrfsTimespec
	Reserved     [8]uint64
}
type BtrfsRootRef struct {
	Dirid    uint64
	Sequence uint64
	Len      uint16
}
type BtrfsFileExtentItem struct {
	Generation    uint64
	RamBytes      uint64
	Compression   uint8
	Encryption    uint8
	OtherEncoding uint16
	Type          uint8
	DiskBytenr    uint64
	DiskNumBytes  uint64
	Offset        uint64
	NumBytes      uint64
}
type BtrfsCsumItem struct {
	Csum uint8
}
type BtrfsQgroupStatusItem struct {
	Version    uint64
	Generation uint64
	Flags      uint64
	Scan       uint64
}
type BtrfsBlockGroupItem struct {
	Used     uint64
	Objectid uint64
	Flags    uint64
}
type BtrfsQgroupInfoItem struct {
	Generation           uint64
	Referenced           uint64
	ReferencedCompressed uint64
	Exclusive            uint64
	ExclusiveCompressed  uint64
}
type BtrfsQgroupLimitItem struct {
	Flags         uint64
	MaxReferenced uint64
	MaxExclusive  uint64
	RsvReferenced uint64
	RsvExclusive  uint64
}
type BtrfsSpaceInfo struct {
	Flags       uint64
	TotalBytes  uint64
	BytesUsed   uint64
	BytesPinned uint64
	Full        int32
	List        *list.List
}
type BtrfsBlockGroupCache struct {
	Cache        CacheExtent
	Key          BtrfsKey
	Item         BtrfsBlockGroupItem
	SpaceInfo    *BtrfsSpaceInfo
	FreeSpaceCtl *BtrfsFreeSpaceCtl
	Pinned       uint64
	Flags        uint64
	Cached       int32
	Ro           int32
}
type BtrfsFsInfo struct {
	Fsid                   [16]uint8
	ChunkTreeUuid          [16]uint8
	FsRoot                 *BtrfsRoot
	ExtentRoot             *BtrfsRoot
	TreeRoot               *BtrfsRoot
	ChunkRoot              *BtrfsRoot
	DevRoot                *BtrfsRoot
	CsumRoot               *BtrfsRoot
	QuotaRoot              *BtrfsRoot
	FsRootTree             RbRoot
	LogRootTree            *BtrfsRoot
	ExtentCache            ExtentIoTree
	FreeSpaceCache         ExtentIoTree
	BlockGroupCache        ExtentIoTree
	PinnedExtents          ExtentIoTree
	PendingDel             ExtentIoTree
	ExtentIns              ExtentIoTree
	MappingTree            BtrfsMappingTree
	Generation             uint64
	LastTransCommitted     uint64
	AvailDataAllocBits     uint64
	AvailMetadataAllocBits uint64
	AvailSystemAllocBits   uint64
	DataAllocProfile       uint64
	MetadataAllocProfile   uint64
	SystemAllocProfile     uint64
	AllocStart             uint64
	RunningTransaction     *BtrfsTransHandle
	SuperCopy              *BtrfsSuperBlock
	FsMutex                Mutex
	SuperBytenr            uint64
	TotalPinned            uint64
	ExtentOps              *BtrfsExtentOps
	DirtyCowonlyRoots      *list.List
	RecowEbs               *list.List
	FsDevices              *BtrfsFsDevices
	SpaceInfo              *list.List
	SystemAllocs           int32
	FreeExtentHook         *[0]byte
	FsckExtentCache        **llrb.LLRB
	CorruptBlocks          **llrb.LLRB
}
type BtrfsRoot struct {
	Node           *ExtentBuffer
	CommitRoot     *ExtentBuffer
	RootItem       BtrfsRootItem
	RootKey        BtrfsKey
	FsInfo         *BtrfsFsInfo
	Objectid       uint64
	LastTrans      uint64
	Sectorsize     uint64
	Nodesize       uint64
	Leafsize       uint64
	Stripesize     uint64
	RefCows        int32
	TrackDirty     int32
	Type           uint64
	HighestInode   uint64
	LastInodeAlloc uint64
	DirtyList      *list.List
	RbNode         RbNode
}

func NewFakeBtrfsRoot(rc *RecoverControl) *BtrfsRoot {
	return &BtrfsRoot{
		FsInfo: &BtrfsFsInfo{
			FsDevices: rc.FsDevices,
			MappingTree: BtrfsMappingTree{
				Tree: llrb.New(),
			},
		},
	}

}

type ExtentState struct {
	Node     CacheExtent
	Start    uint64
	End      uint64
	Refs     int32
	State    uint64
	Xprivate uint64
}
type ExtentBuffer struct {
	CacheNode CacheExtent
	Start     uint64
	DevBytenr uint64
	Len       uint64
	Tree      *ExtentIoTree
	Lru       *list.List
	Recow     *list.List
	Refs      int32
	Flags     int32
	Fd        int32
	Data      []byte
}
type BtrfsFreeSpace struct {
	Index  RbNode
	Offset uint64
	Bytes  uint64
	Bitmap *uint64
	List   *list.List
}
type BtrfsFreeSpaceCtl struct {
	FreeSpaceOffset RbRoot
	FreeSpace       uint64
	ExtentsThresh   int32
	FreeExtents     int32
	TotalBitmaps    int32
	Unit            int32
	Start           uint64
	Private         *byte
}
type BtrfsIoctlVolArgs struct {
	Fd   int64
	Name [4088]int8
}
type BtrfsQgroupLimit struct {
	Flags         uint64
	MaxReferenced uint64
	MaxExclusive  uint64
	RsvReferenced uint64
	RsvExclusive  uint64
}
type BtrfsQgroupInherit struct {
	Flags         uint64
	NumGroups     uint64
	NumRefCopies  uint64
	NumExclCopies uint64
	Lim           BtrfsQgroupLimit
	Qgroups       [0]uint64
}
type BtrfsIoctlQgroupLimitArgs struct {
	Qgroupid uint64
	Lim      BtrfsQgroupLimit
}
type BtrfsIoctlVolArgsV2 struct {
	Fd      int64
	Transid uint64
	Flags   uint64
	Anon0   [32]byte
	Name    [4040]int8
}
type BtrfsScrubProgress struct {
	DataExtentsScrubbed uint64
	TreeExtentsScrubbed uint64
	DataBytesScrubbed   uint64
	TreeBytesScrubbed   uint64
	ReadErrors          uint64
	CsumErrors          uint64
	VerifyErrors        uint64
	NoCsum              uint64
	CsumDiscards        uint64
	SuperErrors         uint64
	MallocErrors        uint64
	UncorrectableErrors uint64
	CorrectedErrors     uint64
	LastPhysical        uint64
	UnverifiedErrors    uint64
}
type BtrfsIoctlScrubArgs struct {
	Devid    uint64
	Start    uint64
	End      uint64
	Flags    uint64
	Progress BtrfsScrubProgress
	Unused   [109]uint64
}
type BtrfsIoctlDevReplaceStartParams struct {
	Srcdevid                  uint64
	ContReadingFromSrcdevMode uint64
	SrcdevName                [1025]uint8
	TgtdevName                [1025]uint8
}
type BtrfsIoctlDevReplaceStatusParams struct {
	ReplaceState               uint64
	Progress_1000              uint64
	TimeStarted                uint64
	TimeStopped                uint64
	NumWriteErrors             uint64
	NumUncorrectableReadErrors uint64
}
type BtrfsIoctlDevReplaceArgs struct {
	Cmd    uint64
	Result uint64
	Anon0  [2072]byte
	Spare  [64]uint64
}
type BtrfsIoctlDevInfoArgs struct {
	Devid      uint64
	Uuid       [16]uint8
	BytesUsed  uint64
	TotalBytes uint64
	Unused     [379]uint64
	Path       [1024]uint8
}
type BtrfsIoctlFsInfoArgs struct {
	MaxId      uint64
	NumDevices uint64
	Fsid       [16]uint8
	Reserved   [124]uint64
}
type BtrfsBalanceArgs struct {
	Profiles uint64
	Usage    uint64
	Devid    uint64
	Pstart   uint64
	Pend     uint64
	Vstart   uint64
	Vend     uint64
	Target   uint64
	Flags    uint64
	Limit    uint64
	Unused   [7]uint64
}
type BtrfsBalanceProgress struct {
	Expected   uint64
	Considered uint64
	Completed  uint64
}
type BtrfsIoctlBalanceArgs struct {
	Flags  uint64
	State  uint64
	Data   BtrfsBalanceArgs
	Meta   BtrfsBalanceArgs
	Sys    BtrfsBalanceArgs
	Stat   BtrfsBalanceProgress
	Unused [72]uint64
}
type BtrfsIoctlSearchKey struct {
	TreeId      uint64
	MinObjectid uint64
	MaxObjectid uint64
	MinOffset   uint64
	MaxOffset   uint64
	MinTransid  uint64
	MaxTransid  uint64
	MinType     uint64
	MaxType     uint64
	NrItems     uint64
	Unused      uint64
	Unused1     uint64
	Unused2     uint64
	Unused3     uint64
	Unused4     uint64
}
type BtrfsIoctlSearchHeader struct {
	Transid  uint64
	Objectid uint64
	Offset   uint64
	Type     uint64
	Len      uint64
}
type BtrfsIoctlSearchArgs struct {
	Key BtrfsIoctlSearchKey
	Buf [3976]int8
}
type BtrfsIoctlInoLookupArgs struct {
	Treeid   uint64
	Objectid uint64
	Name     [4080]int8
}
type BtrfsIoctlDefragRangeArgs struct {
	Start        uint64
	Len          uint64
	Flags        uint64
	ExtentThresh uint64
	CompressType uint64
	Unused       [4]uint64
}
type BtrfsIoctlSpaceInfo struct {
	Flags      uint64
	TotalBytes uint64
	UsedBytes  uint64
}
type BtrfsIoctlSpaceArgs struct {
	SpaceSlots  uint64
	TotalSpaces uint64
	Spaces      [0]BtrfsIoctlSpaceInfo
}
type BtrfsDataContainer struct {
	BytesLeft    uint64
	BytesMissing uint64
	ElemCnt      uint64
	ElemMissed   uint64
	Val          [0]uint64
}
type BtrfsIoctlInoPathArgs struct {
	Inum     uint64
	Size     uint64
	Reserved [4]uint64
	Fspath   uint64
}
type BtrfsIoctlLogicalInoArgs struct {
	Logical  uint64
	Size     uint64
	Reserved [4]uint64
	Inodes   uint64
}
type BtrfsIoctlTimespec struct {
	Sec  uint64
	Nsec uint64
}
type BtrfsIoctlReceivedSubvolArgs struct {
	Uuid     [16]int8
	Stransid uint64
	Rtransid uint64
	Stime    BtrfsIoctlTimespec
	Rtime    BtrfsIoctlTimespec
	Flags    uint64
	Reserved [16]uint64
}
type BtrfsIoctlSendArgs struct {
	SendFd            int64
	CloneSourcesCount uint64
	CloneSources      *uint64
	ParentRoot        uint64
	Flags             uint64
	Reserved          [4]uint64
}
type BtrfsIoctlGetDevStats struct {
	Devid  uint64
	Items  uint64
	Flags  uint64
	Values [5]uint64
	Unused [121]uint64
}
type BtrfsIoctlQuotaCtlArgs struct {
	Cmd    uint64
	Status uint64
}
type BtrfsIoctlQuotaRescanArgs struct {
	Flags    uint64
	Progress uint64
	Reserved [6]uint64
}
type BtrfsIoctlQgroupAssignArgs struct {
	Assign uint64
	Src    uint64
	Dst    uint64
}
type BtrfsIoctlQgroupCreateArgs struct {
	Create   uint64
	Qgroupid uint64
}
type BtrfsIoctlCloneRangeArgs struct {
	SrcFd      int64
	SrcOffset  uint64
	SrcLength  uint64
	DestOffset uint64
}
type VmaShared struct {
	TreeNode int32
}
type VmAreaStruct struct {
	Pgoff  uint64
	Start  uint64
	End    uint64
	Shared VmaShared
}
type Page struct {
	Index uint64
}
type RadixTreeRoot struct{}
type BtrfsCorruptBlock struct {
	Cache CacheExtent
	Key   BtrfsKey
	Level int32
}
type BtrfsStreamHeader struct {
	Magic   [13]int8
	Version uint64
}
type BtrfsCmdHeader struct {
	Len uint64
	Cmd uint16
	Crc uint64
}
type BtrfsTlvHeader struct {
	Type uint16
	Len  uint16
}
type SubvolInfo struct {
	RootId       uint64
	Uuid         [16]uint8
	ParentUuid   [16]uint8
	ReceivedUuid [16]uint8
	Ctransid     uint64
	Otransid     uint64
	Stransid     uint64
	Rtransid     uint64
	Path         *int8
}
type SubvolUuidSearch struct {
	Fd int32
}
type BtrfsTransHandle struct {
	Transid           uint64
	AllocExcludeStart uint64
	AllocExcludeNr    uint64
	BlocksReserved    uint64
	BlocksUsed        uint64
	BlockGroup        *BtrfsBlockGroupCache
}
type UlistIterator struct {
	List **list.List
}
type UlistNode struct {
	Val  uint64
	Aux  uint64
	List *list.List
	Node RbNode
}
type Ulist struct {
	Nnodes uint64
	Nodes  *list.List
	Root   RbRoot
}
type BtrfsDevice struct {
	DevList        *list.List
	DevRoot        *BtrfsRoot
	FsDevices      *BtrfsFsDevices
	TotalIos       uint64
	Fd             int32
	Writeable      int32
	Name           *int8
	Label          *int8
	TotalDevs      uint64
	SuperBytesUsed uint64
	Devid          uint64
	TotalBytes     uint64
	BytesUsed      uint64
	IoAlign        uint64
	IoWidth        uint64
	SectorSize     uint64
	Type           uint64
	Uuid           [16]uint8
}
type BtrfsFsDevices struct {
	Fsid        [16]uint8
	LatestDevid uint64
	LatestTrans uint64
	LowestDevid uint64
	LatestBdev  int32
	LowestBdev  int32
	Devices     *list.List
	List        *list.Element
	Seeding     int32
	Seed        *BtrfsFsDevices
}
type BtrfsBioStripe struct {
	Dev      *BtrfsDevice
	Physical uint64
}
type BtrfsMultiBio struct {
	Error      int32
	NumStripes int32
	Stripes    [0]BtrfsBioStripe
}
type MapLookup struct {
	CacheExtent
	Type       uint64
	IoAlign    int32
	IoWidth    int32
	StripeLen  int32
	SectorSize int32
	NumStripes int32
	SubStripes int32
	Stripes    []BtrfsBioStripe
}

func (x *MapLookup) Less(y llrb.Item) bool {
	return (x.Start + x.Size) <= y.(*MapLookup).Start
}

type BtrfsExtentOps struct {
	AllocExtent *[0]byte
	FreeExtent  *[0]byte
}
type DeviceScan struct {
	Rc  *RecoverControl
	Dev *BtrfsDevice
	Fd  int
}
type ExtentRecord struct {
	CacheExtent
	Generation uint64
	Csum       [BTRFS_CSUM_SIZE]uint8
	Devices    [BTRFS_MAX_MIRRORS](*BtrfsDevice)
	Offsets    [BTRFS_MAX_MIRRORS]uint64
	Nmirrors   uint32
}

// NewExtentRecord creates a new extenet record from the block header and the Leaf size in recover control
func NewExtentRecord(header *BtrfsHeader, rc *RecoverControl) *ExtentRecord {
	return &ExtentRecord{
		CacheExtent: CacheExtent{Start: header.Bytenr, Size: uint64(rc.Leafsize)},
		Generation:  header.Generation,
		Csum:        header.Csum,
	}
}

func (x *ExtentRecord) Less(y llrb.Item) bool {
	return (x.Start + x.Size) <= y.(*ExtentRecord).Start
}

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

func NewRecoverControl(verbose bool, yes bool) *RecoverControl {

	rc := &RecoverControl{
		Chunk:   llrb.New(),
		EbCache: llrb.New(),
		Bg: BlockGroupTree{Tree: llrb.New(),
			Block_Groups: list.New(),
		},
		Devext: DeviceExtentTree{Tree: llrb.New(),
			ChunkOrphans:  list.New(),
			DeviceOrphans: list.New(),
		},
		Verbose:          verbose,
		Yes:              yes,
		GoodChunks:       list.New(),
		BadChunks:        list.New(),
		UnrepairedChunks: list.New(),
		FsDevices: &BtrfsFsDevices{
			Devices: list.New(),
		},
	}
	rc.FsDevices.Devices.PushBack(&BtrfsDevice{
		Devid: 1,
	},
	)
	return rc
	//	pthreadMutexInit(&rc->rcLock, NULL);
}
