package leindex

import "encoding/binary"

func IndexBE32(buf []byte, min, max uint32) int {
	// find a given 32bits big endian value in buf based on a range
	// useful for finding timestamps that may hint to other things

	minB := make([]byte, 4)
	maxB := make([]byte, 4)
	binary.BigEndian.PutUint32(minB, min)
	binary.BigEndian.PutUint32(maxB, max)

	return IndexRange(buf, minB, maxB)
}

func IndexRange(buf, min, max []byte) int {
	// find a big endian value in buf between min and max (we assume len(min)==len(max))
	itemln := len(min)
	bufln := len(buf)

	// return -1 if buffer is smaller than an item
	if bufln < itemln {
		return -1
	}

main:
	for i := 0; i < bufln-itemln; i++ {
		state := 0
		for j := 0; j < itemln; j++ {
			t := buf[i+j]
			//log.Printf("[%d] i=%d t=0x%x state=%d", j, i, t, state)
			switch state {
			case 0:
				// initial
				if t < min[j] || t > max[j] {
					continue main
				}
				if t > min[j] && t < max[j] {
					return i
				}
				if t == min[j] && t == max[j] {
					// stay in state 0
					state = 0
				} else if t == min[j] {
					state = -1
				} else {
					state = 1
				}
			case 1:
				// previous digit was maximum, need to check if lower than max or equal
				if t > max[j] {
					continue main
				}
				if t < max[j] {
					return i
				}
			case -1:
				// previous digit was minimum, need to check if higher than min or equal
				if t < min[j] {
					continue main
				}
				if t > min[j] {
					return i
				}
			}
		}
		return i
	}
	return -1
}
