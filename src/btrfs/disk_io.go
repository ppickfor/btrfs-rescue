package btrfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"syscall"
	//	"unsafe"
)

var (
	crc32c  = crc32.MakeTable(crc32.Castagnoli)
	
)

func btrfs_scan_fs_devices(
	fd int,
	path string,
	fs_devices **Btrfs_fs_devices,
	sb_bytenr uint64,
	run_ioctl bool,
	super_recover bool) int {
	//	u64 total_devs;
	//	int ret;
	//	if (!sb_bytenr)
	//		sb_bytenr = BTRFS_SUPER_INFO_OFFSET;
	//
	//	ret = btrfs_scan_one_device(fd, path, fs_devices,
	//				    &total_devs, sb_bytenr, super_recover);
	//	if (ret) {
	//		fprintf(stderr, "No valid Btrfs found on %s\n", path);
	//		return ret;
	//	}
	//
	//	if (total_devs != 1) {
	//		ret = btrfs_scan_for_fsid(run_ioctl);
	//		if (ret)
	//			return ret;
	//	}
	return 0
}

// calculate byte offset of superblock mirror in partition
func Btrfs_sb_offset(
	mirror int) uint64 {

	var start uint64 = 16 * 1024
	if mirror != 0 {
		return start << uint64(BTRFS_SUPER_MIRROR_SHIFT*mirror)
	}
	return BTRFS_SUPER_INFO_OFFSET
}

// reads super block into sb at offset sb_bytenr
// if super_recover is != 0 the read superblock backups ad find latest generation
func btrfs_read_dev_super(
	fd int,
	sb *Btrfs_super_block,
	sb_bytenr uint64,
	super_recover bool) bool {

	var (
		fsid                [BTRFS_FSID_SIZE]uint8
		fsid_is_initialized bool = false
		buf                 *Btrfs_super_block
		//		i                   int
		//		ret                 int
		max_super int    = 1
		transid   uint64 = 0
		bytenr    int64
		//		err                 error
	)
	buf = new(Btrfs_super_block)
	if super_recover {
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
		if buf.Bytenr != uint64(sb_bytenr) ||
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

		bytenr = int64(Btrfs_sb_offset(i))
		fmt.Printf("bytenr: %v\n", bytenr)
		ret, _ := syscall.Pread(fd, bytebuf, bytenr)
		//fmt.Printf("err: %v, ret: %v. bytebuf: %v\n", err, ret, bytebuf)
		if ret < size {
			break
		}
		bytebr = bytes.NewReader(bytebuf)
		_ = binary.Read(bytebr, binary.LittleEndian, buf)
		if buf.Bytenr != uint64(bytenr) {
			fmt.Printf("bad bytent: should be %v not %v\n", bytenr, buf.Bytenr)
			fmt.Printf("Super block:\n%+v\n", buf)
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
		if !check_super(fd,uint64(bytenr), buf) {
//			continue
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

func csum_tree_block_size(buf *Extent_buffer, csum_size uint16,
	verify bool, silent bool) bool {
	var (
		crc  uint32 = ^uint32(0)
		csum uint32
	)
	bytebr := bytes.NewReader(buf.Data)
	binary.Read(bytebr, binary.LittleEndian, csum)

	crc = crc32.Update(crc, crc32c, buf.Data[BTRFS_CSUM_SIZE:buf.Len])

	if csum != crc {
		if verify {
			if !silent {
				fmt.Errorf("checksum verify failed on %llu found %08X wanted %08X\n",
					buf.Start,
					crc,
					csum)
			}
			return false
		}
	} else {
		bytewbr := new(bytes.Buffer)
		err := binary.Write(bytewbr, binary.LittleEndian, crc)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
		copy(buf.Data, bytewbr.Bytes())
	}
	return true
}

func verify_tree_block_csum_silent(buf *Extent_buffer, csum_size uint16) bool {
	return csum_tree_block_size(buf, csum_size, true, true)
}

func check_super(fd int ,bytenr uint64, sb *Btrfs_super_block) bool {

	switch {
	case sb.Bytenr != bytenr:
		fmt.Printf("Bytenr mismatch calculated %v have %v\n", sb.Bytenr, bytenr)
		return false
	case sb.Magic != BTRFS_MAGIC:
		fmt.Printf("Magic mismatch calculated %v have %v\n", sb.Magic, BTRFS_MAGIC)
		return false
	}
	var (
		crc  uint32 = ^uint32(0)
//		crc  uint32 = uint32(0)
		csum uint32
	)
	fmt.Printf("Reading crc32c\n")
	bytebr := bytes.NewReader(sb.Csum[:])
	binary.Read(bytebr, binary.LittleEndian, &csum)

	fmt.Printf("Calculating crc32c\n")

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, *sb)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	mybytes := make([]byte,4096)
	ret, _ := syscall.Pread(fd, mybytes, int64(bytenr))
	if ret != 4096 {
		fmt.Printf("Pread failed\n")
		
	}

	crc = crc32.Checksum( mybytes[BTRFS_CSUM_SIZE:],crc32c)
	
	if csum != crc {
		fmt.Printf("Crc mismatch calculated %08x have %08x\nLen bytes: %v\n%v\n", crc, csum,len(mybytes),mybytes)
		return true
	}
	fmt.Printf("Crc match calculated %08x have %08x\n", crc, csum)
	return true
	
}
