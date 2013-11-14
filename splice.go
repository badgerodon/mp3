package mp3

import (
	"github.com/badgerodon/ioutil"
	"io"
	"sort"
	"time"
)

type (
	spliceEntry struct {
		nanoseconds int64
		reader      io.Reader
	}
	splicer struct {
		reader      io.Reader
		splice      []spliceEntry
		nanoseconds int64
	}
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

func Splice(src io.ReadSeeker, splice map[time.Duration]io.ReadSeeker) (io.ReadSeeker, error) {
	times := []time.Duration{}
	for k, _ := range splice {
		times = append(times, k)
	}
	sort.Sort(durationSorter(times))

	pieces, err := Slice(src, times...)
	if err != nil {
		return nil, err
	}

	return ioutil.MultiReadSeeker(pieces...), nil
}
