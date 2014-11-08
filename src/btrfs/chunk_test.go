package btrfs

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestNewChunkRecord(t *testing.T) {

	item := BtrfsItem{
		Key: BtrfsDiskKey{
			Objectid: 1,
			Type:     2,
			Offset:   3,
		},
		Offset: 25,
		Size:   80,
	}
	btrfsChunk := BtrfsChunk{
		Length:     6,
		Owner:      7,
		StripeLen:  8,
		Type:       9,
		IoAlign:    10,
		IoWidth:    11,
		SectorSize: 12,
		NumStripes: 1,
		SubStripes: 0,
	}
	stripe := Stripe{
		Devid:  1,
		Offset: 100,
		Uuid:   [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
	}
	bytewbr := new(bytes.Buffer)
	var err error
	err = binary.Write(bytewbr, binary.LittleEndian, item)
	if err != nil {
		t.Errorf("binary.Write failed:", err)
	}
	//	fmt.Printf("offset: %d\n", bytewbr.Len())
	err = binary.Write(bytewbr, binary.LittleEndian, btrfsChunk)
	if err != nil {
		t.Errorf("binary.Write failed:", err)
	}
	//	fmt.Printf("offset: %d\n", bytewbr.Len())
	err = binary.Write(bytewbr, binary.LittleEndian, stripe)
	if err != nil {
		t.Errorf("binary.Write failed:", err)
	}
	//	fmt.Printf("offset: %d\n", bytewbr.Len())
	chunk := NewChunkRecord(99, &item, bytewbr.Bytes())
	if chunk == nil {
		t.Errorf("No new ChunkRecord\n")
	}
	if chunk.Dextents == nil {
		t.Errorf("No chunk.Dextents list ceated\n")
	}
	chunk.Dextents = nil
	x := &ChunkRecord{
		CacheExtent: CacheExtent{Objectid: 0x0, Start: 0x3, Size: 0x6},
		List:        (*list.Element)(nil),
		//		Dextents:    (*list.List)(0xc208024540),
		BgRec:      (*BlockGroupRecord)(nil),
		Generation: 0x63,
		Objectid:   0x1,
		Type:       0x2,
		Offset:     0x3,
		Owner:      0x7,
		Length:     0x6,
		TypeFlags:  0x9,
		StripeLen:  0x8,
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
	if !reflect.DeepEqual(chunk, x) {
		t.Errorf("chunk: values not correct \nchunk=%+v\nx=%+v\n", chunk, x)
	}

}
