package main

import (
	. "btrfs"
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"strconv"
	"sync"

	"code.google.com/p/go.net/context"
	"github.com/petar/GoLLRB/llrb"
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
	writeBytenrFlag *string = flag.String("w", "", "The new scanned block bytenrs")
	wg              sync.WaitGroup
	wfile           *os.File
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
	// create a new file to get list of valid bytnrs for this run
	if len(*writeBytenrFlag) != 0 {
		var err error
		fmt.Printf("\nWriting new bytenr file to %v\n", *writeBytenrFlag)
		wfile, err = os.Create(*writeBytenrFlag)
		if err != nil {
			log.Fatal(err)
		}
		defer wfile.Close()
	}
	rc := NewRecoverControl(true, false)
	RecoverPrepare(rc, *deviceFlag)
	bytenrChan := make(chan uint64, 20)
	csumBlockChan := make(chan treeBlock, 20)
	headerBlockchan := make(chan treeBlock, 20)
	itemBlockChan := make(chan itemBlock, 20)
	//	var items *[]BtrfsItem
	wg.Add(4)
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
		i := 0
	countFromParams:
		for bytenr := start; bytenr < end; bytenr += size {
			select {
			case <-ctx.Done():
				fmt.Printf("Done countFromParams\n")
				break countFromParams
			case bytenrChan <- bytenr:
				i++
			}
		}
		fmt.Printf("Read %d uint64s\n", i)
	} else {
		file, err := os.Open(*bytenrFlag)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		i := 0
	readFromFile:
		for ; scanner.Scan(); i++ {
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
		fmt.Printf("Read %d uint64s\n", i)
	}
	close(bytenrChan)
	wg.Wait()
	checkChunks(rc.Chunk, &rc.Bg, &rc.Devext, rc.GoodChunks, rc.BadChunks, false)

	btrfsRecoverChunks(rc)

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
	fmt.Printf("\nAll %d Chunks\n", rc.Chunk.Len())
	rc.Chunk.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
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

	size := uint64(rc.Leafsize)
	fsid := rc.Fsid[:]
	off0 := BtrfsSbOffset(0)
	off1 := BtrfsSbOffset(1)
	off2 := BtrfsSbOffset(2)
	byteblock := make([]byte, size)
	defer close(csumBlockChan)
	defer wg.Done()
	i := 0
loop:
	for {

		select {
		case <-ctx.Done():
			fmt.Printf("Done byteConsumer\n")
			break loop
		case bytenr, ok := <-bytenrsChan:
			//			fmt.Printf("got byte %08x, staus %v\r", bytenr, ok)
			if ok {
				switch bytenr {
				case off0, off1, off2:
					// superblock
					continue loop
				default:
					//					fmt.Printf("size: %d, byteblock: %v\n", size, byteblock)
					goodBlock, err := BtrfsReadTreeblock(rc.Fd, bytenr, size, fsid, &byteblock)
					if goodBlock {
						treeBlock := treeBlock{bytenr: bytenr, byteblock: make([]byte, size)}
						copy(treeBlock.byteblock, byteblock)
						csumBlockChan <- treeBlock
					} else {
						if err != nil {
							fmt.Println(err)
							fmt.Printf("byteConsumer BtrfsReadTreeblock failed %v\n", err)
							fmt.Printf("byteConsumer read %d uint64s\n", i)
							cancel()
							break loop
						} else {
							i++
						}
					}
				}
			} else {
				fmt.Printf("\n\n cancel byteConsumer %v\n\n", ok)
				fmt.Printf("byteConsumer %d blocks with bad fsid\n", i)
				//				cancel()
				break loop
			}
		}
	}
}

