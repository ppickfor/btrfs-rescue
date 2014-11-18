package main

import (
	. "btrfs"
	"testing"
)

func TestCalcNumStripes(t *testing.T) {
	var x uint16
	x = calcNumStripes(0)
	if x != 1 {
		t.Errorf("calcNumStripes(0) != 1, = %v", x)
	}
	x = calcNumStripes(BTRFS_BLOCK_GROUP_RAID0)
	if x != 0 {
		t.Errorf("calcNumStripes(BTRFS_BLOCK_GROUP_RAID0) != 0, = %v", x)
	}
	x = calcNumStripes(BTRFS_BLOCK_GROUP_RAID1)
	if x != 2 {
		t.Errorf("calcNumStripes(BTRFS_BLOCK_GROUP_RAID1) != 2, = %v", x)
	}
	x = calcNumStripes(0)
	if x != 1 {
		t.Errorf("calcNumStripes(0) != 1, = %v", x)
	}
}
