package mp3

import (
	"bufio"
	"io"
)

type (
	Frames struct {
		src    *bufio.Reader
		offset int64
		header FrameHeader
		err    error
	}
)

func GetFrames(src io.ReadSeeker) (*Frames, error) {
	stripped, err := Stripped(src)
	if err != nil {
		return nil, err
	}
	return &Frames{
		src: bufio.NewReader(stripped),
	}, nil
}

func (this *Frames) Next() bool {
	for {
		bs, err := this.src.Peek(4)
		if err != nil {
			this.err = err
			return false
		}
		if this.header.Parse(bs) {
			for i := 0; i < this.header.Size(); i++ {
				_, err = this.src.ReadByte()
				if err != nil {
					this.err = err
					return false
				}
			}
			this.offset += int64(this.header.Size())
			return true
		} else {
			this.src.ReadByte()
		}
	}
}

func (this *Frames) Header() FrameHeader {
	return this.header
}

func (this *Frames) Offset() int64 {
	return this.offset
}

func (this *Frames) Error() error {
	if this.err == io.EOF {
		return nil
	}
	return this.err
}
