package audio

import (
	"math"
)

func ConverSignalToEnergy(channeledSignal [][]AudioInputPoint, sampleRateMs float64, frameDuration float64) [][]AudioInputPoint {
	channels := make([][]AudioInputPoint, len(channeledSignal))

	for i, channel := range channeledSignal {
		channels[i] = convertChannelToEnergy(channel, sampleRateMs, frameDuration)
	}

	return channels
}

func convertChannelToEnergy(signalArray []AudioInputPoint, sampleRateMs float64, frameDuration float64) []AudioInputPoint {
	sampleCount := len(signalArray)
	samplesPerFrame := int(frameDuration / sampleRateMs)
	overlapStep := samplesPerFrame / 2

	energySignals := make([]AudioInputPoint, 0)

	for i := 0; i < sampleCount; i += overlapStep {
		lastIndex := min(sampleCount-1, i+samplesPerFrame)
		frameSamples := signalArray[i:lastIndex]

		frameEnergy := 0.0
		for _, sample := range frameSamples {
			frameEnergy += math.Pow(sample.Y, 2.0)
		}

		// assign energy to timeline - middle of the frame
		frameMiddleIndex := i + ((lastIndex - i) / 2)
		frameMiddleMs := sampleRateMs * float64(frameMiddleIndex)
		energySignals = append(energySignals, AudioInputPoint{X: frameMiddleMs, Y: frameEnergy})
	}

	return energySignals
}
