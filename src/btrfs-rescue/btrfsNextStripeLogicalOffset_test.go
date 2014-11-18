package main

import (
	. "btrfs"
	"container/list"
	"testing"
)

func TestBtrfsNextStripeLogicalOffset(t *testing.T) {
	chunk := &ChunkRecord{
		CacheExtent: CacheExtent{Objectid: 0x0, Start: 0x3, Size: 0x6},
		List:        (*list.Element)(nil),
		//		Dextents:    (*list.List)(0xc208024540),
		BgRec:      (*BlockGroupRecord)(nil),
		Generation: 0x63,
		Objectid:   0x1,
		Type:       0x2,
		Offset:     1024,
		Owner:      0x7,
		Length:     0x6,
		TypeFlags:  0x9,
		StripeLen:  BTRFS_STRIPE_LEN,
		NumStripes: 0x1,
		SubStripes: 0x0,
		IoAlign:    0xa,
		IoWidth:    0xb,
		SectorSize: 0xc,
		Stripes: []Stripe{
			Stripe{
				Devid:  0x1,
				Offset: 0x64,
				Uuid:   [16]uint8{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf},
			},
		},
	}
	x := btrfsNextStripeLogicalOffset(chunk, 1024*1000)
	if x != 1049600 {
		t.Errorf("btrfsNextStripeLogicalOffset(chunk, 1024*1000) != 1049600 With chunk offset 1024, = %d", x)
	}

}
