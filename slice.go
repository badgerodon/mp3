package mp3

import (
	"github.com/badgerodon/ioutil"
	"io"
	"time"
)

func Slice(src io.ReadSeeker, cutPoints ...time.Duration) ([]io.ReadSeeker, error) {
	pieces := make([]io.ReadSeeker, 0, len(cutPoints)+1)
	frames, err := GetFrames(src)
	if err != nil {
		return nil, err
	}

	length, err := src.Seek(0, 2)
	if err != nil {
		return nil, err
	}

	var elapsed time.Duration
	var lastOffset int64
	for frames.Next() {
		if len(cutPoints) == 0 {
			pieces = append(pieces, ioutil.NewSectionReader(src, frames.Offset(), length-frames.Offset()))
			break
		}

		hdr := frames.Header()
		elapsed += hdr.Duration()

		if cutPoints[0] <= elapsed {
			pieces = append(pieces, ioutil.NewSectionReader(src, lastOffset, frames.Offset()))
			lastOffset = frames.Offset()
			cutPoints = cutPoints[1:]
		}
	}

	return pieces, frames.Error()
}
