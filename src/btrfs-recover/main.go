package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

const (
	APP_VERSION    = "0.1"
	BTRFS_SIZE     = 4096
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

func main() {
	var buf []byte = make([]byte, BTRFS_SIZE)

	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	fd, err := syscall.Open(*deviceFlag, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	for offset := int64(*startblocksFlag * BTRFS_SIZE); offset < BTRFS_SIZE*(*blocksFlag); offset += BTRFS_SIZE {
		n, err := syscall.Pread(fd, buf, offset)
		if n == 0 {
			break
		}
		fmt.Printf("Read %v bytes from %v @%v\n", n, fd, offset)
		if err != nil {
			log.Fatalln(os.NewSyscallError("pread64", err))
		}
	}

}
