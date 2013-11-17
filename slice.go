package mp3

import (
	//"fmt"
	"github.com/badgerodon/ioutil"
	"io"
	"time"
)

func Slice(src io.ReadSeeker, cutPoints ...time.Duration) ([]io.ReadSeeker, error) {
	pieces := make([]io.ReadSeeker, 0, len(cutPoints)+1)

	stripped, err := Stripped(src)
	if err != nil {
		return nil, err
	}

	start := stripped.Offset()
	end := stripped.Offset() + stripped.Length()

	frames, err := GetFrames(src)
	if err != nil {
		return nil, err
	}

	var elapsed time.Duration
	lastOffset := start
	for frames.Next() {
		if len(cutPoints) == 0 {
			break
		}

		elapsed += frames.Header().Duration

		if cutPoints[0] <= elapsed {
			piece := ioutil.NewSectionReader(src, lastOffset, frames.Offset()+frames.Header().Size-lastOffset)
			pieces = append(pieces, piece)
			lastOffset = frames.Offset() + frames.Header().Size
			cutPoints = cutPoints[1:]
		}
	}

	pieces = append(pieces, ioutil.NewSectionReader(src, lastOffset, end-lastOffset))

	return pieces, frames.Error()
}
