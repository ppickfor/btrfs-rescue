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
	BTRFS_NUM_BACKUP_ROOTS    = 4
	BTRFS_MAGIC               = 0x4D5F53665248425F /* ascii _BHRfS_M, no null */
	BTRFS_HEADER_FLAG_WRITTEN = (1 << 0)
	BTRFS_HEADER_FLAG_RELOC   = (1 << 1)
	BTRFS_SUPER_FLAG_SEEDING  = (1 << 32)
	BTRFS_SUPER_FLAG_METADUMP = (1 << 33)
)

type (
	u64  uint64
	u8   byte
	le64 uint64
	le32 uint32
	char byte
	le16 uint16
	u16  uint16
	u32  uint32
)

// The flag package provides a default help printer via -h switch
var (
	versionFlag     *bool   = flag.Bool("v", false, "Print the version number.")
	deviceFlag      *string = flag.String("d", DEFAULT_DEVICE, "The device to scan")
	blocksFlag      *int64  = flag.Int64("n", DEFAULT_BLOCKS, "The number of BTRFS_SIZE blocks to read")
	startblocksFlag *int64  = flag.Int64("s", 0, "The number of BTRFS_SIZE blocks to start at")

	btrfs_csum_sizes = []int{4, 0}
)

//
func btrfs_sb_offset(mirror int) uint64 {
	var start uint64 = 16 * 1024
	if mirror != 0 {
		return start << uint64(BTRFS_SUPER_MIRROR_SHIFT*mirror)
	}
	return BTRFS_SUPER_INFO_OFFSET
}

type rb_root struct {
}
type cache_tree struct {
	root rb_root
}

type list_head struct {
	next *list_head
	prev *list_head
}

type btrfs_root_backup struct {
	Tree_root     le64
	Tree_root_gen le64

	Chunk_root     le64
	Chunk_root_gen le64

	Extent_root     le64
	Extent_root_gen le64

	Fs_root     le64
	Fs_root_gen le64

	Dev_root     le64
	Dev_root_gen le64

	Csum_root     le64
	Csum_root_gen le64

	Total_bytes le64
	Bytes_used  le64
	Num_devices le64
	/* future */
	Unsed_64 [4]le64

	Tree_root_level   u8
	Chunk_root_level  u8
	Extent_root_level u8
	Fs_root_level     u8
	Dev_root_level    u8
	Csum_root_level   u8
	/* future and to align */
	Unused_8 [10]u8
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
	Csum [BTRFS_CSUM_SIZE]u8
	/* the first 3 fields must match struct btrfs_header */
	Fsid   [BTRFS_FSID_SIZE]u8 /* FS specific uuid */
	Bytenr le64                /* this block number */
	Flags  le64

	/* allowed to be different from the btrfs_header from here own down */
	Magic      le64
	Generation le64
	Root       le64
	Chunk_root le64
	Log_root   le64

	/* this will help find the new super based on the log root */
	Log_root_transid      le64
	Total_bytes           le64
	Bytes_used            le64
	Root_dir_objectid     le64
	Num_devices           le64
	Sectorsize            le32
	Nodesize              le32
	Leafsize              le32
	Stripesize            le32
	Sys_chunk_array_size  le32
	Chunk_root_generation le64
	Compat_flags          le64
	Compat_ro_flags       le64
	Incompat_flags        le64
	Csum_type             le16
	Root_level            u8
	Chunk_root_level      u8
	Log_root_level        u8
	Dev_item              btrfs_dev_item

	Label [BTRFS_LABEL_SIZE]char

	Cache_generation     le64
	Uuid_tree_generation le64

	/* future expansion */
	Reserved        [30]le64
	Sys_chunk_array [BTRFS_SYSTEM_CHUNK_ARRAY_SIZE]u8
	Super_roots     [BTRFS_NUM_BACKUP_ROOTS]btrfs_root_backup
}
type btrfs_fs_devices struct {
	fsid [BTRFS_FSID_SIZE]u8 /* FS specific uuid */

	/* the device with this id has the most recent copy of the super */
	latest_devid u64
	latest_trans u64
	lowest_devid u64
	latest_bdev  int
	lowest_bdev  int
	devices      list_head
	list         list_head

	seeding int
	seed    *btrfs_fs_devices
}
type block_group_tree struct {
	tree         cache_tree
	block_groups list_head
}
type device_extent_tree struct {
	tree cache_tree
	/*
	 * The idea is:
	 * When checking the chunk information, we move the device extents
	 * that has its chunk to the chunk's device extents list. After the
	 * check, if there are still some device extents in no_chunk_orphans,
	 * it means there are some device extents which don't belong to any
	 * chunk.
	 *
	 * The usage of no_device_orphans is the same as the first one, but it
	 * is for the device information check.
	 */
	no_chunk_orphans  list_head
	no_device_orphans list_head
}
type recover_control struct {
	verbose int
	yes     int

	csum_size             u16
	sectorsize            u32
	leafsize              u32
	generation            u64
	chunk_root_generation u64

	fs_devices *btrfs_fs_devices

	chunk    cache_tree
	bg       block_group_tree
	devext   device_extent_tree
	eb_cache cache_tree

	good_chunks       list_head
	bad_chunks        list_head
	unrepaired_chunks list_head
	//pthread_mutex_t rc_lock;
}

