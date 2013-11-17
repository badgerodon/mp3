package mp3

import (
	"bytes"
	"fmt"
	"github.com/badgerodon/ioutil"
	"io"
)

func getFirstFrameOffset(src io.ReadSeeker) (int64, error) {
	var id3v2 ID3V2Header
	_, err := src.Seek(0, 0)
	if err != nil {
		return 0, err
	}

	bs := make([]byte, 10)

	_, err = io.ReadAtLeast(src, bs, 10)
	if err != nil {
		return 0, fmt.Errorf("Failed to read first 10 bytes: %v", err)
	}

	err = id3v2.Parse(bs)
	if err == nil {
		return id3v2.Size, nil
	}

	return 0, nil
}

// Skips both the ID3V2 tags and optional VBR headers
func getFirstRealFrameOffset(src io.ReadSeeker) (int64, error) {
	var hdr FrameHeader
	var xing XingHeader

	off, err := getFirstFrameOffset(src)
	if err != nil {
		return 0, err
	}

	_, err = src.Seek(off, 0)
	if err != nil {
		return 0, err
	}

	bs := make([]byte, 8192)

	_, err = io.ReadAtLeast(src, bs, 4)
	if err != nil {
		return 0, err
	}

	err = hdr.Parse(bs)
	if err != nil {
		return 0, err
	}

	if xing.Parse(bs[:int(hdr.Size())]) {
		return off + hdr.Size(), nil
	}

	return off, nil
}

func getLastFrameEnd(src io.ReadSeeker) (int64, error) {
	end, err := src.Seek(-128, 2)
	if err != nil {
		return 0, err
	}

	bs := make([]byte, 3)
	_, err = io.ReadAtLeast(src, bs, 3)
	if err != nil {
		return 0, err
	}
	if !bytes.Equal(bs, []byte("TAG")) {
		end += 128
	}
	return end, nil
}

func Stripped(src io.ReadSeeker) (*ioutil.SectionReader, error) {
	o1, err := getFirstRealFrameOffset(src)
	if err != nil {
		return nil, err
	}
	o2, err := getLastFrameEnd(src)
	if err != nil {
		return nil, err
	}
	return ioutil.NewSectionReader(src, o1, o2-o1), nil
}
