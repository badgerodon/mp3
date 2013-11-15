package mp3

import (
	"github.com/badgerodon/ioutil"
	"io"
	"sort"
	"time"
)

type durationSorter []time.Duration

func (this durationSorter) Len() int {
	return len(this)
}
func (this durationSorter) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func (this durationSorter) Less(i, j int) bool {
	return this[i] < this[j]
}

// Take a source MP3 and insert all the splice members into it (at the specified durations)
func Splice(src io.ReadSeeker, splice map[time.Duration]io.ReadSeeker) (io.ReadSeeker, error) {
	// Get the times
	spliceTimes := []time.Duration{}
	for k, _ := range splice {
		spliceTimes = append(spliceTimes, k)
	}
	sort.Sort(durationSorter(spliceTimes))

	// Slice up the src into len(splice)+1 pieces
	sliced, err := Slice(src, spliceTimes...)
	if err != nil {
		return nil, err
	}

	// Insert splice members between the slices
	pieces := []io.ReadSeeker{sliced[0]}
	for i := 1; i < len(sliced); i++ {
		pieces = append(pieces, splice[spliceTimes[i-1]], pieces[i])
	}

	// Treat all the pieces as one big ReadSeeker
	return ioutil.NewMultiReadSeeker(pieces...), nil
}