// csumByteblock check crc of the treeBlocks byte buffer implemented as a channel filter
func csumByteblock(ctx context.Context, cancel context.CancelFunc, in <-chan (treeBlock), out chan<- (treeBlock)) {
	//	var inc,outc,ins,outs int
	csum := uint32(0)
	defer close(out)
	defer wg.Done()
	i := 0
	var last uint64
	for {
		i++
		select {
		case <-ctx.Done():
			fmt.Printf("Done csumByteblock\n")
			return
		case treeBlock, ok := <-in:
			if ok {
				byteblock := treeBlock.byteblock
				bytenr := treeBlock.bytenr
				if bytenr > 0 {
					last = bytenr
				}
				bytebr := bytes.NewReader(byteblock)
				binary.Read(bytebr, binary.LittleEndian, &csum)
				crc := crc32.Checksum((byteblock)[BTRFS_CSUM_SIZE:], Crc32c)
				if crc != csum {
					bytenr := treeBlock.bytenr
					fmt.Printf("crc32c mismatch @%08x have %08x expected %08x\n", bytenr, csum, crc)
				} else {
					if wfile != nil {
						fmt.Fprintf(wfile, "%d\n", bytenr)
					}
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
				fmt.Printf("csumByteblock: read %d treeblocks\n", i)
				fmt.Printf("csumByteblock: last bytenr %d\n", last)
				//				cancel()
				return
			}
		}
	}
}

// headerConsumer reads blocks from headerBlockchan and processes them via detailBlock
func headerConsumer(ctx context.Context, cancel context.CancelFunc, headerBlockchan <-chan (treeBlock), itemBlockChan chan itemBlock, rc *RecoverControl) {

	// process treeblock from goroutine via channel
	defer close(itemBlockChan)
	defer wg.Done()
	i := 0
	for {
		i++
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
				fmt.Printf("headerConsumer read %d treeblocks\n", i)
				//				cancel()
				return
			}
		}
	}
}

// processItems reads items from itemBlockChan and processes the leaves
func processItems(ctx context.Context, cancel context.CancelFunc, itemBlockChan chan (itemBlock), rc *RecoverControl) {

	defer wg.Done()
	// process items in new treeblock
	i := 0
	for {
		i++
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
					//			fmt.Printf("Leaf @%08x: Items: %d\r", bytenr, leaf.Header.Nritems)
					switch owner {
					case BTRFS_EXTENT_TREE_OBJECTID, BTRFS_DEV_TREE_OBJECTID:
						/* different tree use different generation */
						if generation <= rc.Generation {
							items := make([]BtrfsItem, nritems)
							_ = binary.Read(bytereader, binary.LittleEndian, items)
							extractMetadataRecord(rc, generation, items, infoByteBlock)
						}
					case BTRFS_CHUNK_TREE_OBJECTID:
						if generation <= rc.ChunkRootGeneration {
							items := make([]BtrfsItem, nritems)
							_ = binary.Read(bytereader, binary.LittleEndian, items)
							extractMetadataRecord(rc, generation, items, infoByteBlock)
						}
					}
				}
			} else {
				fmt.Printf("\n\n cancel processItems %+v\n\n", ok)
				fmt.Printf("processItems read %d itemblocks\n", i)
				//				cancel()
				return
			}

		}
	}
}