func btrfs_super_csum_size(s *btrfs_super_block) int {
	t := s.Csum_type
	//BUG_ON(t >= ARRAY_SIZE(btrfs_csum_sizes));
	return btrfs_csum_sizes[t]
}

// reads super block into sb at offset sb_bytenr
// if super_recover is != 0 the read superblock backups ad find latest generation
func btrfs_read_dev_super(fd int, sb *btrfs_super_block, sb_bytenr u64, super_recover int) bool {
	var (
		fsid                [BTRFS_FSID_SIZE]u8
		fsid_is_initialized bool = false
		buf                 *btrfs_super_block
		//		i                   int
		//		ret                 int
		max_super int  = 1
		transid   le64 = 0
		bytenr    int64
		//		err                 error
	)
	buf = new(btrfs_super_block)
	if super_recover != 0 {
		max_super = BTRFS_SUPER_MIRROR_MAX
	}
	var size = binary.Size(buf)
	var bytebuf = make([]byte, size)
	var bytebr = bytes.NewReader(bytebuf)
	// dont look like this will be executed
	if sb_bytenr != BTRFS_SUPER_INFO_OFFSET {
		ret, _ := syscall.Pread(fd, bytebuf, int64(sb_bytenr))
		if ret < size {
			return false
		}
		_ = binary.Read(bytebr, binary.LittleEndian, buf)
		if buf.Bytenr != le64(sb_bytenr) ||
			buf.Magic != BTRFS_MAGIC {
			return false
		}
		*sb = *buf
		return true
	}

	/*
	* we would like to check all the supers, but that would make
	* a btrfs mount succeed after a mkfs from a different FS.
	* So, we need to add a special mount option to scan for
	* later supers, using BTRFS_SUPER_MIRROR_MAX instead
	 */

	for i := 0; i < max_super; i++ {
		fmt.Printf("i: %v\n", i)

		bytenr = int64(btrfs_sb_offset(i))
		fmt.Printf("bytenr: %v\n", bytenr)
		ret, _ := syscall.Pread(fd, bytebuf, bytenr)
		//fmt.Printf("err: %v, ret: %v. bytebuf: %v\n", err, ret, bytebuf)
		if ret < size {
			break
		}
		bytebr = bytes.NewReader(bytebuf)
		_ = binary.Read(bytebr, binary.LittleEndian, buf)
		if buf.Bytenr != le64(bytenr) {
			fmt.Printf("bad bytent: should be %v not %v\n",bytenr,buf.Bytenr)
			continue
		}

		/* if magic is NULL, the device was removed */
		if buf.Magic == 0 && i == 0 {
			return false
		}
		if buf.Magic != BTRFS_MAGIC {
			fmt.Printf("bad magic %x not %x\n",buf.Magic,BTRFS_MAGIC)
			continue
		}

		if !fsid_is_initialized {

			fsid = buf.Fsid
			fsid_is_initialized = true
		} else if fsid != buf.Fsid {
			/*
			 * the superblocks (the original one and
			 * its backups) contain data of different
			 * filesystems -> the super cannot be trusted
			 */
			fmt.Printf("bad fsid %x not %x\n",fsid ,buf.Fsid)
			continue
		}

		if buf.Generation > transid {
			*sb = *buf
//			fmt.Printf("buf: %+v\n\n\nsb inside: %+v\n", buf, sb)
			transid = buf.Generation
		}
	}
	if transid > 0 {
		return true
	} else {
		return false
	}

}

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
func main() {
	
	var empty []byte = make([]byte, BTRFS_SIZE)
	var sb btrfs_super_block
	var rc recover_control

	flag.Parse() // Scan the arguments list

	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
	}

	fd, err := syscall.Open(*deviceFlag, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	ret := btrfs_read_dev_super(fd, &sb, BTRFS_SUPER_INFO_OFFSET, 1)
	if ! ret  {
		log.Fatal("read super block error\n")
	}
	fmt.Printf("SB: %+v\n", sb)
	rc.sectorsize = u32(sb.Sectorsize)
	rc.leafsize = u32(sb.Leafsize)
	rc.generation = u64(sb.Generation)
	rc.chunk_root_generation = u64(sb.Chunk_root_generation)
	rc.csum_size = u16(btrfs_super_csum_size(&sb))
	fmt.Printf("RC: %+v\n", rc)
	var buf []byte = make([]byte, rc.leafsize)
	/* if seed, the result of scanning below will be partial */
	if (sb.Flags & BTRFS_SUPER_FLAG_SEEDING) != 0 {
		fmt.Errorf("this device is seed device\n")
		ret = false
		goto fail_free_sb
	}
	for bytenr := int64(*startblocksFlag * BTRFS_SIZE); bytenr < BTRFS_SIZE*(*blocksFlag); bytenr += int64(rc.leafsize) {
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
