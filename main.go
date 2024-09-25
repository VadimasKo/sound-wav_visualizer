package main

import (
	"fmt"
	"os"
	"vadimasKo/wav_visualizer/audio"
	"vadimasKo/wav_visualizer/filePicker"
	// "vadimasKo/wav_visualizer/visualizer"
)

func main() {
	audioFilePath := filePicker.PickWavFile()

	file, err := os.Open(audioFilePath)
	if err != nil || file == nil {
		fmt.Errorf("failed to open file: %w", err)
		return
	}
	defer file.Close()

	_, audioProps, err := audio.DecodeWav(file)
	if err != nil {
		fmt.Errorf("failed to decode .wav: %w", err)
		return
	}
	// visualizer.CreateWavelineChart(audioFilePath)
	fmt.Println("FileName:", audioProps.FileName)
	fmt.Println("QuantizationPeriod:", audioProps.QuantizationPeriod)
	fmt.Println("ChannelCount:", audioProps.ChannelCount)
	fmt.Println("Depth:", audioProps.Depth)
	fmt.Println("SampleRate:", audioProps.SampleRate)

	fmt.Println("Duration:", audioProps.Duration.Seconds())
}
