package audio

import (
	"fmt"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// AudioFileProperties holds information about the audio file.
type AudioFileProperties struct {
	QuantizationPeriod float64
	FileName           string
	ChannelCount       int
	Depth              uint16
	SampleRate         uint32
}

// ProcessWAV reads a WAV file and processes the samples with a callback function.
func ProcessWAV(filePath string, callback func(startIndex int, samples []int) error) (*AudioFileProperties, error) {
	// Open the WAV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file")
	}

	props := &AudioFileProperties{
		QuantizationPeriod: 1.0 / float64(decoder.SampleRate),
		FileName:           filePath,
		ChannelCount:       int(decoder.NumChans),
		Depth:              decoder.BitDepth,
		SampleRate:         uint32(decoder.SampleRate),
	}

	// Prepare a buffer for samples
	bufSize := 4096
	buf := &audio.IntBuffer{Data: make([]int, bufSize), Format: &audio.Format{
		NumChannels: int(decoder.NumChans),
		SampleRate:  int(decoder.SampleRate),
	}}

	startIndex := 0
	var n int

	for {
		n, err = decoder.PCMBuffer(buf)
		if err != nil {
			return nil, fmt.Errorf("failed to read samples: %w", err)
		}
		if n == 0 {
			break // No more samples to read
		}

		// Trim the buffer to the actual number of samples read
		if n < len(buf.Data) {
			buf.Data = buf.Data[:n]
		}

		// Process the samples with the callback function
		if err := callback(startIndex, buf.Data); err != nil {
			return nil, fmt.Errorf("callback error: %w", err)
		}
		startIndex += n
	}

	return props, nil
}
