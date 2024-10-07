package visualizer

import (
	"fmt"
	"strconv"
	"wav_visualizer/pcm_energy_visualizer/audio"

	"github.com/NimbleMarkets/ntcharts/canvas"
	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/linechart/wavelinechart"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var defaultStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("5"))

var graphLineStyle1 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("5"))

var axisStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fff"))

var labelStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#fff"))

var propertyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("5")).
	Padding(0, 1)

type Model struct {
	Wlc        wavelinechart.Model
	zM         *zone.Manager
	audioProps *audio.AudioFileProperties
}

func (m Model) Init() tea.Cmd {
	m.Wlc.Draw()
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	forwardMsg := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "down", "left", "right", "pgup", "pgdown":
			forwardMsg = true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			if m.zM.Get(m.Wlc.ZoneID()).InBounds(msg) { // switch to canvas 1 if clicked on it
				m.Wlc.Focus()
			} else {
				m.Wlc.Blur()
			}
		}
		forwardMsg = true
	}

	// wavelinechart handles mouse events
	if forwardMsg {
		if m.Wlc.Focused() {
			m.Wlc, _ = m.Wlc.Update(msg)
			m.Wlc.DrawAll()
		}
	}

	return m, nil
}

func (m Model) View() string {
	// Format each property manually
	fileName := propertyStyle.Render(fmt.Sprintf("FileName: %s", m.audioProps.FileName))
	channelCount := propertyStyle.Render(fmt.Sprintf("ChannelCount: %d", m.audioProps.ChannelCount))
	depth := propertyStyle.Render(fmt.Sprintf("Depth: %dbit", m.audioProps.Depth))
	sampleRate := propertyStyle.Render(fmt.Sprintf("SampleRate: %dhz", m.audioProps.SampleRate))

	// Combine properties into a single formatted string
	top := fmt.Sprintf(
		"%s\n%s | %s | %s\n",
		fileName,
		channelCount,
		depth,
		sampleRate,
	)

	s := lipgloss.JoinHorizontal(lipgloss.Top, defaultStyle.Render(top+m.Wlc.View())+"\n")
	return m.zM.Scan(s) // Call zone Manager.Scan() at root modeli
}

func PlotMultiChannelData(m *Model, points [][]canvas.Float64Point) {
	colors := []lipgloss.Color{
		lipgloss.Color("22"),
		lipgloss.Color("4"),
		lipgloss.Color("7"),
		lipgloss.Color("2"),
	}

	for i, channelPoints := range points {
		dataSetId := strconv.Itoa(i)

		fmt.Printf("Starting plotting for channel %s, total points: %d\n", dataSetId, len(channelPoints))

		m.Wlc.SetDataSetStyles(
			dataSetId,
			runes.ThinLineStyle,
			lipgloss.NewStyle().Foreground(colors[i%len(colors)]),
		)

		for pointIndex, point := range channelPoints {
			fmt.Printf("Adding point #%d to channel %s: (X: %f, Y: %f)\n", pointIndex+1, dataSetId, point.X, point.Y)

			// m.Wlc.Plot(point)
			m.Wlc.PlotDataSet(dataSetId, point)
		}

		fmt.Printf("Finished plotting channel %s, plotted %d points.\n", dataSetId, len(channelPoints))
	}

	fmt.Println("Finalizing drawing of all datasets.")
	m.Wlc.DrawAll()
}

func WavelineModel(properties *audio.AudioFileProperties) *Model {
	width := 64
	height := 11
	xStep := 1
	yStep := 1
	minXValue := 0.0
	maxXValue := properties.Duration.Seconds()
	minYValue := -1.0
	maxYValue := 1.0

	// create new bubblezone Manager to enable mouse suport to zoom in and out of chart
	zoneManager := zone.New()

	wlc := wavelinechart.New(width, height)
	wlc.AxisStyle = axisStyle
	wlc.LabelStyle = labelStyle
	wlc.SetXStep(xStep)
	wlc.SetYStep(yStep)
	wlc.SetXYRange(minXValue, maxXValue, minYValue, maxYValue)
	wlc.SetViewXYRange(0, 60, -60, 60)
	wlc.SetStyles(runes.ThinLineStyle, graphLineStyle1)
	wlc.SetZoneManager(zoneManager)
	wlc.Focus()
	m := Model{wlc, zoneManager, properties}
	return &m
}
