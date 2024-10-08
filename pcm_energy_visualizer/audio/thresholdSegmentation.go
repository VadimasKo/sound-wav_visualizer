package audio

type AudioSegment struct {
	StartSampleIndex int
	EndSampleIndex   int
}

func GetSignalThresholdSegments(channeledSignal [][]AudioInputPoint, threshold float64) [][]AudioSegment {
	channelSegments := make([][]AudioSegment, len(channeledSignal))

	for i, channel := range channeledSignal {
		channelSegments[i] = getThresholdSegments(channel, threshold)
	}

	return channelSegments
}

func getThresholdSegments(audioSignal []AudioInputPoint, threshold float64) []AudioSegment {
	signalCount := len(audioSignal)
	audioSegments := make([]AudioSegment, 0)

	for startSearchIndex := 0; startSearchIndex < signalCount; startSearchIndex++ {
		endIndex := 0
		startIndex := -1

		if audioSignal[startSearchIndex].Y <= threshold {
			continue
		}
		startIndex = startSearchIndex

		for endSearchIndex := startSearchIndex + 1; endSearchIndex < signalCount; endSearchIndex++ {
			if audioSignal[endSearchIndex].Y < threshold || endSearchIndex == signalCount-1 {
				endIndex = endSearchIndex
				startSearchIndex = endIndex
				break
			}
		}

		if startIndex != -1 {
			newSegment := AudioSegment{StartSampleIndex: startIndex, EndSampleIndex: endIndex}
			audioSegments = append(audioSegments, newSegment)
		}
	}

	return audioSegments
}
