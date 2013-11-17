package mp3

import (
	"fmt"
	"github.com/badgerodon/ioutil"
	"io"
)

type (
	Frames struct {
		src           *ioutil.SectionReader
		offset, count int64
		header        FrameHeader
		err           error
	}
)

func GetFrames(src io.ReadSeeker) (*Frames, error) {
	stripped, err := Stripped(src)
	if err != nil {
		return nil, fmt.Errorf("failed to strip src: %v", err)
	}
	return &Frames{
		src:    stripped,
		offset: 0,
	}, nil
}

func (this *Frames) Next() bool {
	var err error
	if this.count > 0 {
		// skip to next frame
		this.offset, err = this.src.Seek(this.offset+this.header.Size(), 0)
		if err != nil {
			this.err = fmt.Errorf("premature end of frame: %v", err)
			return false
		}
	}
	bs := make([]byte, 4)
	_, err = io.ReadAtLeast(this.src, bs, 4)
	if err != nil {
		this.err = err
		return false
	}
	err = this.header.Parse(bs)
	if err != nil {
		this.err = err
		return false
	}
	this.count++
	return true
}

func (this *Frames) Header() *FrameHeader {
	return &this.header
}

func (this *Frames) Offset() int64 {
	return this.src.Offset() + this.offset
}

func (this *Frames) Error() error {
	if this.err == io.EOF {
		return nil
	}
	return this.err
}
