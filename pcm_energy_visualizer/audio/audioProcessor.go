package audio

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

type AudioInputPoint struct {
	X float64
	Y float64
}

type AudioFileProperties struct {
	QuantizationPeriod float64
	FileName           string
	ChannelCount       int
	Depth              uint16
	SampleRate         uint32
	Duration           time.Duration
}

type AudioProcessor struct {
	FileProperties   *AudioFileProperties
	PointsChannel    chan [][]AudioInputPoint
	buf              *audio.IntBuffer
	decoder          *wav.Decoder
	maximumAmplitude float64
}

func calculateBufferSize(bitDepth uint16, channelCount int, sampleCount int) int {
	// Determine bytes per sample based on bit depth
	var bytesPerSample int
	switch bitDepth {
	case 16:
		bytesPerSample = 2 // 16 bits = 2 bytes
	case 24:
		bytesPerSample = 3 // 24 bits = 3 bytes
	case 32:
		bytesPerSample = 4 // 32 bits = 4 bytes
	default:
		panic("Unsupported bit depth")
	}

	totalBytes := sampleCount * bytesPerSample * channelCount

	return totalBytes
}

func getMaxAmplitude(bitDepth uint16) float64 {
	switch bitDepth {
	case 16:
		return float64(math.MaxUint16)
	case 24:
		return 16777215
	case 32:
		return float64(math.MaxUint32)
	default:
		panic("Unsupported bit depth")
	}
}

// ======= Constructor
func StartAudioProcessing(file *os.File, filePath string) (*AudioProcessor, error) {
	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("invalid WAV file")
	}

	duration, err := decoder.Duration()
	if err != nil {
		return nil, fmt.Errorf("encountered error while getting duration")
	}

	props := &AudioFileProperties{
		QuantizationPeriod: 1.0 / float64(decoder.SampleRate),
		FileName:           filePath,
		ChannelCount:       int(decoder.NumChans),
		Depth:              decoder.BitDepth,
		SampleRate:         uint32(decoder.SampleRate),
		Duration:           duration,
	}

	sampleCount := 1024
	bufferSize := calculateBufferSize(props.Depth, props.ChannelCount, sampleCount)

	buf := &audio.IntBuffer{Data: make([]int, bufferSize), Format: &audio.Format{
		NumChannels: int(decoder.NumChans),
		SampleRate:  int(decoder.SampleRate),
	}}

	ap := &AudioProcessor{
		FileProperties:   props,
		buf:              buf,
		decoder:          decoder,
		PointsChannel:    make(chan [][]AudioInputPoint, props.ChannelCount),
		maximumAmplitude: getMaxAmplitude(props.Depth),
	}

	return ap, nil
}

// ========== SAMPLE PARSING
func (ap *AudioProcessor) parseSamples(startIndex int, samples []int) {
	points := make([][]AudioInputPoint, ap.FileProperties.ChannelCount)

	for i := range points {
		points[i] = make([]AudioInputPoint, len(samples)/int(ap.FileProperties.ChannelCount))
	}

	for i, sample := range samples {
		channel := i % int(ap.FileProperties.ChannelCount)
		sampleIndex := i / int(ap.FileProperties.ChannelCount)
		if sampleIndex >= len(points[channel]) {
			return
		}
		amplitude := float64(sample) / ap.maximumAmplitude
		timeInSeconds := float64(startIndex+sampleIndex) / float64(ap.FileProperties.SampleRate)
		points[channel][sampleIndex] = AudioInputPoint{X: timeInSeconds, Y: amplitude}
	}

	ap.PointsChannel <- points
}

func (ap *AudioProcessor) ProcessAudioFile() {
	startIndex := 0
	var wg sync.WaitGroup

	for {
		n, err := ap.decoder.PCMBuffer(ap.buf)
		if err != nil {
			return
		}
		if n == 0 {
			break // No more samples to read
		}

		if n < len(ap.buf.Data) {
			ap.buf.Data = ap.buf.Data[:n]
		}
		wg.Add(1)
		go func(startIndex int, samples []int) {
			defer wg.Done()
			ap.parseSamples(startIndex, samples)
		}(startIndex, ap.buf.Data)

		startIndex += n
	}

	go func() {
		wg.Wait()
		close(ap.PointsChannel)
	}()
}
