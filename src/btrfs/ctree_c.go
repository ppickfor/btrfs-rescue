package btrfs

import ()

var (
	btrfs_csum_sizes = []uint16{4, 0}
)

func Btrfs_super_csum_size(s *Btrfs_super_block) uint16 {
	t := s.Csum_type
	//BUG_ON(t >= ARRAY_SIZE(btrfs_csum_sizes));
	return btrfs_csum_sizes[t]
}