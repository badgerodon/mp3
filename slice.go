package mp3

import (
	"fmt"
	"github.com/badgerodon/ioutil"
	"io"
	"time"
)

func Slice(src io.ReadSeeker, cutPoints ...time.Duration) ([]io.ReadSeeker, error) {
	pieces := make([]io.ReadSeeker, 0, len(cutPoints)+1)

	_, err := src.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to seek src to beginning: %v", err)
	}

	frames, err := GetFrames(src)
	if err != nil {
		return nil, err
	}

	length, err := src.Seek(0, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to seek src to end: %v", err)
	}

	var elapsed time.Duration
	var lastOffset int64
	for frames.Next() {
		if len(cutPoints) == 0 {
			break
		}

		elapsed += frames.Header().Duration()

		if cutPoints[0] <= elapsed {
			piece := ioutil.NewSectionReader(src, lastOffset, frames.Offset()-lastOffset)
			pieces = append(pieces, piece)
			lastOffset = frames.Offset()
			cutPoints = cutPoints[1:]
		}
	}

	pieces = append(pieces, ioutil.NewSectionReader(src, lastOffset, length-lastOffset))

	return pieces, frames.Error()
}
