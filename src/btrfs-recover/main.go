package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	//"strings"
	"encoding/binary"
	"syscall"
)

const (
	APP_VERSION    = "0.1"
	BTRFS_SIZE     = 4096
	DEFAULT_BLOCKS = 1024
	DEFAULT_DEVICE = "/dev/sda"

	BTRFS_SUPER_INFO_OFFSET = (64 * 1024)
	BTRFS_SUPER_INFO_SIZE   = 4096

	BTRFS_SUPER_MIRROR_MAX   = 3
	BTRFS_SUPER_MIRROR_SHIFT = 12

	BTRFS_CSUM_SIZE = 32

	BTRFS_FSID_SIZE = 16
	BTRFS_UUID_SIZE = 16

	BTRFS_SYSTEM_CHUNK_ARRAY_SIZE = 2048
	BTRFS_LABEL_SIZE              = 256

	/*
	 * just in case we somehow lose the roots and are not able to mount,
	 * we store an array of the roots from previous transactions
	 * in the super.
	 */
	BTRFS_NUM_BACKUP_ROOTS = 4
	BTRFS_MAGIC            = 0x4D5F53665248425F /* ascii _BHRfS_M, no null */
)

type (
	u64  uint64
	u8   byte
	le64 uint64
	le32 uint32
	char byte
	le16 uint16
)

// The flag package provides a default help printer via -h switch
var (
	versionFlag     *bool   = flag.Bool("v", false, "Print the version number.")
	deviceFlag      *string = flag.String("d", DEFAULT_DEVICE, "The device to scan")
	blocksFlag      *int64  = flag.Int64("n", DEFAULT_BLOCKS, "The number of BTRFS_SIZE blocks to read")
	startblocksFlag *int64  = flag.Int64("s", 0, "The number of BTRFS_SIZE blocks to start at")
	sb              btrfs_super_block
)

func btrfs_sb_offset(mirror int) uint64 {
	var start uint64 = 16 * 1024
	if mirror != 0 {
		return start << uint64(BTRFS_SUPER_MIRROR_SHIFT*mirror)
	}
	return BTRFS_SUPER_INFO_OFFSET
}

type btrfs_root_backup struct {
	tree_root     le64
	tree_root_gen le64

	chunk_root     le64
	chunk_root_gen le64

	extent_root     le64
	extent_root_gen le64

	fs_root     le64
	fs_root_gen le64

	dev_root     le64
	dev_root_gen le64

	csum_root     le64
	csum_root_gen le64

	total_bytes le64
	bytes_used  le64
	num_devices le64
	/* future */
	unsed_64 [4]le64

	tree_root_level   u8
	chunk_root_level  u8
	extent_root_level u8
	fs_root_level     u8
	dev_root_level    u8
	csum_root_level   u8
	/* future and to align */
	unused_8 [10]u8
}

type btrfs_dev_item struct {
	/* the internal btrfs device id */
	Devid le64

	/* size of the device */
	Total_bytes le64

	/* bytes used */
	Bytes_used le64

	/* optimal io alignment for this device */
	Io_align le32

	/* optimal io width for this device */
	Io_width le32

	/* minimal io size for this device */
	Sector_size le32

	/* type and info about this device */
	Type le64

	/* expected generation for this device */
	Generation le64

	/*
	 * starting byte of this partition on the device,
	 * to allowr for stripe alignment in the future
	 */
	Start_offset le64

	/* grouping information for allocation decisions */
	Dev_group le32

	/* seek speed 0-100 where 100 is fastest */
	Seek_speed u8

	/* bandwidth 0-100 where 100 is fastest */
	Bandwidth u8

	/* btrfs generated uuid for this device */
	Uuid [BTRFS_UUID_SIZE]u8

	/* uuid of FS who owns this device */
	Fsid [BTRFS_UUID_SIZE]u8
}

