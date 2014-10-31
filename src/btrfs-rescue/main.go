package main

import (
	"btrfs"
	"bytes"
	"flag"
	"fmt"
	//	"log"
	"os"
	//	"syscall"
	"encoding/binary"
	"encoding/hex"
	//	"strings"
	"bufio"
	"github.com/petar/GoLLRB/llrb"
	"hash/crc32"
	"log"
	"regexp"
	"strconv"
)

const (
	APP_VERSION = "0.1"

	DEFAULT_BLOCKS = 1024
	DEFAULT_DEVICE = "/dev/sda"
)

// The flag package provides a default help printer via -h switch
var (
	versionFlag     *bool   = flag.Bool("v", false, "Print the version number.")
	deviceFlag      *string = flag.String("d", DEFAULT_DEVICE, "The device to scan")
	blocksFlag      *int64  = flag.Int64("n", DEFAULT_BLOCKS, "The number of BTRFS_SIZE blocks to read")
	startblocksFlag *int64  = flag.Int64("s", 0, "The number of BTRFS_SIZE blocks to start at")
	fsidFlag        *string = flag.String("f", "671e120db12f497da4c42e41b3d30fec", "The Hex fsid of the file system to match")
	bytenrFlag      *string = flag.String("b", "", "The previously scanned block bytenrs")
)

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

type Cache_Extent_record struct {
	Objectid uint64
	Start    uint64
	Size     uint64
}

type Extent_record struct {
	Cache      Cache_Extent_record
	Generation uint64
	Csum       [btrfs.BTRFS_CSUM_SIZE]uint8
	devices    [btrfs.BTRFS_MAX_MIRRORS](*btrfs.Btrfs_device)
	offsets    [btrfs.BTRFS_MAX_MIRRORS]uint64
	nmirrors   uint32
}

// create a new extenet record from the block header and the Leaf size in recover control
func NewExtent_record(header *btrfs.Btrfs_header, rc *btrfs.Recover_control) *Extent_record {
	return &Extent_record{
		Cache:      Cache_Extent_record{Start: header.Bytenr, Size: uint64(rc.Leafsize)},
		Generation: header.Generation,
		Csum:       header.Csum,
	}
}
func (x *Extent_record) Less(y llrb.Item) bool {
	return (x.Cache.Start + x.Cache.Size) <= y.(*Extent_record).Cache.Start
}

type Block_group_record struct {
	Cache      Cache_Extent_record
	List       btrfs.List_head
	Generation uint64
	Objectid   uint64
	Type       uint8
	Offset     uint64
	Flags      uint64
}

// create a new block group record from the a block heade, a block group ietem key and the items byte buffer
func NewBlock_group_record(generation uint64, item *btrfs.Btrfs_item, itemsBuf []byte) *Block_group_record {

	key := item.Key
	itemPtr := itemsBuf[item.Offset:]
	bytereader := bytes.NewReader(itemPtr)
	bgItem := new(btrfs.Btrfs_block_group_item)
	_ = binary.Read(bytereader, binary.LittleEndian, bgItem)
	//	fmt.Printf("bgItem: %+v\n", bgItem)
	this := &Block_group_record{
		Cache:      Cache_Extent_record{Start: key.Objectid, Size: key.Offset},
		Generation: generation,
		Objectid:   key.Objectid,
		Type:       key.Type,
		Offset:     key.Offset,
		Flags:      bgItem.Flags,
		//		List:	btrfs.List_head{Next: &List, Prev: &List},
	}
	this.List = btrfs.List_head{Next: &this.List, Prev: &this.List}
	return this
}
func (x *Block_group_record) Less(y llrb.Item) bool {
	return (x.Cache.Start + x.Cache.Size) <= y.(*Block_group_record).Cache.Start
}

type Device_extent_record struct {
	Cache       Cache_Extent_record
	Chunk_list  btrfs.List_head
	Device_list btrfs.List_head
	Generation  uint64
	Objectid    uint64
	Type        uint8

	Offset         uint64
	Chunk_objectid uint64
	Chunk_offset   uint64
	Length         uint64
}

