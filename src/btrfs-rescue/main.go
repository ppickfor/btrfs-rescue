package main

import (
	. "btrfs"
	"bufio"
	"bytes"
	"code.google.com/p/go.net/context"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/petar/GoLLRB/llrb"
	"hash/crc32"
	"log"
	"os"
	"strconv"
)

const (
	APP_VERSION = "0.1"

	DEFAULT_BLOCKS = 1024
	DEFAULT_DEVICE = "/dev/sda1"
)

// The flag package provides a default help printer via -h switch
var (
	versionFlag     *bool   = flag.Bool("v", false, "Print the version number.")
	deviceFlag      *string = flag.String("d", DEFAULT_DEVICE, "The device to scan")
	blocksFlag      *int64  = flag.Int64("n", DEFAULT_BLOCKS, "The number of BTRFS_SIZE blocks to read")
	startblocksFlag *int64  = flag.Int64("s", 0, "The number of BTRFS_SIZE blocks to start at")
	bytenrFlag      *string = flag.String("b", "", "The previously scanned block bytenrs")
)

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

// reads super block into sb at offset sbBytenr
// if superRecover is != 0 then read superblock backups and find latest generation
func main() {
	// ctx is the Context for this handler. Calling cancel closes the
	// ctx.Done channel, which is the cancellation signal for requests
	// started by this handler.
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel() // Cancel ctx as soon as
	flag.Parse()   // Scan the arguments list
	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}
	rc := NewRecoverControl(true, false)
	RecoverPrepare(rc, *deviceFlag)
	bytenrChan := make(chan uint64, 40)
	csumBlockChan := make(chan treeBlock, 40)
	headerBlockchan := make(chan treeBlock, 40)
	itemBlockChan := make(chan itemBlock, 40)
	//	var items *[]BtrfsItem
	go byteConsumer(ctx, cancel, bytenrChan, csumBlockChan, rc)
	go csumByteblock(ctx, cancel, csumBlockChan, headerBlockchan)
	go headerConsumer(ctx, cancel, headerBlockchan, itemBlockChan, rc)
	go processItems(ctx, cancel, itemBlockChan, rc)
	//	treeBlock := new(treeBlock)
	//			fmt.Printf("rc: %+v\n", rc)
	size := uint64(rc.Leafsize)
	//	fsid := rc.Fsid[:]
	if len(*bytenrFlag) == 0 {
		// read treeblock from disk and pass to channel for processing
		start := uint64(*startblocksFlag) * size
		end := uint64(*blocksFlag) * size
	countFromParams:
		for bytenr := start; bytenr < end; bytenr += size {
			select {
			case <-ctx.Done():
				fmt.Printf("Done countFromParams\n")
				break countFromParams
			case bytenrChan <- bytenr:
			}
		}
	} else {
		file, err := os.Open(*bytenrFlag)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
	readFromFile:
		for i := 0; scanner.Scan(); i++ {
			bytenr, err := strconv.ParseUint(scanner.Text(), 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			select {
			case <-ctx.Done():
				fmt.Printf("Done readFromFile\n")
				break readFromFile
			case bytenrChan <- bytenr:
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("\nAll %d Extent buffers\n", rc.EbCache.Len())
	rc.EbCache.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
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

// byteConsumer reads bytenr from bytenrschan reads the treeblock and sends it to csumBlockChan
func byteConsumer(ctx context.Context, cancel context.CancelFunc, bytenrsChan <-chan uint64, csumBlockChan chan<- treeBlock, rc *RecoverControl) {

	treeBlock := new(treeBlock)
	size := uint64(rc.Leafsize)
	fsid := rc.Fsid[:]
	off0 := BtrfsSbOffset(0)
	off1 := BtrfsSbOffset(1)
	off2 := BtrfsSbOffset(2)
	treeBlock.byteblock = make([]byte, size)
	byteblock := make([]byte, size)
	defer close(csumBlockChan)
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Done byteConsumer\n")
			return
		case bytenr, ok := <-bytenrsChan:
			//			fmt.Printf("got byte %08x, staus %v\r", bytenr, ok)
			if ok {
				switch bytenr {
				case off0, off1, off2:
					continue loop
				default:
					//					fmt.Printf("size: %d, byteblock: %v\n", size, byteblock)
					ok, err := BtrfsReadTreeblock(rc.Fd, bytenr, size, fsid, &byteblock)
					if err != nil {
						fmt.Println(err)
						fmt.Printf("\n\n cancel byteConsumer %v\n\n", err)
						cancel()
						break loop
					}
					if ok {
						treeBlock.bytenr = bytenr
						copy(treeBlock.byteblock, byteblock)
						csumBlockChan <- *treeBlock
					}
				}
			} else {
				fmt.Printf("\n\n cancel byteConsumer %v\n\n", ok)
				cancel()
				return
			}
		}
	}
	//	for bytenr := range bytenrsChan {

}

// csumByteblock check crc of the treeBlocks byte buffer implemented as a channel filter
func csumByteblock(ctx context.Context, cancel context.CancelFunc, in <-chan (treeBlock), out chan<- (treeBlock)) {
	//	var inc,outc,ins,outs int
	csum := uint32(0)
	defer close(out)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Done csumByteblock\n")
			return
		case treeBlock, ok := <-in:
			if ok {
				byteblock := treeBlock.byteblock
				bytebr := bytes.NewReader(byteblock)
				binary.Read(bytebr, binary.LittleEndian, &csum)
				crc := crc32.Checksum((byteblock)[BTRFS_CSUM_SIZE:], Crc32c)
				if crc != csum {
					bytenr := treeBlock.bytenr
					fmt.Printf("crc32c mismatch @%08x have %08x expected %08x\n", bytenr, csum, crc)
				} else {
					out <- treeBlock
				}
				//		fmt.Printf("read treeblock @%d, %v\n",bytenr,(*byteblock)[0:4])
				//		inc++
				//		outc++
				//		ins= ins+len(in)
				//		outs=outs+len(out)
				//		fmt.Printf("csumByteblock Chan len in: %03.2f out: %03.2f\r",float64(ins)/float64(inc),float64(outs)/float64(outc))
			} else {
				fmt.Printf("\n\n cancel csumByteblock, %v\n\n", ok)
				cancel()
				return
			}
		}
	}
}