type btrfs_super_block struct {
	csum [BTRFS_CSUM_SIZE]u8
	/* the first 3 fields must match struct btrfs_header */
	fsid   [BTRFS_FSID_SIZE]u8 /* FS specific uuid */
	bytenr le64                /* this block number */
	flags  le64

	/* allowed to be different from the btrfs_header from here own down */
	magic      le64
	generation le64
	root       le64
	chunk_root le64
	log_root   le64

	/* this will help find the new super based on the log root */
	log_root_transid      le64
	total_bytes           le64
	bytes_used            le64
	root_dir_objectid     le64
	num_devices           le64
	sectorsize            le32
	nodesize              le32
	leafsize              le32
	stripesize            le32
	sys_chunk_array_size  le32
	chunk_root_generation le64
	compat_flags          le64
	compat_ro_flags       le64
	incompat_flags        le64
	csum_type             le16
	root_level            u8
	chunk_root_level      u8
	log_root_level        u8
	dev_item              btrfs_dev_item

	label [BTRFS_LABEL_SIZE]char

	cache_generation     le64
	uuid_tree_generation le64

	/* future expansion */
	reserved        [30]le64
	sys_chunk_array [BTRFS_SYSTEM_CHUNK_ARRAY_SIZE]u8
	super_roots     [BTRFS_NUM_BACKUP_ROOTS]btrfs_root_backup
}

func btrfs_read_dev_super(fd int, sb *btrfs_super_block, sb_bytenr u64,
	super_recover int) int {
	var (
		fsid                [BTRFS_FSID_SIZE]u8
		fsid_is_initialized bool = false
		buf                 btrfs_super_block
		//		i                   int
		//		ret                 int
		max_super int  = 1
		transid   le64 = 0
		bytenr    u64
		//		err                 error
	)

	if super_recover != 0 {
		max_super = BTRFS_SUPER_MIRROR_MAX
	}
	var size = binary.Size(buf)
	var bytebuf = make([]byte, size)
	var bytebr = bytes.NewReader(bytebuf)
	if sb_bytenr != BTRFS_SUPER_INFO_OFFSET {
		ret, _ := syscall.Pread(fd, bytebuf, int64(sb_bytenr))
		if ret < size {
			return -1
		}
		_ = binary.Read(bytebr, binary.LittleEndian, buf)
		if buf.bytenr != le64(sb_bytenr) ||
			buf.magic != BTRFS_MAGIC {
			return -1
		}

		return 0
	}

	/*
	* we would like to check all the supers, but that would make
	* a btrfs mount succeed after a mkfs from a different FS.
	* So, we need to add a special mount option to scan for
	* later supers, using BTRFS_SUPER_MIRROR_MAX instead
	 */

	for i := 0; i < max_super; i++ {
		bytenr = u64(btrfs_sb_offset(i))
		ret, _ := syscall.Pread(fd, bytebuf, int64(sb_bytenr))
		if ret < size {
			break
		}
		_ = binary.Read(bytebr, binary.LittleEndian, buf)
		if buf.bytenr != le64(bytenr) {
			continue
		}

		/* if magic is NULL, the device was removed */
		if buf.magic == 0 && i == 0 {
			return -1
		}
		if buf.magic == BTRFS_MAGIC {
			continue
		}

		if !fsid_is_initialized {

			fsid = buf.fsid
			fsid_is_initialized = true
		} else if fsid != buf.fsid {
			/*
			 * the superblocks (the original one and
			 * its backups) contain data of different
			 * filesystems -> the super cannot be trusted
			 */
			continue
		}

		if buf.generation > transid {
			sb = &buf
			transid = buf.generation
		}
	}
	if transid > 0 {
		return 0
	} else {
		return -1
	}

}

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
func main() {
	var buf []byte = make([]byte, BTRFS_SIZE)
	var empty []byte = make([]byte, BTRFS_SIZE)

	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	fd, err := syscall.Open(*deviceFlag, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	for bytenr := int64(*startblocksFlag * BTRFS_SIZE); bytenr < BTRFS_SIZE*(*blocksFlag); bytenr += BTRFS_SIZE {
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
		fmt.Printf("Read %v bytes from %v @%v\n", n, fd, bytenr)
		for _, elem := range buf {
			if isAlphaNum(elem) {
				fmt.Print(string(elem))
			} else {
				fmt.Print(".")
			}

		}
		fmt.Println("")
	}

}
