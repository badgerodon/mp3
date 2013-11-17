package mp3

import (
	"fmt"
	"time"
)

type (
	Version     byte
	Layer       byte
	ChannelMode byte
	Emphasis    byte
	FrameHeader struct {
		Version         Version
		Layer           Layer
		Protection      bool
		Bitrate         int
		SampleRate      int
		Pad             bool
		Private         bool
		ChannelMode     ChannelMode
		IntensityStereo bool
		MSStereo        bool
		CopyRight       bool
		Original        bool
		Emphasis        Emphasis

		Size     int64
		Samples  int
		Duration time.Duration
	}
)

const (
	MPEG25 Version = iota
	MPEGReserved
	MPEG2
	MPEG1
)

const (
	LayerReserved Layer = iota
	Layer3
	Layer2
	Layer1
)

const (
	EmphNone Emphasis = iota
	Emph5015
	EmphReserved
	EmphCCITJ17
)

const (
	Stereo ChannelMode = iota
	JointStereo
	DualChannel
	SingleChannel
)

var (
	bitrates = map[Version]map[Layer][15]int{
		MPEG1: { // MPEG 1
			Layer1: {0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448}, // Layer1
			Layer2: {0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384},    // Layer2
			Layer3: {0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320},     // Layer3
		},
		MPEG2: { // MPEG 2, 2.5
			Layer1: {0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256}, // Layer1
			Layer2: {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},      // Layer2
			Layer3: {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},      // Layer3
		},
	}
	sampleRates = map[Version][3]int{
		MPEG1:        {44100, 48000, 32000},
		MPEG2:        {22050, 24000, 16000},
		MPEG25:       {11025, 12000, 8000},
		MPEGReserved: {0, 0, 0},
	}
	samplesPerFrame = map[Version]map[Layer]int{
		MPEG1: {
			Layer1: 384,
			Layer2: 1152,
			Layer3: 1152,
		},
		MPEG2: {
			Layer1: 384,
			Layer2: 1152,
			Layer3: 576,
		},
	}
	slotSize = map[Layer]int{
		LayerReserved: 0,
		Layer3:        1,
		Layer2:        1,
		Layer1:        4,
	}
)

func init() {
	bitrates[MPEG25] = bitrates[MPEG2]
	samplesPerFrame[MPEG25] = samplesPerFrame[MPEG2]
}

func (this *FrameHeader) Parse(bs []byte) error {
	this.Size = 0
	this.Samples = 0
	this.Duration = 0

	if len(bs) < 4 {
		return fmt.Errorf("not enough bytes")
	}
	if bs[0] != 0xFF || (bs[1]&0xE0) != 0xE0 {
		return fmt.Errorf("missing sync word, got: %x, %x", bs[0], bs[1])
	}
	this.Version = Version((bs[1] >> 3) & 0x03)
	if this.Version == MPEGReserved {
		return fmt.Errorf("reserved mpeg version")
	}

	this.Layer = Layer(((bs[1] >> 1) & 0x03))
	if this.Layer == LayerReserved {
		return fmt.Errorf("reserved layer")
	}

	this.Protection = (bs[1] & 0x01) != 0x01

	bitrateIdx := (bs[2] >> 4) & 0x0F
	if bitrateIdx == 0x0F {
		return fmt.Errorf("invalid bitrate: %v", bitrateIdx)
	}
	this.Bitrate = bitrates[this.Version][this.Layer][bitrateIdx] * 1000
	if this.Bitrate == 0 {
		return fmt.Errorf("invalid bitrate: %v", bitrateIdx)
	}

	sampleRateIdx := (bs[2] >> 2) & 0x03
	if sampleRateIdx == 0x03 {
		return fmt.Errorf("invalid sample rate: %v", sampleRateIdx)
	}
	this.SampleRate = sampleRates[this.Version][sampleRateIdx]

	this.Pad = ((bs[2] >> 1) & 0x01) == 0x01

	this.Private = (bs[2] & 0x01) == 0x01

	this.ChannelMode = ChannelMode(bs[3]>>6) & 0x03

	// todo: mode extension

	this.CopyRight = (bs[3]>>3)&0x01 == 0x01

	this.Original = (bs[3]>>2)&0x01 == 0x01

	this.Emphasis = Emphasis(bs[3] & 0x03)
	if this.Emphasis == EmphReserved {
		return fmt.Errorf("reserved emphasis")
	}

	this.Size = this.size()
	this.Samples = this.samples()
	this.Duration = this.duration()

	return nil
}

func (this *FrameHeader) samples() int {
	return samplesPerFrame[this.Version][this.Layer]
}

func (this *FrameHeader) size() int64 {
	bps := float64(this.samples()) / 8
	fsize := (bps * float64(this.Bitrate)) / float64(this.SampleRate)
	if this.Pad {
		fsize += float64(slotSize[this.Layer])
	}
	return int64(fsize)
}

func (this *FrameHeader) duration() time.Duration {
	ms := (1000 / float64(this.SampleRate)) * float64(this.samples())
	return time.Duration(time.Duration(float64(time.Millisecond) * ms))
}
