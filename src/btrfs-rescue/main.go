package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"btrfs"
	"syscall"
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
)

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
func main() {

	var empty []byte = make([]byte, btrfs.BTRFS_SIZE)
	var sb btrfs.Btrfs_super_block
	var rc btrfs.Recover_control

	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	fd, err := syscall.Open(*deviceFlag, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	ret := btrfs.Btrfs_read_dev_super(fd, &sb, btrfs.BTRFS_SUPER_INFO_OFFSET, 1)
	if !ret {
		log.Fatal("read super block error\n")
	}
	fmt.Printf("\nSB: %+v\n", sb)
	rc.Sectorsize = btrfs.U32(sb.Sectorsize)
	rc.Leafsize = btrfs.U32(sb.Leafsize)
	rc.Generation = btrfs.U64(sb.Generation)
	rc.Chunk_root_generation = btrfs.U64(sb.Chunk_root_generation)
	rc.Csum_size = btrfs.U16(btrfs.Btrfs_super_csum_size(&sb))
	fmt.Printf("\nRC: %+v\n", rc)
	var buf []byte = make([]byte, rc.Leafsize)
	/* if seed, the result of scanning below will be partial */
	if (sb.Flags & btrfs.BTRFS_SUPER_FLAG_SEEDING) != 0 {
		fmt.Errorf("this device is seed device\n")
		ret = false
		goto fail_free_sb
	}
	for bytenr := int64(*startblocksFlag * btrfs.BTRFS_SIZE); bytenr < btrfs.BTRFS_SIZE*(*blocksFlag); bytenr += int64(rc.Leafsize) {
		n, err := syscall.Pread(fd, buf, bytenr)
		if err != nil {
			log.Fatalln(os.NewSyscallError("pread64", err))
		}
		if n == 0 {
			break
		}
		if bytes.Equal(buf, empty) {
			continue
		}
		//		fmt.Printf("Read %v bytes from %v @%v\n", n, fd, bytenr)
		//		for _, elem := range buf {
		//			if isAlphaNum(elem) {
		//				fmt.Print(string(elem))
		//			} else {
		//				fmt.Print(".")
		//			}
		//
		//		}
		//		fmt.Println("")
	}
fail_free_sb:
}
