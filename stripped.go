package mp3

import (
	"bytes"
	"fmt"
	"github.com/badgerodon/ioutil"
	"io"
)

func Stripped(src io.ReadSeeker) (io.ReadSeeker, error) {
	var pos int64
	var id3v2 ID3V2Header

	bs := make([]byte, 8192)

	_, err := src.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	_, err = io.ReadAtLeast(src, bs, 10)
	if err != nil {
		return nil, fmt.Errorf("Failed to read ID3: %v", err)
	}

	if id3v2.Parse(bs[:10]) {
		pos, err = src.Seek(int64(id3v2.Size), 0)
		if err != nil {
			return nil, fmt.Errorf("Failed to skip ID3: %v", err)
		}
	}

	var hdr FrameHeader
	for {
		_, err = io.ReadAtLeast(src, bs, 4)
		if err != nil {
			return nil, fmt.Errorf("Failed to read frame: %v", err)
		}

		if hdr.Parse(bs[:4]) {
			break
		}

		// Move one byte forward
		pos, err = src.Seek(1, 1)
		if err != nil {
			return nil, fmt.Errorf("Failed to seek: %v", err)
		}
	}

	_, err = io.ReadAtLeast(src, bs, hdr.Size())
	if err != nil {
		return nil, fmt.Errorf("Failed to read initial frame: %v, size: %v", err, hdr.Size())
	}

	// Skip the whole VBR header if we can
	var xing XingHeader
	if xing.Parse(bs[:hdr.Size()]) {
		pos, err = src.Seek(pos+int64(hdr.Size()), 0)
		if err != nil {
			return nil, err
		}
	}

	// Skip any ending ID3 tag
	end, err := src.Seek(-128, 2)
	if err != nil {
		return nil, err
	}

	_, err = io.ReadAtLeast(src, bs, 3)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bs[:3], []byte("TAG")) {
		end += 128
	}

	return ioutil.NewSectionReader(src, pos, end-pos), nil
}
