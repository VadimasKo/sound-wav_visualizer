package visualizer

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/NimbleMarkets/ntcharts/canvas"
	"github.com/NimbleMarkets/ntcharts/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/linechart/wavelinechart"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var randomFloat64Point1 canvas.Float64Point
var randomFloat64Point2 canvas.Float64Point

var defaultStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63")) // purple

var graphLineStyle1 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("4")) // blue

var graphLineStyle2 = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9")) // red

var axisStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3")) // yellow

var labelStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("6")) // cyan

type model struct {
	wlc wavelinechart.Model
	zM  *zone.Manager
}

func (m model) Init() tea.Cmd {
	m.wlc.Draw()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	addPoint := false
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
		default:
			addPoint = true
		}
	}
	if addPoint {
		// generate a random points within the given X,Y value ranges
		dx := m.wlc.MaxX() - m.wlc.MinX()
		dy := m.wlc.MaxY() - m.wlc.MinY()
		xRand1 := rand.Float64()*dx + m.wlc.MinX()
		yRand1 := rand.Float64()*dy + m.wlc.MinY()
		xRand2 := rand.Float64()*dx + m.wlc.MinX()
		yRand2 := rand.Float64()*dy + m.wlc.MinY()
		randomFloat64Point1 = canvas.Float64Point{X: xRand1, Y: yRand1}
		randomFloat64Point2 = canvas.Float64Point{X: xRand2, Y: yRand2}

		//  wavelinechart 1 plots random point 1 to default data set
		m.wlc.Plot(randomFloat64Point1)
		m.wlc.Draw()
	}

	// wavelinechart handles mouse events
	if forwardMsg {
		m.wlc, _ = m.wlc.Update(msg)
		m.wlc.DrawAll()
	}

	return m, nil
}

func (m model) View() string {
	top := fmt.Sprintf("file name:_ \n")
	top += fmt.Sprintf("channel count:_ sampling rate:_ quantization depth:_ \n")
	s := lipgloss.JoinHorizontal(lipgloss.Top, defaultStyle.Render(top+m.wlc.View())) + "\n"
	s += "`q/ctrl+c` to quit\n"
	// s += "pgup/pdown/mouse wheel scroll to zoom in and out\n"
	// s += "+arrow keys to move view while zoomed in\n"
	return m.zM.Scan(s) // call zone Manager.Scan() at root model
}

// The signal graph should include the name of the selected file
// and the main audio quality indicators: number of channels,
// the main parameters of the channels, such as sampling rate, quantization depth.

func TermGraph() {
	width := 92
	height := 16

	// generated based on audio file
	xStep := 1
	yStep := 1
	minXValue := 0.0
	maxXValue := 10.0
	minYValue := -5.0
	maxYValue := 5.0

	// create new bubblezone Manager to enable mouse support to zoom in and out of chart
	zoneManager := zone.New()

	// wavelinechart 1 created with New() and SetStyle()
	wlc := wavelinechart.New(width, height)
	wlc.AxisStyle = axisStyle
	wlc.LabelStyle = labelStyle
	wlc.SetXStep(xStep)
	wlc.SetYStep(yStep)
	wlc.SetXYRange(minXValue, maxXValue, minYValue, maxYValue)     // set expected ranges (can be less than or greater than displayed)
	wlc.SetViewXYRange(minXValue, maxXValue, minYValue, maxYValue) // setting displayed ranges fails unless setting expected values first
	wlc.SetStyles(runes.ThinLineStyle, graphLineStyle1)            // graphLineStyle1 replaces linechart rune style
	wlc.SetZoneManager(zoneManager)
	wlc.Focus()

	m := model{wlc, zoneManager}
	if _, err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
