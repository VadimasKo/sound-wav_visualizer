package audio

import (
	"math"
)

func ConvertSignalToZCRGraph(channeledSignal [][]AudioInputPoint, sampleRateMs float64, frameDuration float64) [][]AudioInputPoint {
	channels := make([][]AudioInputPoint, len(channeledSignal))

	for i, channel := range channeledSignal {
		channels[i] = convertChannelToZCRGraph(channel, sampleRateMs, frameDuration)
	}

	return channels
}

func convertChannelToZCRGraph(signalArray []AudioInputPoint, sampleRateMs float64, frameDuration float64) []AudioInputPoint {
	sampleCount := len(signalArray)
	samplesPerFrame := int(frameDuration / sampleRateMs)
	overlapStep := samplesPerFrame / 2

	crossingSignals := make([]AudioInputPoint, 0)

	for i := 0; i < sampleCount; i += overlapStep {
		lastIndex := min(sampleCount-1, i+samplesPerFrame)
		frameSamples := signalArray[i:lastIndex]

		if len(frameSamples) < 2 {
			continue
		}

		zcr := 0.0
		for sampleIndex := 1; sampleIndex < len(frameSamples); sampleIndex++ {
			signCurrent := math.Copysign(1, frameSamples[sampleIndex].Y)
			signPrev := math.Copysign(1, frameSamples[sampleIndex-1].Y)

			dif := math.Abs(signCurrent - signPrev)
			zcr += dif
		}

		frameZcr := zcr / (2.0 * float64(len(frameSamples)))

		// assign energy to timeline - middle of the frame
		frameMiddleIndex := i + ((lastIndex - i) / 2)
		frameMiddleMs := sampleRateMs * float64(frameMiddleIndex)
		crossingSignals = append(crossingSignals, AudioInputPoint{X: frameMiddleMs, Y: frameZcr})
	}

	return crossingSignals
}
