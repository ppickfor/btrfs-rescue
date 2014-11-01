package btrfs

import ()

var (
	btrfsCsumSizes = []uint16{4, 0}
)

func BtrfsSuperCsumSize(s *BtrfsSuperBlock) uint16 {
	t := s.CsumType
	//BUG_ON(t >= ARRAY_SIZE(btrfsCsumSizes));
	return btrfsCsumSizes[t]
}
