package leindex_test

import (
	"encoding/binary"
	"testing"

	"github.com/KarpelesLab/leindex"
)

func TestBE(t *testing.T) {
	buf := make([]byte, 4096)

	binary.BigEndian.PutUint32(buf[:4], 0x98765432)
	binary.BigEndian.PutUint32(buf[52:56], 1546300799)
	binary.BigEndian.PutUint32(buf[56:60], 1672531201)
	binary.BigEndian.PutUint32(buf[60:64], 1650513840)
	binary.BigEndian.PutUint32(buf[90:94], 0xffffffff)
	binary.BigEndian.PutUint32(buf[30:34], 0x12345678)
	binary.BigEndian.PutUint32(buf[3030:3034], 123456789)
	buf[8] = 0x1f

	tests := []testVec{
		{0x98765000, 0x98766000, 0},
		{1546300800, 1672531200, 60},
		{123456789, 123456790, 3030},
		{1000, 1000, -1},
		{0x0000ffff, 0x0000ffff, 88},
		{0x00ffff00, 0x00ffff00, -1},
		{0xffff0000, 0xffff0000, 92},
		{0xff000000, 0xff000000, 93},
	}

	for _, v := range tests {
		res := leindex.IndexBE32(buf, v.min, v.max)
		if res != v.expect {
			if res == -1 {
				t.Errorf("could not find expected value %x~%x at %d", v.min, v.max, v.expect)
			} else {
				n := binary.BigEndian.Uint32(buf[res : res+4])
				t.Errorf("expected %d but got %d (value=%x) while looking for %x~%x", v.expect, res, n, v.min, v.max)
			}
		}
	}
}
