package leindex

import (
	"encoding/binary"
)

func IndexLE32(buf []byte, min, max uint32) int {
	// find a given 32bits little endian value in buf based on a range
	// useful for finding timestamps that may hint to other things

	minB := make([]byte, 4)
	maxB := make([]byte, 4)
	binary.LittleEndian.PutUint32(minB, min)
	binary.LittleEndian.PutUint32(maxB, max)

	return IndexLEbin(buf, minB, maxB)
}

func IndexLEbin(buf, min, max []byte) int {
	// find a little endian value in buf between min and max (we assume len(min)==len(max))
	itemln := len(min)
	bufln := len(buf)

	revmin := make([]byte, itemln)
	revmax := make([]byte, itemln)
	for i := 0; i < itemln; i++ {
		revmin[itemln-i-1] = min[i]
		revmax[itemln-i-1] = max[i]
	}

main:
	for i := itemln - 1; i < bufln; i++ {
		state := 0
		for j := 0; j < itemln; j++ {
			t := buf[i-j]
			//log.Printf("[%d] t=%d state=%d", j, t, state)
			switch state {
			case 0:
				// initial
				if t < revmin[j] || t > revmax[j] {
					continue main
				}
				if t > revmin[j] && t < revmax[j] {
					return i - itemln + 2
				}
				if t == revmin[j] && t == revmax[j] {
					// stay in state 0
					state = 0
				} else if t == revmin[j] {
					state = -1
				} else {
					state = 1
				}
			case 1:
				// previous digit was maximum, need to check if lower than max or equal
				if t > revmax[j] {
					continue main
				}
				if t < revmax[j] {
					return i - itemln + 2
				}
			case -1:
				// previous digit was minimum, need to check if higher than min or equal
				if t < revmin[j] {
					continue main
				}
				if t > revmin[j] {
					return i - itemln + 2
				}
			}
		}
		return i - itemln + 1
	}
	return -1
}
