package main

import (
	"fmt"
	"os"
	"sync"
	"wav_visualizer/pcm_energy_visualizer/audio"
	"wav_visualizer/pcm_energy_visualizer/chart"
	"wav_visualizer/pcm_energy_visualizer/filePicker"

	"github.com/NimbleMarkets/ntcharts/canvas"
)

func main() {
	audioFilePath := filePicker.PickWavFile()

	file, err := os.Open(audioFilePath)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	ap, err := audio.StartAudioProcessing(file, audioFilePath)
	if err != nil {
		fmt.Printf("failed to decode .wav: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ap.ProcessAudioFile()
	}()

	processedPoints := make([][]canvas.Float64Point, ap.FileProperties.ChannelCount)
	for i := range processedPoints {
		processedPoints[i] = make([]canvas.Float64Point, 2000)
	}

	wg.Wait()
	for channeledPoints := range ap.PointsChannel {
		for channelId, points := range channeledPoints {
			for _, point := range points {

				if point.X > ap.FileProperties.Duration.Seconds() {
					break
				}

				processedPoints[channelId] = append(processedPoints[channelId], point)
			}
		}
	}

	examples := chart.LineExamples{}
	chart.LineExamples.Examples(examples)
}