// create a new device extent record from the a block header, a device extent item key and the items byte buffer
func NewDevice_extent_record(generation uint64, item *btrfs.Btrfs_item, itemsBuf []byte) *Device_extent_record {

	key := item.Key
	itemPtr := itemsBuf[item.Offset:]
	bytereader := bytes.NewReader(itemPtr)
	deItem := new(btrfs.Btrfs_dev_extent)
	_ = binary.Read(bytereader, binary.LittleEndian, deItem)
	//	fmt.Printf("bgItem: %+v\n", bgItem)
	this := &Device_extent_record{
		Cache:          Cache_Extent_record{Objectid: key.Objectid, Start: key.Offset, Size: deItem.Length},
		Generation:     generation,
		Objectid:       key.Objectid,
		Type:           key.Type,
		Offset:         key.Offset,
		Chunk_objectid: deItem.Chunk_Objectid,
		Chunk_offset:   deItem.Chunk_Offset,
		Length:         deItem.Length,
	}
	this.Chunk_list = btrfs.List_head{Next: &this.Chunk_list, Prev: &this.Chunk_list}
	this.Device_list = btrfs.List_head{Next: &this.Device_list, Prev: &this.Device_list}
	return this
}

func (x *Device_extent_record) Less(y llrb.Item) bool {
	switch {
	case x.Cache.Objectid < y.(*Device_extent_record).Cache.Objectid:
		return true
	case x.Cache.Objectid > y.(*Device_extent_record).Cache.Objectid:
		return false
	}
	return (x.Cache.Start + x.Cache.Size) <= y.(*Device_extent_record).Cache.Start
}

type treeBlock struct {
	bytenr    uint64
	byteblock []byte
}

type itemBlock struct {
	Owner         uint64
	Nritems       uint32
	Generation    uint64
	InfoByteBlock []byte
	Level         uint8
}

func list_add_tail(xnew, head *btrfs.List_head) {
	__list_add(xnew, head.Prev, head)
}
func __list_add(xnew,
	prev,
	next *btrfs.List_head) {
	next.Prev = xnew
	xnew.Next = next
	xnew.Prev = prev
	prev.Next = xnew
}
func __list_del(prev *btrfs.List_head, next *btrfs.List_head) {
	next.Prev = prev
	prev.Next = next
}
func list_del_init(entry *btrfs.List_head) {

	__list_del(entry.Prev, entry.Next)
	entry.Next = entry
	entry.Prev = entry
}

