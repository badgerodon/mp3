package mp3

import (
	"github.com/badgerodon/ioutil"
	"os"
	"testing"
	"time"
)

func TestSlice(t *testing.T) {
	nearlyEqual := func(a, b time.Duration) bool {
		diff := a - b
		if diff < 0 {
			diff = -diff
		}
		return diff < 100*time.Millisecond
	}

	f, err := os.Open("440hz.mp3")
	if err != nil {
		t.Fatalf("Error opening sample file: %v", err)
	}
	defer f.Close()

	slices, err := Slice(f, time.Second*1, time.Second*2)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if len(slices) != 3 {
		t.Fatalf("Expected 3 slices, got: %v", err)
	}
	len1, err := Length(slices[0])
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !nearlyEqual(len1, time.Second) {
		t.Errorf("Expected %v to be %v, got %v", slices[0], time.Second, len1)
	}
	len2, err := Length(slices[1])
	if err != nil {
		t.Errorf("Expected no error for %v, got: %v", slices[1], err)
	}
	if !nearlyEqual(len2, time.Second) {
		t.Errorf("Expected %v to be %v, got %v", slices[1], time.Second, len2)
	}
	len3, err := Length(slices[2])
	if err != nil {
		t.Errorf("Expected no error for %v, got: %v", slices[2], err)
	}
	if !nearlyEqual(len3, 3*time.Second) {
		t.Errorf("Expected the third slice to be %v, got %v", 3*time.Second, len3)
	}

}

func TestSlicedAndConcatenated(t *testing.T) {
	f, err := os.Open("440hz.mp3")
	if err != nil {
		t.Fatalf("Error opening sample file: %v", err)
	}
	defer f.Close()

	slices, err := Slice(f, time.Second, 2*time.Second, 3*time.Second)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	concatenated := ioutil.NewMultiReadSeeker(slices...)

	l1, err := Length(f)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	l2, err := Length(concatenated)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if l1 != l2 {
		t.Errorf("Expected concatenated (%v) to be the same length as original (%v)", l2, l1)
	}

}
