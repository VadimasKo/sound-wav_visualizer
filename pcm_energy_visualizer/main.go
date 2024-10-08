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
var frameDuration = 20.0 / 1000.0 // 20ms
var energyThreshold = 0.5

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

	timeChartOptions := chart.ChartOptions{
		Title:           "Time series chart",
		ChanneledPoints: processedPoints,
		Segments:        make([][]audio.AudioSegment, 0),
	}

	// ====== ZRC Chart
	zeroCrossingRate := audio.ConvertSignalToZCRGraph(
		processedPoints,
		ap.FileProperties.QuantizationPeriod,
		frameDuration,
	)

	zeroCrosingRateChartOptions := chart.ChartOptions{
		Title:           "Zero crossing rate chart",
		ChanneledPoints: zeroCrossingRate,
		Segments:        make([][]audio.AudioSegment, 0),
	}

	// ====== energy chart
	energyPoints := audio.ConverSignalToEnergy(
		processedPoints,
		ap.FileProperties.QuantizationPeriod,
		frameDuration,
	)

	energyChartOptions := chart.ChartOptions{
		Title:           "Energy chart",
		ChanneledPoints: energyPoints,
		Segments:        make([][]audio.AudioSegment, 0),
	}

	// ====== energy chart with segments
	energySegments := audio.GetSignalThresholdSegments(energyPoints, energyThreshold)
	segmentedEnergyChartOptions := chart.ChartOptions{
		Title:           "Segmented Energy Chart",
		ChanneledPoints: energyPoints,
		Segments:        energySegments,
	}

	page := components.NewPage()
	page.AddCharts(
		timeChartOptions.CreateAudioLineChart(),
		energyChartOptions.CreateAudioLineChart(),
		zeroCrosingRateChartOptions.CreateAudioLineChart(),
		segmentedEnergyChartOptions.CreateAudioLineChart(),
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
