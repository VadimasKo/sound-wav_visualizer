package main

import (
	"fmt"
	"os"
	"sync"
	"vadimasKo/wav_visualizer/audio"
	"vadimasKo/wav_visualizer/filePicker"
	"vadimasKo/wav_visualizer/visualizer"

	"github.com/NimbleMarkets/ntcharts/canvas"
	tea "github.com/charmbracelet/bubbletea"
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

	m := visualizer.WavelineModel(ap.FileProperties)

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
			fmt.Printf("Processing Channel %d, Total Points: %d\n", channelId, len(points))

			for _, point := range points {
				fmt.Printf("Channel: %d | Time (X): %.2f seconds | Amplitude (Y): %.2f\n", channelId, point.X, point.Y)

				if point.X > ap.FileProperties.Duration.Seconds() {
					fmt.Printf("Warning: Time value out of range for Channel %d: %.2f\n", channelId, point.X)
					break
				}

				if point.Y > 1 || point.Y < -1 { // Example of outlier detection
					fmt.Printf("Warning: Amplitude value out of range for Channel %d: %.2f\n", channelId, point.Y)
				}

				processedPoints[channelId] = append(processedPoints[channelId], point)
			}
		}
	}

	var a = 2 + 2
	println("%d", a)

	visualizer.PlotMultiChannelData(m, processedPoints)

	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		return
	}
}
