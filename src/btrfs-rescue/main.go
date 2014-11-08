package main

import (
	. "btrfs"
	"bufio"
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"strconv"
	"sync"

	"code.google.com/p/go.net/context"
	"github.com/monnand/GoLLRB/llrb"
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
	CheckChunks(rc.Chunk, &rc.Bg, &rc.Devext, rc.GoodChunks, rc.BadChunks, false)

	BtrfsRecoverChunks(rc)

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
							ExtractMetadataRecord(rc, generation, items, infoByteBlock)
						}
					case BTRFS_CHUNK_TREE_OBJECTID:
						if generation <= rc.ChunkRootGeneration {
							items := make([]BtrfsItem, nritems)
							_ = binary.Read(bytereader, binary.LittleEndian, items)
							ExtractMetadataRecord(rc, generation, items, infoByteBlock)
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
