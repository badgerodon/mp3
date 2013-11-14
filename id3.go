package mp3

type (
	ID3V2Header struct {
		Version struct {
			Major, Revision byte
		}
		Flags byte
		Size  uint
	}
)

func (this *ID3V2Header) Parse(bs []byte) bool {
	if len(bs) < 10 {
		return false
	}
	if bs[0] != 'I' || bs[1] != 'D' || bs[2] != '3' {
		return false
	}

	this.Version.Major = bs[3]
	this.Version.Revision = bs[4]
	this.Flags = bs[5]
	if bs[6] >= 0x80 || bs[7] >= 0x80 || bs[8] >= 0x80 || bs[9] >= 0x80 {
		return false
	}

	this.Size = uint(bs[9]) + (uint(bs[8]) << 7) + (uint(bs[7]) << 14) + (uint(bs[6]) << 21)

	return true
}
