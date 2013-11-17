package mp3

import (
	"fmt"
)

type (
	ID3V2Header struct {
		Version struct {
			Major, Revision byte
		}
		Flags byte
		Size  int64
	}
)

func (this *ID3V2Header) Parse(bs []byte) error {
	if len(bs) < 10 {
		return fmt.Errorf("Expected at least 10 bytes, got: %v", len(bs))
	}
	if bs[0] != 'I' || bs[1] != 'D' || bs[2] != '3' {
		return fmt.Errorf("Expected ID3 head to start with ID3, got: %v", bs[:3])
	}

	this.Version.Major = bs[3]
	this.Version.Revision = bs[4]
	this.Flags = bs[5]
	if bs[6] >= 0x80 || bs[7] >= 0x80 || bs[8] >= 0x80 || bs[9] >= 0x80 {
		return fmt.Errorf("Invalid size, got: %v", bs[6:])
	}

	this.Size = int64(bs[9]) +
		int64(uint(bs[8])<<7) +
		int64(uint(bs[7])<<14) +
		int64(uint(bs[6])<<21) +
		10

	return nil
}
