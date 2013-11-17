package mp3

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestStripped(t *testing.T) {
	f, err := os.Open("440hz.mp3")
	if err != nil {
		t.Fatalf("Error opening sample: %v", err)
	}
	defer f.Close()

	off, err := getFirstFrameOffset(f)
	if err != nil {
		t.Error(err)
	}

	if off != 33 {
		t.Errorf("Expected offset to be %v, got %v", 33, off)
	}

	off, err = getFirstRealFrameOffset(f)
	if err != nil {
		t.Error(err)
	}

	if off != 33 {
		t.Errorf("Expected offset to be %v, got %v", 33, off)
	}

	end, err := getLastFrameEnd(f)
	if err != nil {
		t.Error(err)
	}

	if end != 30439 {
		t.Errorf("Expected offset to be %v, got %v", 30439, end)
	}

	limited, err := Stripped(f)
	if err != nil {
		t.Error(err)
	}
	bs, err := ioutil.ReadAll(limited)
	if err != nil {
		t.Error(err)
	}
	if len(bs) != 30406 {
		t.Errorf("Expected %v bytes, got %v", 30406, len(bs))
	}

	limited, err = Stripped(limited)
	if err != nil {
		t.Error(err)
	}

}
