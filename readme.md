## About 🔍
This repository contains projects centered around manipulation of sound signals.

## wav_visualizer
This is tui application, which reads .wav file and displays signal in time graph.

1. User selects `.wav` file using `bubbletea file-picker`
2. Application, decodes main information about file (bitDepth, sampleRate...)
3. Creates `waveline` chart model which will handle point rendering
4. Starts reading `.wav` file in batches and processing it using go routines
5. Go routine sample processing:
   - Normalizes signal value [-1, 1] based on maximum value of bitDepth
   - Calculates signal time stamp, based on sampleIndex and channelCount
   - Aggregates parsed results inside 2d points channel
6. After agregation points are added to `waveline` model
7. Result is displayed utilising `bubbletea` framework.
<img width="733" alt="Screenshot 2024-09-28 at 16 30 40" src="https://github.com/user-attachments/assets/cdaad0da-4791-4728-9fc8-3f61eef3c207">

### ...
