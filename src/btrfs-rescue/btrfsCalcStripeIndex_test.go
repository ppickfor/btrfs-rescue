package main

import (
	. "btrfs"
	"testing"
)

func TestBtrfsCalcStripeIndex(t *testing.T) {
	chunk := &ChunkRecord{
		Offset:     0,
		StripeLen:  BTRFS_STRIPE_LEN,
		NumStripes: 0,
		SubStripes: 0,
		TypeFlags:  0,
	}
	logical := uint64(0)
	index := uint16(0)
	index = btrfsCalcStripeIndex(chunk, logical)
	if index != 0 {
		t.Errorf("index !=0")
	}

	chunk.TypeFlags = BTRFS_BLOCK_GROUP_RAID0
	chunk.NumStripes = 2
	chunk.Offset = 10000
	index = btrfsCalcStripeIndex(chunk, 30000+64*1024)
	if index != 1 {
		t.Errorf("index !=1, is %v", index)
	}
}