// reads super block into sb at offset sb_bytenr
// if super_recover is != 0 then read superblock backups and find latest generation
func main() {

	var rc btrfs.Recover_control

	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}
	btrfs.Init_recover_control(&rc, true, false)
	btrfs.Recover_prepare(&rc, *deviceFlag)
	//	header := new(btrfs.Btrfs_header)
	nonhex := regexp.MustCompile("[^0-9a-fA-F]*")
	fsidstr := nonhex.ReplaceAllString(*fsidFlag, "")
	fsid, err := hex.DecodeString(fsidstr)
	bytenrs := make([]uint64, 5874688)
	off0 := btrfs.Btrfs_sb_offset(0)
	off1 := btrfs.Btrfs_sb_offset(1)
	off2 := btrfs.Btrfs_sb_offset(2)

	csumBlockchan := make(chan treeBlock, 40)
	headerBlockchan := make(chan treeBlock, 40)
	itemBlockchan := make(chan itemBlock, 40)

	if err != nil {
		fmt.Errorf("fsidFlag:%v is not hex, %v\n", *fsidFlag, err)
		os.Exit(1)
	}
	//	var items *[]btrfs.Btrfs_item
	go processItems(&rc, itemBlockchan)
	go csumByteblock(csumBlockchan, headerBlockchan)
	if len(*bytenrFlag) == 0 {
		// read treeblock from disk and pass to channel for processing
		go func() {
			treeBlock := new(treeBlock)
			//			fmt.Printf("rc: %+v\n", rc)
			size := uint64(rc.Leafsize)
			fsid := rc.Fsid[:]
		loop:
			for bytenr := uint64(*startblocksFlag) * size; bytenr < uint64(*blocksFlag)*size; bytenr += size {
				switch bytenr {
				case off0, off1, off2:
					continue loop
				default:
					byteblock := make([]byte, size)
					ok, err := btrfs.Btrfs_read_treeblock(rc.Fd, bytenr, size, fsid, &byteblock)
					if err != nil {
						fmt.Println(err)
						break loop
					}
					if ok {
						treeBlock.bytenr = bytenr
						treeBlock.byteblock = byteblock
						csumBlockchan <- *treeBlock
					}
				}
			}
			close(csumBlockchan)
			close(itemBlockchan)
		}()

	} else {
		file, err := os.Open(*bytenrFlag)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for i := 0; scanner.Scan(); i++ {
			bytenr, err := strconv.ParseUint(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			bytenrs[i] = bytenr
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
		// read treeblock from disk and pass to channel for processing
		go func() {
			treeBlock := new(treeBlock)
			size := uint64(rc.Leafsize)
			fsid := rc.Fsid[:]
		loop:
			for _, bytenr := range bytenrs {
				switch bytenr {
				case off0, off1, off2:
					continue loop
				default:
					byteblock := make([]byte, size)
					//					fmt.Printf("size: %d, byteblock: %v\n", size, byteblock)
					ok, err := btrfs.Btrfs_read_treeblock(rc.Fd, bytenr, size, fsid, &byteblock)
					if err != nil {
						fmt.Println(err)
						break loop
					}
					if ok {
						treeBlock.bytenr = bytenr
						treeBlock.byteblock = byteblock
						csumBlockchan <- *treeBlock
					}
				}
			}
			close(csumBlockchan)
			close(itemBlockchan)
		}()

	}
	// process treeblock from goroutine via channel
	for treeBlock := range headerBlockchan {
		//			treeBlock = treeBlock
		//		fmt.Printf("treeBlock: %+v\n", treeBlock)
		//		fmt.Printf("from chan treeblock: @%d, %v\n", treeBlock.bytent, treeBlock.byteblock[0:4])
		detailBlock(&treeBlock, itemBlockchan, &rc, fsid)
	}
	fmt.Printf("\nAll %d Extent buffers\n", rc.Eb_cache.Len())
	rc.Eb_cache.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		fmt.Printf("%+v\n", i)
		return true
	})
	fmt.Printf("\nAll %d Block Groups\n", rc.Bg.Tree.Len())
	rc.Bg.Tree.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		fmt.Printf("%+v\n", i)
		return true
	})
	fmt.Printf("\nAll %d Device Extents\n", rc.Devext.Tree.Len())
	rc.Devext.Tree.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		fmt.Printf("%+v\n", i)
		return true
	})
}

func csumByteblock(in, out chan (treeBlock)) {
	for treeBlock := range in {
		byteblock := treeBlock.byteblock
		csum := uint32(0)
		bytebr := bytes.NewReader(byteblock)
		binary.Read(bytebr, binary.LittleEndian, &csum)
		crc := crc32.Checksum((byteblock)[btrfs.BTRFS_CSUM_SIZE:], btrfs.Crc32c)
		if crc != csum {
			bytenr := treeBlock.bytenr
			fmt.Printf("crc32c mismatch @%08x have %08x expected %08x\n", bytenr, csum, crc)
		} else {
			out <- treeBlock
		}
		//		fmt.Printf("read treeblock @%d, %v\n",bytenr,(*byteblock)[0:4])
	}
	close(out)
}

