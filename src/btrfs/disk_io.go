package btrfs

import (
	"bytes"
	"fmt"
	"encoding/binary"
	"syscall"
)

//
func btrfs_sb_offset(mirror int) uint64 {
	var start uint64 = 16 * 1024
	if mirror != 0 {
		return start << uint64(BTRFS_SUPER_MIRROR_SHIFT*mirror)
	}
	return BTRFS_SUPER_INFO_OFFSET
}

// reads super block into sb at offset sb_bytenr
// if super_recover is != 0 the read superblock backups ad find latest generation
func Btrfs_read_dev_super(fd int, sb *Btrfs_super_block, sb_bytenr U64, super_recover int) bool {
	var (
		fsid                [BTRFS_FSID_SIZE]U8
		fsid_is_initialized bool = false
		buf                 *Btrfs_super_block
		//		i                   int
		//		ret                 int
		max_super int  = 1
		transid   Le64 = 0
		bytenr    int64
		//		err                 error
	)
	buf = new(Btrfs_super_block)
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
		if buf.Bytenr != Le64(sb_bytenr) ||
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
		if buf.Bytenr != Le64(bytenr) {
			fmt.Printf("bad bytent: should be %v not %v\n", bytenr, buf.Bytenr)
			continue
		}

		/* if magic is NULL, the device was removed */
		if buf.Magic == 0 && i == 0 {
			return false
		}
		if buf.Magic != BTRFS_MAGIC {
			fmt.Printf("bad magic %x not %x\n", buf.Magic, BTRFS_MAGIC)
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
			fmt.Printf("bad fsid %x not %x\n", fsid, buf.Fsid)
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
