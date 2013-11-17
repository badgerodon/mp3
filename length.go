package mp3

import (
	"fmt"
	"io"
	"time"
)

func Length(src io.ReadSeeker) (time.Duration, error) {
	frames, err := GetFrames(src)
	if err != nil {
		return 0, fmt.Errorf("failed to parse frames: %v", err)
	}

	duration := time.Duration(0)

	for frames.Next() {
		duration += frames.Header().Duration()
	}

	return duration, frames.Error()
}