// build cache structures for selected block records and items at physical bytenr
func detailBlock(treeBlock *treeBlock, itemBlockchan chan (itemBlock), rc *btrfs.Recover_control, fsid []byte) {

	byteblock := treeBlock.byteblock
	header := new(btrfs.Btrfs_header)
	bytereader := bytes.NewReader(byteblock)
	_ = binary.Read(bytereader, binary.LittleEndian, header)
	bytenr := treeBlock.bytenr

	tree := rc.Eb_cache

	er := NewExtent_record(header, rc)

again:
	if exists := tree.Get(er); exists != nil {
		if exists.(*Extent_record).Generation > er.Generation {
			//				fmt.Printf("Exists:%d\r", tree.Len())
			return
		}
		if exists.(*Extent_record).Generation == er.Generation {
			if exists.(*Extent_record).Cache.Start != er.Cache.Start ||
				exists.(*Extent_record).Cache.Size != er.Cache.Size ||
				bytes.Compare(exists.(*Extent_record).Csum[:], er.Csum[:]) != 0 {
				//							exists but different
				fmt.Printf("detailBlock: Exists but dif %+v\n", er)
				return
			} else {
				//					fmt.Printf("Mirror:%d\r", tree.Len())
				//				fmt.Printf("mirror %+v\n", exists.(*Extent_record))
				if exists.(*Extent_record).nmirrors < 3 {
					exists.(*Extent_record).devices[exists.(*Extent_record).nmirrors] = nil
					exists.(*Extent_record).offsets[exists.(*Extent_record).nmirrors] = bytenr
				}
				exists.(*Extent_record).nmirrors++
				tree.ReplaceOrInsert(exists)
				return
			}
		}
		tree.Delete(er)
		//			fmt.Printf("Worse:%d\n", tree.Len())
		goto again
	}
	tree.InsertNoReplace(er)
	itemBlock := itemBlock{Generation: header.Generation,
		Level:         header.Level,
		Nritems:       header.Nritems,
		Owner:         header.Owner,
		InfoByteBlock: byteblock[len(byteblock)-bytereader.Len():],
	}
	itemBlockchan <- itemBlock

}

func processItems(rc *btrfs.Recover_control, itemBlockchan chan (itemBlock)) {
	// process items in new treeblock
	for itemBlock := range itemBlockchan {
		level := itemBlock.Level
		nritems := itemBlock.Nritems
		owner := itemBlock.Owner
		generation := itemBlock.Generation
		infoByteBlock := itemBlock.InfoByteBlock
		bytereader := bytes.NewReader(infoByteBlock)
		if level == 0 {
			// Leaf
			items := make([]btrfs.Btrfs_item, nritems)

			_ = binary.Read(bytereader, binary.LittleEndian, items)
			//			fmt.Printf("Leaf @%08x: Items: %d\r", bytenr, leaf.Header.Nritems)
			switch owner {
			case btrfs.BTRFS_EXTENT_TREE_OBJECTID, btrfs.BTRFS_DEV_TREE_OBJECTID:
				/* different tree use different generation */
				//			if header.Generation <= rc.Generation {
				extract_metadata_record(rc, generation, items, infoByteBlock)
				//			}
			case btrfs.BTRFS_CHUNK_TREE_OBJECTID:
				//			if header.Generation <= rc.Chunk_root_generation {
				extract_metadata_record(rc, generation, items, infoByteBlock)
				//			}
			}

		} else {
			// Node
			//			node := new(btrfs.Btrfs_node)
			//			node.Header = *header
			//			node.Ptrs = make([]btrfs.Btrfs_key_ptr, header.Nritems)
			//			_ = binary.Read(bytereader, binary.LittleEndian, node.Ptrs)
			//			fmt.Printf("Node @%08x: Pointers: %d\r", bytenr, node.Header.Nritems)
		}

	}
}

