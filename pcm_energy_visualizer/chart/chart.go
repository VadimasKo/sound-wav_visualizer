package chart

import (
	"fmt"
	"wav_visualizer/pcm_energy_visualizer/audio"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type ChartOptions struct {
	Title           string
	ChanneledPoints [][]audio.AudioInputPoint
	Segments        [][]audio.AudioSegment
}

func (options *ChartOptions) CreateAudioLineChart() *charts.Line {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: options.Title}),
		charts.WithInitializationOpts(opts.Initialization{Width: "1920px"}),
	)

	for channelIndex, points := range options.ChanneledPoints {
		xPoints := make([]float64, len(points))
		convertedPoints := make([]opts.LineData, len(points))
		channelName := fmt.Sprintf("channel %d", channelIndex)

		for i, point := range points {
			xPoints[i] = point.X
			convertedPoints[i] = opts.LineData{Value: point.Y}
		}

		line.SetXAxis(xPoints)
		line.AddSeries(channelName, convertedPoints)
		line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

		if len(options.Segments) == 0 {
			continue
		}

		channelSegments := options.Segments[channelIndex]

		for segmentIndex, segment := range channelSegments {
			segmentName := fmt.Sprintf("segment %d", segmentIndex)

			// fmt.Printf("segment %d %f %f \n", segmentIndex, segment.Start, segment.End)

			markLineStart := opts.MarkLineNameXAxisItem{
				Name:  segmentName + " (start)",
				XAxis: segment.StartSampleIndex,
			}

			markLineEnd := opts.MarkLineNameXAxisItem{
				Name:  segmentName + " (end)",
				XAxis: segment.EndSampleIndex,
			}

			line.SetSeriesOptions(
				charts.WithMarkLineNameXAxisItemOpts(markLineStart),
				charts.WithMarkLineNameXAxisItemOpts(markLineEnd),
			)
		}

	}

	return line
}