// detailBlock build cache structures for selected block records and items at physical bytenr
func detailBlock(treeBlock *treeBlock, itemBlockChan chan<- (itemBlock), rc *RecoverControl) {

	byteblock := treeBlock.byteblock
	header := BtrfsHeader{}
	bytereader := bytes.NewReader(byteblock)
	_ = binary.Read(bytereader, binary.LittleEndian, &header)
	bytenr := treeBlock.bytenr
	tree := rc.EbCache
	er := NewExtentRecord(&header, rc)
again:
	if i := tree.Get(er); i != nil {
		exists := i.(*ExtentRecord)
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
					// TODO
					//					exists.Devices[exists.Nmirrors] = nil
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
	//	er.Devices[0] = nil
	er.Offsets[0] = bytenr
	er.Nmirrors++
	tree.InsertNoReplace(er)
	itemBlock := itemBlock{Generation: header.Generation,
		Level:         header.Level,
		Nritems:       header.Nritems,
		Owner:         header.Owner,
		InfoByteBlock: byteblock[len(byteblock)-bytereader.Len():],
	}
	itemBlockChan <- itemBlock
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
			processChunkItem(rc.Chunk, generation, &item, itemBuf)
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
	if i := bgCache.Tree.Get(rec); i != nil {
		exists := i.(*BlockGroupRecord)
		/*check the generation and replace if needed*/
		if exists.Generation > rec.Generation {
			return
		}
		if exists.Generation == rec.Generation {
			// TODO
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
	if i := devextCache.Tree.Get(rec); i != nil {
		exists := i.(*DeviceExtentRecord)
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
	rec.ChunkList = devextCache.ChunkOrphans.PushBack(rec)
	rec.DeviceList = devextCache.DeviceOrphans.PushBack(rec)
}

// processChunkItem: inserts new chunks into the chunk cache based on generation
func processChunkItem(chunkCache *llrb.LLRB, generation uint64, item *BtrfsItem, itemBuf []byte) {
	//	fmt.Printf("found chunk item\n")
	rec := NewChunkRecord(generation, item, itemBuf)
again:
	if i := chunkCache.Get(rec); i != nil {
		exists := i.(*ChunkRecord)
		if exists.Generation > rec.Generation {
			return
		}
		if exists.Generation == rec.Generation {
			numStripes := rec.NumStripes
			if exists.NumStripes != numStripes {
				fmt.Printf("processChunkItem: same generation but different details %+v\n", rec)
				// TODO
				// compare everything from generation to end of the record (varies by num stripes)
				//int num_stripes = rec->num_stripes;
				//			int rec_size = btrfs_chunk_record_size(num_stripes);
				//			int offset = offsetof(struct chunk_record, generation);
				//
				//			if (exist->num_stripes != rec->num_stripes ||
				//			    memcmp(((void *)exist) + offset,
				//				   ((void *)rec) + offset,
				//				   rec_size - offset))
				//				ret = -EEXIST;
				//			goto free_out;
				return
			}
		}
		chunkCache.Delete(exists)
		goto again
	}
	chunkCache.InsertNoReplace(rec)
}

// checkChunks: calculates the length of a single stripe based on the type of block group a total length and the number of stripes.
func calcStripeLength(typeFlags, length uint64, numStripes uint16) uint64 {
	var stripeSize uint64
	switch {
	case typeFlags&BTRFS_BLOCK_GROUP_RAID0 != 0:
		stripeSize = length
		stripeSize /= uint64(numStripes)
	case typeFlags&BTRFS_BLOCK_GROUP_RAID10 != 0:
		stripeSize = length * 2
		stripeSize /= uint64(numStripes)
	case typeFlags&BTRFS_BLOCK_GROUP_RAID5 != 0:
		stripeSize = length
		stripeSize /= (uint64(numStripes) - 1)
	case typeFlags&BTRFS_BLOCK_GROUP_RAID6 != 0:
		stripeSize = length
		stripeSize /= (uint64(numStripes) - 2)
	default:
		stripeSize = length
	}
	return stripeSize
}

// checkChunkRefs: checks that the chunk has a block group and that it matches the chunk
// for each stripe it also checks that the stripe exist in the device extent and that it matches
func checkChunkRefs(chunkRec *ChunkRecord, blockGroupTree *BlockGroupTree, devExtentTree *DeviceExtentTree, silent bool) bool {

	ret := true
	if i := blockGroupTree.Tree.Get(&BlockGroupRecord{CacheExtent: CacheExtent{Start: chunkRec.Offset, Size: chunkRec.Length}}); i != nil {
		blockGroupRec := i.(*BlockGroupRecord)
		if chunkRec.Length != blockGroupRec.Offset || chunkRec.Offset != blockGroupRec.Objectid || chunkRec.TypeFlags != blockGroupRec.Flags {
			if !silent {
				fmt.Printf("Chunk[%d, %d, %d]: length(%d), offset(%d), type(%d) mismatch with block group[%d, %d, %d]: offset(%d), objectid(%d), flags(%d)\n",
					chunkRec.Objectid,
					chunkRec.Type,
					chunkRec.Offset,
					chunkRec.Length,
					chunkRec.Offset,
					chunkRec.TypeFlags,
					blockGroupRec.Objectid,
					blockGroupRec.Type,
					blockGroupRec.Offset,
					blockGroupRec.Offset,
					blockGroupRec.Objectid,
					blockGroupRec.Flags)
				ret = false
			} else {
				chunkRec.List = nil
				chunkRec.BgRec = blockGroupRec
				//					list_del_init(&block_group_rec->list);
				//			chunk_rec->bg_rec = block_group_rec;
			}
		} else {
			if !silent {
				fmt.Printf("Chunk[%d, %d, %d]: length(%d), offset(%d), type(%d) is not found in block group\n",
					chunkRec.Objectid,
					chunkRec.Type,
					chunkRec.Offset,
					chunkRec.Length,
					chunkRec.Offset,
					chunkRec.TypeFlags)
				ret = false

			}
		}
	}
	length := calcStripeLength(chunkRec.TypeFlags, chunkRec.Length, chunkRec.NumStripes)
	for i := uint16(0); i < chunkRec.NumStripes; i++ {
		devid := chunkRec.Stripes[i].Devid
		offset := chunkRec.Stripes[i].Offset
		if item := devExtentTree.Tree.Get(&DeviceExtentRecord{CacheExtent: CacheExtent{Objectid: devid, Start: offset, Size: length}}); item != nil {
			devExtentRec := item.(*DeviceExtentRecord)
			if devExtentRec.Objectid != devid ||
				devExtentRec.Offset != offset ||
				devExtentRec.ChunkOffset != chunkRec.Offset ||
				devExtentRec.Length != length {
				if !silent {
					fmt.Printf(
						"Chunk[%d, %d, %d] stripe[%d, %d] dismatch dev extent[%d, %d, %d] local [%d, %d, %d]\n",
						chunkRec.Objectid,
						chunkRec.Type,
						chunkRec.Offset,
						chunkRec.Stripes[i].Devid,
						chunkRec.Stripes[i].Offset,
						devExtentRec.Objectid,
						devExtentRec.Offset,
						devExtentRec.Length,
						devid,
						offset,
						length)
					ret = false
				} else {
					devExtentRec.ChunkList = chunkRec.Dextents.PushBack(devExtentRec.ChunkList)
					//									list_move(&dev_extent_rec->chunk_list,
					//					  &chunk_rec->dextents);
				}
			}
		} else {
			if !silent {
				fmt.Printf(
					"Chunk[%d, %d, %d] stripe[%d, %d] is not found in dev extent\n",
					chunkRec.Objectid,
					chunkRec.Type,
					chunkRec.Offset,
					chunkRec.Stripes[i].Devid,
					chunkRec.Stripes[i].Offset)
				ret = false
			}
		}
	}
	return ret
}

// checkChunks: walks the chunk cache and checks the chunks references incrementing the good and bad chunk lists accordingliy
// if silent is false interates the block groups and device extent chunk oprphans list printing a could not find message
func checkChunks(chunkCache *llrb.LLRB, blockGroupTree *BlockGroupTree, deviceExtentTree *DeviceExtentTree, good, bad *list.List, silent bool) bool {

	var ret bool
	chunkCache.AscendGreaterOrEqual(llrb.Inf(-1), func(i llrb.Item) bool {
		chunkRec := i.(*ChunkRecord)
		err := checkChunkRefs(chunkRec, blockGroupTree, deviceExtentTree, silent)
		if err {
			ret = err
			if bad != nil {
				chunkRec.List = bad.PushBack(chunkRec)
			}
		} else {
			if good != nil {
				chunkRec.List = good.PushBack(chunkRec)
			}
		}
		return true
	})
	if !silent {
		for item := blockGroupTree.Block_Groups.Front(); item != nil; item = item.Next() {
			i := item.Value
			bgRec := i.(*BlockGroupRecord)
			fmt.Printf("Block group[%d, %d] (flags = %d) didn't find the relative chunk.\n",
				bgRec.Objectid,
				bgRec.Offset,
				bgRec.Flags)
		}
		for item := deviceExtentTree.ChunkOrphans.Front(); item != nil; item = item.Next() {
			i := item.Value
			dextRec := i.(*DeviceExtentRecord)
			fmt.Printf("Device extent[%d, %d, %d] didn't find the relative chunk.\n",
				dextRec.Objectid,
				dextRec.Offset,
				dextRec.Length)
		}
	}
	return ret
}

//btrfsGetDeviceExtents: counts the chunks in the orphan device extents list that match theis chunk and returns a list of them
func btrfsGetDeviceExtents(chunkObject uint64, orphanDevexts *list.List, ret_list *list.List) uint16 {

	count := uint16(0)
	for element := orphanDevexts.Front(); element != nil; element = element.Next() {
		i := element.Value
		devExt := i.(*DeviceExtentRecord)
		if devExt.ChunkOffset == chunkObject {
			//			list_move_tail(&devext->chunk_list, ret_list);
			count++
		}
	}
	return count
}

// calcSubNstripes caqclulates the number of sub stripes as 2 if this is a rad10 block 1 otherwise
func calcSubNstripes(Type uint64) uint16 {
	if Type&BTRFS_BLOCK_GROUP_RAID10 == 0 {
		return 2
	} else {
		return 1
	}
}

// btrfsRebuildOrderedMetaChunkStripes: rebuild stripes for raid block groups with existing devextents
func btrfsRebuildOrderedMetaChunkStripes(rc *RecoverControl, chunk *ChunkRecord) (error, bool) {
	var (
		start  = chunk.Offset
		end    = chunk.Offset + chunk.Length
		mirror uint32
	)
	cache := rc.EbCache.Get(&ExtentRecord{CacheExtent: CacheExtent{Start: start, Size: chunk.Length}})
	if cache == nil {
		return btrfsRebuildUnorderedChunkStripes(rc, chunk)
	}
	devExts := chunk.Dextents
	chunk.Dextents = list.New()
again:
	er := cache.(*ExtentRecord)
	index := btrfsCalcStripeIndex(chunk, er.CacheExtent.Start)
	if chunk.Stripes[index].Devid != 0 {
		goto next
	}
	for i := devExts.Front(); i != nil && i.Value != nil; i.Next() {
		devExt := i.Value.(*DeviceExtentRecord)
		if isExtentRecordInDeviceExtent(er, devExt, &mirror) {
			chunk.Stripes[index].Devid = devExt.Objectid
			chunk.Stripes[index].Offset = devExt.Offset
			chunk.Stripes[index].Uuid = er.Devices[mirror].Uuid
			index++
			//			list_move(&devext->chunk_list, &chunk->dextents);
			devExt.ChunkList = chunk.Dextents.PushBack(devExt.ChunkList)
		}
	}
next:
	start = btrfsNextStripeLogicalOffset(chunk, er.CacheExtent.Start)
	if start > end {
		goto noExtentRecord
	}
	cache = rc.EbCache.Get(&ExtentRecord{CacheExtent: CacheExtent{Start: start, Size: end - start}})
	if cache != nil {
		goto again
	}
noExtentRecord:
	if devExts.Len() == 0 {
		return nil, true
	}
	if chunk.TypeFlags&(BTRFS_BLOCK_GROUP_RAID5|BTRFS_BLOCK_GROUP_RAID6) == 0 {
		//		/* Fixme: try to recover the order by the parity block. */
		//		list_splice_tail(&devexts, &chunk->dextents);
		return errors.New("-EINVAL"), false
	}
	/* There is no data on the lost stripes, we can reorder them freely. */
	for index := uint16(0); index < chunk.NumStripes; index++ {
		if chunk.Stripes[index].Devid != 0 {
			continue
		}
		//		devext = list_first_entry(&devexts,
		//					  struct device_extent_record,
		//					   chunk_list);
		//		list_move(&devext->chunk_list, &chunk->dextents);
		if i := devExts.Front(); i != nil {
			devExt := i.Value.(*DeviceExtentRecord)
			devExt.ChunkList = chunk.Dextents.PushBack(devExt.ChunkList)
			chunk.Stripes[index].Devid = devExt.Objectid
			chunk.Stripes[index].Offset = devExt.Offset
			device := btrfsFindDeviceByDevid(rc.FsDevices, devExt.Objectid, 0)
			if device == nil {
				chunk.Dextents.PushBackList(devExts)
				return errors.New("-EINVAL"), false
			}
			chunk.Stripes[index].Uuid = device.Uuid
		}
	}
	return nil, true
}

// btrfsRebuildUnorderedChunkStripes rebuild stimple stripes
func btrfsRebuildUnorderedChunkStripes(rc *RecoverControl, chunk *ChunkRecord) (error, bool) {
	item := chunk.Dextents.Front()
	for i := uint16(0); item != nil && item.Value != nil && i < chunk.NumStripes; i++ {
		devExt := item.Value.(*DeviceExtentRecord)
		chunk.Stripes[i].Devid = devExt.Objectid
		chunk.Stripes[i].Offset = devExt.Offset
		device := btrfsFindDeviceByDevid(rc.FsDevices, devExt.Objectid, 0)
		if device == nil {
			return errors.New("-EINVAL"), false
		}
		chunk.Stripes[i].Uuid = device.Uuid
		item = item.Next()
	}
	return nil, true
}

// btrfsRebuildChunkStripes rebuild the rtripes fro the chunk
func btrfsRebuildChunkStripes(rc *RecoverControl, chunk *ChunkRecord) (error, bool) {
	/*
	 * All the data in the system metadata chunk will be dropped,
	 * so we need not guarantee that the data is right or not, that
	 * is we can reorder the stripes in the system metadata chunk.
	 */
	switch {
	case (chunk.TypeFlags&BTRFS_BLOCK_GROUP_METADATA) != 0 && (chunk.TypeFlags&BTRFS_ORDERED_RAID) != 0:
		return btrfsRebuildOrderedMetaChunkStripes(rc, chunk)
	case (chunk.TypeFlags&BTRFS_BLOCK_GROUP_DATA) != 0 && (chunk.TypeFlags&BTRFS_ORDERED_RAID) != 0:
		return nil, true /* Be handled after the fs is opened. */
	}
	return btrfsRebuildUnorderedChunkStripes(rc, chunk)
}
func btrfsVerifyDeviceExtents(bg *BlockGroupRecord, devExts *list.List, nDevExts uint16) bool {
	var (
		stripeLength       uint64
		expectedNumStripes uint16
	)
	expectedNumStripes = calcNumStripes(bg.Flags)
	if expectedNumStripes != 0 && expectedNumStripes != nDevExts {
		return false
	}
	stripeLength = calcStripeLength(bg.Flags, bg.Offset, nDevExts)
	for i := devExts.Front(); i != nil && i.Value != nil; i = i.Next() {
		devExt := i.Value.(*DeviceExtentRecord)
		if devExt.Length != stripeLength {
			return false
		}
	}
	return true
}

// btrfsRecoverChunks create the chunks by block group
func btrfsRecoverChunks(rc *RecoverControl) bool {
	ret := true
	devExts := list.New()
	for element := rc.Bg.Block_Groups.Front(); element != nil; element = element.Next() {
		i := element.Value
		bg := i.(*BlockGroupRecord)
		nstripes := btrfsGetDeviceExtents(bg.Objectid, rc.Devext.ChunkOrphans, devExts)
		chunk := &ChunkRecord{
			Dextents:    list.New(),
			BgRec:       bg,
			CacheExtent: CacheExtent{Start: bg.Objectid, Size: bg.Offset},
			Objectid:    BTRFS_FIRST_CHUNK_TREE_OBJECTID,
			Type:        BTRFS_CHUNK_ITEM_KEY,
			Offset:      bg.Objectid,
			Generation:  bg.Generation,
			Owner:       BTRFS_CHUNK_TREE_OBJECTID,
			StripeLen:   BTRFS_STRIPE_LEN,
			TypeFlags:   bg.Flags,
			IoWidth:     BTRFS_STRIPE_LEN,
			IoAlign:     BTRFS_STRIPE_LEN,
			SectorSize:  rc.Sectorsize,
			SubStripes:  calcSubNstripes(bg.Flags),
		}
		rc.Chunk.InsertNoReplace(chunk)
		if nstripes == 0 {
			chunk.List = rc.BadChunks.PushBack(chunk)
			continue
		}
		chunk.Dextents.PushBackList(devExts)

		if btrfsVerifyDeviceExtents(bg, devExts, nstripes) {
			continue
			rc.BadChunks.PushBack(chunk.List)
		}
		chunk.NumStripes = nstripes
		var err error
		err, ret = btrfsRebuildChunkStripes(rc, chunk)
		switch {
		case err != nil:
			rc.UnrepairedChunks.PushBack(chunk.List)
		case ret:
			rc.GoodChunks.PushBack(chunk.List)
		case !ret:
			rc.BadChunks.PushBack(chunk.List)
		}
	}
	/*
	 * Don't worry about the lost orphan device extents, they don't
	 * have its chunk and block group, they must be the old ones that
	 * we have dropped.
	 */
	return ret
}

func btrfsCalcStripeIndex(chunk *ChunkRecord, logical uint64) uint16 {
	var (
		offset                         = logical - chunk.Offset
		stripeNr, ntDataStripes, index uint16
	)

	stripeNr = uint16(offset / chunk.StripeLen)
	switch {
	case chunk.TypeFlags&BTRFS_BLOCK_GROUP_RAID0 != 0:
		index = stripeNr % chunk.NumStripes
	case chunk.TypeFlags&BTRFS_BLOCK_GROUP_RAID10 != 0:
		index = stripeNr % (chunk.NumStripes / chunk.SubStripes)
		index *= chunk.SubStripes
	case chunk.TypeFlags&BTRFS_BLOCK_GROUP_RAID5 != 0:
		ntDataStripes = chunk.NumStripes - 1
		index = stripeNr % ntDataStripes
		stripeNr /= ntDataStripes
		index = (index + stripeNr) % chunk.NumStripes
	case chunk.TypeFlags&BTRFS_BLOCK_GROUP_RAID6 != 0:
		ntDataStripes = chunk.NumStripes - 2
		index = stripeNr % ntDataStripes
		stripeNr /= ntDataStripes
		index = (index + stripeNr) % chunk.NumStripes
	}
	return index
}

func btrfsFindDeviceByDevid(fsDevices *BtrfsFsDevices, devid uint64, instance int) *BtrfsDevice {
	numFound := 0
	for i := fsDevices.Devices.Front(); i != nil && i.Value != nil; i = i.Next() {
		dev := i.Value.(*BtrfsDevice)
		if dev.Devid == devid && numFound == instance {
			return dev
		}
		numFound++
	}
	return nil
}
func calcNumStripes(flags uint64) uint16 {

	switch {
	case flags&(BTRFS_BLOCK_GROUP_RAID0|
		BTRFS_BLOCK_GROUP_RAID10|
		BTRFS_BLOCK_GROUP_RAID5|
		BTRFS_BLOCK_GROUP_RAID6) != 0:
		return 0

	case flags&(BTRFS_BLOCK_GROUP_RAID1|
		BTRFS_BLOCK_GROUP_DUP) != 0:
		return 2
	default:
		return 1
	}

}

func isExtentRecordInDeviceExtent(er *ExtentRecord, dext *DeviceExtentRecord, mirror *uint32) bool {

	for i := uint32(0); i < er.Nmirrors; i++ {
		if er.Devices[i].Devid == dext.Objectid &&
			er.Offsets[i] >= dext.Offset &&
			er.Offsets[i] < dext.Offset+dext.Length {
			*mirror = i
			return true
		}
	}
	return false
}
func btrfsNextStripeLogicalOffset(chunk *ChunkRecord, logical uint64) uint64 {

	var offset = logical - chunk.Offset
	offset /= chunk.StripeLen
	offset *= chunk.StripeLen
	offset += chunk.StripeLen
	return offset + chunk.Offset
}

// TODO
// Proper support for multiple devices and device scan
