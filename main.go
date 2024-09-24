package main

import (
	"vadimasKo/wav_visualizer/audio"
	"vadimasKo/wav_visualizer/filePicker"
	// "vadimasKo/wav_visualizer/visualizer"
)

func main() {
	audioFile := filePicker.PickWavFile()
	audio.ReadWavFile(audioFile)
	// visualizer.TermGraph()
	print(audioFile)
}
