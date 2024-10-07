package chart

import (
	"io"
	"math/rand"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	itemCntLine = 20000
	fruits      = []string{"Apple", "Banana", "Peach ", "Lemon", "Pear", "Cherry"}
)

func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < itemCntLine; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}

func generateLineData(data []float32) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(data); i++ {
		items = append(items, opts.LineData{Value: data[i]})
	}
	return items
}

func lineSmooth() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "smooth style",
		}),
	)

	line.SetXAxis(fruits).AddSeries("Category A", generateLineItems()).
		SetSeriesOptions(charts.WithLineChartOpts(
			opts.LineChart{
				Smooth: opts.Bool(true),
			}),
		)
	return line
}

type LineExamples struct{}

func (LineExamples) Examples() {
	page := components.NewPage()
	page.AddCharts(
		lineSmooth(),
	)

	dirPath := "examples/html"

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(dirPath + "/chart.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
