package audio

import (
	"fmt"
	"github.com/go-audio/wav"
	"os"
)

func ReadWavFile(audioFile string) {
	// Open the .wav file
	file, err := os.Open(audioFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new decoder
	decoder := wav.NewDecoder(file)

	// Check if it's a valid wav file
	if !decoder.IsValidFile() {
		fmt.Println("Invalid WAV file.")
		return
	}

	// Decode the WAV file and retrieve PCM audio data
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		fmt.Println("Error decoding WAV:", err)
		return
	}

	// Access the PCM audio data
	fmt.Println("Number of samples:", buf.NumFrames())
	fmt.Println("Sample rate:", buf.Format.SampleRate)
	fmt.Println("Channels:", buf.Format.NumChannels)

	// For example, you can access individual samples like this:
	for i, sample := range buf.Data {
		fmt.Printf("Sample %d: %d\n", i, sample)
		if i > 10 { // Limiting the output for readability
			break
		}
	}
}