// headerConsumer reads blocks from headerBlockchan and processes them via detailBlock
func headerConsumer(ctx context.Context, cancel context.CancelFunc, headerBlockchan <-chan (treeBlock), itemBlockChan chan itemBlock, rc *RecoverControl) {

	// process treeblock from goroutine via channel
	defer close(itemBlockChan)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Done headerConsumer\n")
			return
		case treeBlock, ok := <-headerBlockchan:
			if ok {
				//			treeBlock = treeBlock
				//		fmt.Printf("treeBlock: %+v\n", treeBlock)
				//		fmt.Printf("from chan treeblock: @%d, %v\n", treeBlock.bytent, treeBlock.byteblock[0:4])
				detailBlock(&treeBlock, itemBlockChan, rc)
			} else {
				fmt.Printf("\n\n cancel headerConsumer %v\n\n", ok)
				cancel()
				return
			}
		}
	}
}

// detailBlock build cache structures for selected block records and items at physical bytenr
func detailBlock(treeBlock *treeBlock, itemBlockChan chan<- (itemBlock), rc *RecoverControl) {

	byteblock := treeBlock.byteblock
	header := new(BtrfsHeader)
	bytereader := bytes.NewReader(byteblock)
	_ = binary.Read(bytereader, binary.LittleEndian, header)
	bytenr := treeBlock.bytenr
	tree := rc.EbCache
	er := NewExtentRecord(header, rc)
again:
	if tree.Has(er) {
		exists := tree.Get(er).(*ExtentRecord)
		if exists.Generation > er.Generation {
			//				fmt.Printf("Exists:%d\r", tree.Len())
			return
		}
		if exists.Generation == er.Generation {
			if exists.CacheExtent.Start != er.CacheExtent.Start ||
				exists.CacheExtent.Size != er.CacheExtent.Size ||
				bytes.Compare(exists.Csum[:], er.Csum[:]) != 0 {
				//							exists but different
				fmt.Printf("detailBlock: Exists but dif %+v\n", er)
				return
			} else {
				//					fmt.Printf("Mirror:%d\r", tree.Len())
				//				fmt.Printf("mirror %+v\n", exists.(*ExtentRecord))
				if exists.Nmirrors < BTRFS_SUPER_MIRROR_MAX {
					exists.Devices[exists.Nmirrors] = nil
					exists.Offsets[exists.Nmirrors] = bytenr
				}
				exists.Nmirrors++
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
	itemBlockChan <- itemBlock
}

// processItems reads items from itemBlockChan and processes the leaves
func processItems(ctx context.Context, cancel context.CancelFunc, itemBlockChan chan (itemBlock), rc *RecoverControl) {

	// process items in new treeblock
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Done processItems\n")
			return
		case itemBlock, ok := <-itemBlockChan:
			if ok {
				level := itemBlock.Level
				nritems := itemBlock.Nritems
				owner := itemBlock.Owner
				generation := itemBlock.Generation
				infoByteBlock := itemBlock.InfoByteBlock
				bytereader := bytes.NewReader(infoByteBlock)
				if level == 0 {
					// Leaf
					items := make([]BtrfsItem, nritems)

					_ = binary.Read(bytereader, binary.LittleEndian, items)
					//			fmt.Printf("Leaf @%08x: Items: %d\r", bytenr, leaf.Header.Nritems)
					switch owner {
					case BTRFS_EXTENT_TREE_OBJECTID, BTRFS_DEV_TREE_OBJECTID:
						/* different tree use different generation */
						//			if header.Generation <= rc.Generation {
						extractMetadataRecord(rc, generation, items, infoByteBlock)
						//			}
					case BTRFS_CHUNK_TREE_OBJECTID:
						//			if header.Generation <= rc.ChunkRootGeneration {
						extractMetadataRecord(rc, generation, items, infoByteBlock)
						//			}
					}
				}
			} else {
				fmt.Printf("\n\n cancel processItems %+v\n\n", ok)
				cancel()
				return
			}

		}
	}
}

// extractMetadataRecord iterates items of the current block and proccess each blockgroup, chunk and dev extent adding them to caches
func extractMetadataRecord(rc *RecoverControl, generation uint64, items []BtrfsItem, itemBuf []byte) {

	for _, item := range items {
		switch item.Key.Type {
		case BTRFS_BLOCK_GROUP_ITEM_KEY:
			//			pthreadMutexLock(&rc->rcLock);
			//			fmt.Printf("BLOCK Group Item: %+v\n", &leaf.Items[i])
			processBlockGroupItem(&rc.Bg, generation, &item, itemBuf)
			//			pthreadMutexUnlock(&rc->rcLock);
			//			break;
		case BTRFS_CHUNK_ITEM_KEY:
			//			pthreadMutexLock(&rc->rcLock);
			//			ret = processChunkItem(&rc->chunk, leaf, &key, i);
			//			pthreadMutexUnlock(&rc->rcLock);
			//			break;
		case BTRFS_DEV_EXTENT_KEY:
			//			pthreadMutexLock(&rc->rcLock);
			processDeviceExtentItem(&rc.Devext, generation, &item, itemBuf)
			//			pthreadMutexUnlock(&rc->rcLock);
			//			break;
		}
	}
}

// processBlockGroupItem creates a new block group record and update the cache by latest generation for each blockgroup item
func processBlockGroupItem(bgCache *BlockGroupTree, generation uint64, item *BtrfsItem, itemBuf []byte) {

	rec := NewBlockGroupRecord(generation, item, itemBuf)
again:
	if bgCache.Tree.Has(rec) {
		exists := bgCache.Tree.Get(rec).(*BlockGroupRecord)
		/*check the generation and replace if needed*/
		if exists.Generation > rec.Generation {
			return
		}
		if exists.Generation == rec.Generation {
			//			int offset = offsetof(struct blockGroupRecord,
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
			//			goto freeOut;
			fmt.Printf("processBlockGroupItem: same generation %+v\n", rec)
			return
		}
		bgCache.Block_Groups.Remove(exists.List)
		exists.List = nil
		//		listDelInit(&exists.List)
		bgCache.Tree.Delete(exists)
		/*
		 * We must do seach again to avoid the following cache.
		 * /--old bg 1--//--old bg 2--/
		 *        /--new bg--/
		 */
		goto again
	}
	bgCache.Tree.InsertNoReplace(rec)
	rec.List = bgCache.Block_Groups.PushBack(rec)
}

// processDeviceExtentItem creates a new device extent record and update the cache by latest generation
func processDeviceExtentItem(devextCache *DeviceExtentTree, generation uint64, item *BtrfsItem, itemBuf []byte) {

	rec := NewDeviceExtentRecord(generation, item, itemBuf)
again:
	if devextCache.Tree.Has(rec) {
		exists := devextCache.Tree.Get(rec).(*DeviceExtentRecord)
		// check the generation and replace if needed
		if exists.Generation > rec.Generation {
			return
		}
		if exists.Generation == rec.Generation {
			// should not happen
			fmt.Printf("processDeviceExtentItem: same generation %+v\n", rec)
			return
		}
		devextCache.ChunkOrphans.Remove(exists.ChunkList)
		devextCache.DeviceOrphans.Remove(exists.DeviceList)
		exists.ChunkList = nil
		exists.DeviceList = nil
		//		listDelInit(&exists.ChunkList)
		//		listDelInit(&exists.DeviceList)
		devextCache.Tree.Delete(exists)
		/*
		 * We must do seach again to avoid the following cache.
		 * /--old bg 1--//--old bg 2--/
		 *        /--new bg--/
		 */
		goto again
	}
	devextCache.Tree.InsertNoReplace(rec)
	rec.ChunkList = devextCache.ChunkOrphans.PushBack(rec.ChunkList)
	rec.DeviceList = devextCache.DeviceOrphans.PushBack(rec.DeviceList)
}
