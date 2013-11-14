package mp3

import (
	"testing"
)

func TestID3(t *testing.T) {
	valid := [][]byte{
		{'I', 'D', '3', 0, 0, 0, 0, 0, 0, 0, 0},
	}
	invalid := [][]byte{
		nil,
		{'I', 'D', '2', 0, 0, 0, 0, 0, 0, 0, 0},
	}

	var id3 ID3V2Header

	for _, test := range valid {
		if !id3.Parse(test) {
			t.Errorf("Expected valid ID3 tag for: %v", test)
		}
	}

	for _, test := range invalid {
		if id3.Parse(test) {
			t.Errorf("Expected invalid ID3 tag for: %v", test)
		}
	}
}
