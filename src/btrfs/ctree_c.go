package btrfs

import ()

var (
	btrfs_csum_sizes = []int{4, 0}
)

func Btrfs_super_csum_size(s *Btrfs_super_block) int {
	t := s.Csum_type
	//BUG_ON(t >= ARRAY_SIZE(btrfs_csum_sizes));
	return btrfs_csum_sizes[t]
}