package btrfs

import ()

const (
	BTRFS_SIZE              = 4096
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
	U8  byte
	U16 uint16
	U32 uint32
	U64 uint64

	Char byte
	Le16 uint16
	Le32 uint32
	Le64 uint64
)
