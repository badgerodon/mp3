package mp3

import (
	"os"
	"testing"
	"time"
)

func TestLength(t *testing.T) {
	f, err := os.Open("1000hz.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	duration, err := Length(f)
	if err != nil {
		t.Fatal(err)
	}
	if duration != 5067754912 {
		t.Errorf("Expected length to return %v, got %v", time.Duration(5067754912), duration)
	}
}
