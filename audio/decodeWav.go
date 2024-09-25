package audio

import (
	"fmt"
	"os"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type AudioFileProperties struct {
	QuantizationPeriod float64
	FileName           string
	ChannelCount       int
	Depth              uint16
	SampleRate         uint32
	Duration           time.Duration
}

func DecodeWav(file *os.File) (*audio.IntBuffer, *AudioFileProperties, error) {
	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, nil, fmt.Errorf("invalid WAV file")
	}

	duration, err := decoder.Duration()
	if err != nil {
		return nil, nil, fmt.Errorf("encountered error while getting duration")
	}

	props := &AudioFileProperties{
		QuantizationPeriod: 1.0 / float64(decoder.SampleRate),
		FileName:           "TEST",
		ChannelCount:       int(decoder.NumChans),
		Depth:              decoder.BitDepth,
		SampleRate:         uint32(decoder.SampleRate),
		Duration:           duration,
	}

	// Prepare a buffer for samples
	bufSize := 4096
	buf := &audio.IntBuffer{Data: make([]int, bufSize), Format: &audio.Format{
		NumChannels: int(decoder.NumChans),
		SampleRate:  int(decoder.SampleRate),
	}}

	return buf, props, nil
}
