package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"wav_visualizer/pcm_energy_visualizer/audio"
	"wav_visualizer/pcm_energy_visualizer/chart"
	"wav_visualizer/pcm_energy_visualizer/filePicker"

	"github.com/go-echarts/go-echarts/v2/components"
)

var outputPath = "./output"

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

	processedPoints := make([][]audio.AudioInputPoint, ap.FileProperties.ChannelCount)
	for i := range processedPoints {
		processedPoints[i] = make([]audio.AudioInputPoint, 2000)
	}

	wg.Wait()
	for channeledPoints := range ap.PointsChannel {
		for channelId, points := range channeledPoints {
			for _, point := range points {
				processedPoints[channelId] = append(processedPoints[channelId], point)
			}
		}
	}

	// audio energy.go
	// audio energy.go
	// audo  pcm.go

	// create page file w header showing file name and information
	// create time graph
	// create energy graph
	// create pcm graph
	// create segmented energy graph
	// write page to html file

	page := components.NewPage()
	page.AddCharts(
		chart.AudioLineChart("Time series graph", processedPoints),
	)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(outputPath + "/chart.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
