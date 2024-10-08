package chart

import (
	"fmt"
	"wav_visualizer/pcm_energy_visualizer/audio"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func AudioLineChart(title string, channeledPoints [][]audio.AudioInputPoint) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width: "1920px",
		}),
	)

	for channelIndex, points := range channeledPoints {
		xPoints := make([]float64, len(points))
		convertedPoints := make([]opts.LineData, len(points))

		for i, point := range points {
			xPoints[i] = point.X
			convertedPoints[i] = opts.LineData{Value: point.Y}
		}

		channelName := fmt.Sprintf("channel%d", channelIndex)

		line.SetXAxis(xPoints).AddSeries(channelName, convertedPoints).
			SetSeriesOptions(charts.WithLineChartOpts(
				opts.LineChart{
					// Smooth: opts.Bool(true),
				}),
			)

	}

	return line
}