// iterate items of the current block and proccess each blockgroup, chunk and dev extent adding them to caches
func extract_metadata_record(rc *btrfs.Recover_control, generation uint64, items []btrfs.Btrfs_item, itemBuf []byte) {

	for _, item := range items {
		//		fmt.Printf("Leaf Item: %+v\n", &leaf.Items[i])

		switch item.Key.Type {
		case btrfs.BTRFS_BLOCK_GROUP_ITEM_KEY:
			//			pthread_mutex_lock(&rc->rc_lock);
			//			fmt.Printf("BLOCK Group Item: %+v\n", &leaf.Items[i])
			process_block_group_item(&rc.Bg, generation, &item, itemBuf)
			//			pthread_mutex_unlock(&rc->rc_lock);
			//			break;
		case btrfs.BTRFS_CHUNK_ITEM_KEY:
			//			pthread_mutex_lock(&rc->rc_lock);
			//			ret = process_chunk_item(&rc->chunk, leaf, &key, i);
			//			pthread_mutex_unlock(&rc->rc_lock);
			//			break;
		case btrfs.BTRFS_DEV_EXTENT_KEY:
			//			pthread_mutex_lock(&rc->rc_lock);
			process_device_extent_item(&rc.Devext, generation, &item, itemBuf)
			//			pthread_mutex_unlock(&rc->rc_lock);
			//			break;
		}
		//		if (ret)
		//			break;
	}

}

// create a new block group record and update the cache by latest generation for each blockgroup item
func process_block_group_item(
	bg_cache *btrfs.Block_group_tree,
	generation uint64,
	item *btrfs.Btrfs_item,
	itemBuf []byte) {

	rec := NewBlock_group_record(generation, item, itemBuf)
again:
	if exists := bg_cache.Tree.Get(rec); exists != nil {
		/*check the generation and replace if needed*/
		if exists.(*Block_group_record).Generation > rec.Generation {
			return
		}
		//		if (exist->generation > rec->generation)
		//			goto free_out;

		if exists.(*Block_group_record).Generation == rec.Generation {
			//			int offset = offsetof(struct block_group_record,
			//					      generation);
			/*
			 * According to the current kernel code, the following
			 * case is impossble, or there is something wrong in
			 * the kernel code.
			 */
			//			if (memcmp(((void *)exist) + offset,
			//				   ((void *)rec) + offset,
			//				   sizeof(*rec) - offset))
			//				ret = -EEXIST;
			//			goto free_out;
			fmt.Printf("process_block_group_item: same generation %+v\n", rec)
			return

		}
		list_del_init(&exists.(*Block_group_record).List)
		bg_cache.Tree.Delete(exists)
		/*
		 * We must do seach again to avoid the following cache.
		 * /--old bg 1--//--old bg 2--/
		 *        /--new bg--/
		 */
		goto again

	}
	bg_cache.Tree.InsertNoReplace(rec)
	list_add_tail(&rec.List, &bg_cache.Block_Groups)
	//	bg_cache.InsertNoReplace(rec)
}

// create a new device extent record and update the cache by latest generation
func process_device_extent_item(devext_cache *btrfs.Device_extent_tree,
	generation uint64,
	item *btrfs.Btrfs_item,
	itemBuf []byte) {
	rec := NewDevice_extent_record(generation, item, itemBuf)
again:
	if exists := devext_cache.Tree.Get(rec); exists != nil {
		/*check the generation and replace if needed*/
		if exists.(*Device_extent_record).Generation > rec.Generation {
			return
		}
		//		if (exist->generation > rec->generation)
		//			goto free_out;

		if exists.(*Device_extent_record).Generation == rec.Generation {
			// should not happen
			fmt.Printf("process_device_extent_item: same generation %+v\n", rec)
			return
		}
		list_del_init(&exists.(*Device_extent_record).Chunk_list)
		list_del_init(&exists.(*Device_extent_record).Device_list)
		devext_cache.Tree.Delete(exists)
		/*
		 * We must do seach again to avoid the following cache.
		 * /--old bg 1--//--old bg 2--/
		 *        /--new bg--/
		 */
		goto again

	}
	devext_cache.Tree.ReplaceOrInsert(rec)
	//	devext_cache.InsertNoReplace(rec)
	list_add_tail(&rec.Chunk_list, &devext_cache.Chunk_orphans)
	list_add_tail(&rec.Device_list, &devext_cache.Device_orphans)
}
