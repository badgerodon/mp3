package mp3

import (
	"bytes"
	"encoding/binary"
)

type XingHeader struct {
	Frames, Bytes, Quality int
}

// Parse an Xing header from the first frame of an mp3
func (this *XingHeader) Parse(src []byte) bool {
	if len(src) == 0 {
		return false
	}

	// parse header id
	idx := bytes.Index(src, []byte("Xing"))
	if idx < 0 {
		idx = bytes.Index(src, []byte("Info"))
	}
	if idx < 0 {
		return false
	}

	src = src[idx+4:]

	flags := binary.BigEndian.Uint32(src[:4])
	src = src[4:]
	// Frames
	if flags&1 == 1 {
		this.Frames = int(binary.BigEndian.Uint32(src[:4]))
		src = src[4:]
	}
	// Bytes
	if flags&2 == 2 {
		this.Bytes = int(binary.BigEndian.Uint32(src[:4]))
		src = src[4:]
	}
	// TOC
	if flags&4 == 4 {
		src = src[4:]
	}
	// Quality
	if flags&8 == 8 {
		this.Quality = int(binary.BigEndian.Uint32(src[:4]))
	}

	return true
}
