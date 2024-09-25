package visualizer

import (
	"fmt"
	"vadimasKo/wav_visualizer/audio"
	// "github.com/NimbleMarkets/ntcharts/canvas"
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
	wlc        wavelinechart.Model
	zM         *zone.Manager
	audioProps *audio.AudioFileProperties
}

func (m Model) Init() tea.Cmd {
	m.wlc.Draw()
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	forwardMsg := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			// wavelinechart Clear() resets the canvas
			// and ClearData() resets internal data storage
			m.wlc.ClearAllData()
			m.wlc.Clear()
			m.wlc.DrawXYAxisAndLabel()
			m.wlc.Draw()

			return m, nil
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "down", "left", "right", "pgup", "pgdown":
			forwardMsg = true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			if m.zM.Get(m.wlc.ZoneID()).InBounds(msg) { // switch to canvas 1 if clicked on it
				m.wlc.Focus()
			} else {
				m.wlc.Blur()
			}
		}
		forwardMsg = true
	}

	// wavelinechart handles mouse events
	if forwardMsg {
		if m.wlc.Focused() {
			m.wlc, _ = m.wlc.Update(msg)
			m.wlc.DrawAll()
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

	s := lipgloss.JoinHorizontal(lipgloss.Top, defaultStyle.Render(top+m.wlc.View())+"\n")

	return m.zM.Scan(s) // Call zone Manager.Scan() at root modeli
}

func WavelineModel(properties *audio.AudioFileProperties) Model {
	width := 64
	height := 11
	xStep := 1
	yStep := 1
	minXValue := 0.0
	maxXValue := properties.Duration.Seconds()
	minYValue := -2.0
	maxYValue := 2.0

	// create new bubblezone Manager to enable mouse support to zoom in and out of chart
	zoneManager := zone.New()

	// wavelinechart 1 created with New() and SetStyle()
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
	return m
}
