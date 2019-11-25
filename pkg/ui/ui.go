package ui

import (
	"fmt"
	"github.com/robopuff/goradio/pkg/driver"
	"github.com/robopuff/goradio/pkg/stations"

	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	currentlyPlayingFormat = "Currently playing: %s"
	helpFooter             = "k/↑ : Up | j/↓: Down | Enter: Select | p: Pause | m: Mute | s: Stop | +: Louder | -: Quieter | R: Refresh | q: Quit"
)

const (
	colorGray tui.Color = 8
)

var (
	current            = -1
	debug              bool
	w, h               int
	fullLineFormatter  string
	stationsList       *stations.List
	uiPlayingParagraph *widgets.Paragraph
	uiFooterParagraph  *widgets.Paragraph
	uiStationsList     *widgets.List
	uiLoggerList       *widgets.List
	volumeGauge        *widgets.Gauge
	drawables          []tui.Drawable
)

// Init initialize UI
func Init(csvStationsList *stations.List, debugFlag bool) error {
	debug = debugFlag
	if err := tui.Init(); nil != err {
		return err
	}

	w, h = tui.TerminalDimensions()

	stationsList = csvStationsList

	fullLineFormatter = fmt.Sprintf("%%-%ds", w)

	uiPlayingParagraph = widgets.NewParagraph()
	uiPlayingParagraph.SetRect(0, -1, w, 3)
	uiPlayingParagraph.Border = false
	uiPlayingParagraph.TextStyle.Fg = tui.ColorRed
	setCurrentlyPlaying("")

	uiFooterParagraph = widgets.NewParagraph()
	uiFooterParagraph.Text = fmt.Sprintf(fullLineFormatter, helpFooter)
	uiFooterParagraph.WrapText = false
	uiFooterParagraph.PaddingLeft = -1
	uiFooterParagraph.PaddingRight = -1
	uiFooterParagraph.SetRect(0, h-3, w, h)
	uiFooterParagraph.Border = false
	uiFooterParagraph.TextStyle.Fg = tui.ColorBlack
	uiFooterParagraph.TextStyle.Bg = colorGray

	uiLoggerList = widgets.NewList()
	uiLoggerList.Title = "[ sendToLog ]"
	uiLoggerList.SetRect(w/2, 1, w-1, h-2)
	uiLoggerList.TextStyle.Fg = tui.ColorBlue
	uiLoggerList.BorderStyle.Fg = colorGray
	uiLoggerList.SelectedRowStyle.Fg = uiLoggerList.TextStyle.Fg
	uiLoggerList.Rows = []string{""}

	uiStationsList = widgets.NewList()
	uiStationsList.Title = "[ stations ]"
	uiStationsList.TextStyle.Fg = 8
	uiStationsList.TextStyle.Modifier = tui.ModifierBold
	uiStationsList.SelectedRowStyle.Modifier = tui.ModifierBold
	uiStationsList.SelectedRowStyle.Fg = tui.ColorWhite
	uiStationsList.SelectedRowStyle.Bg = colorGray
	uiStationsList.BorderStyle.Fg = colorGray
	uiStationsList.WrapText = false

	volumeGauge = widgets.NewGauge()
	volumeGauge.Border = false
	volumeGauge.SetRect(w-21, -1, w-1, 2)
	volumeGauge.Percent = 25
	volumeGauge.Label = "25"

	if debug {
		uiStationsList.SetRect(0, 1, (w/2)-1, h-2)
	} else {
		uiStationsList.SetRect(0, 1, w, h-2)
	}

	uiStationsList.Rows = csvStationsList.GetRows(uiStationsList.Size().X)

	drawables = []tui.Drawable{
		uiPlayingParagraph,
		volumeGauge,
		uiFooterParagraph,
		uiStationsList,
	}

	if debug {
		drawables = append(drawables, uiLoggerList)
	}

	return nil
}

// Run run UI and events
func Run(d driver.Driver) {
	render()

	uiEvents := tui.PollEvents()

	defer (func() {
		d.Close()
		tui.Close()
	})()

	go manageDriverLogs(d)

	for {
		select {
		case e := <-uiEvents:
			if r := manageKeyboardEvent(e, d); r != 0 {
				return
			}
			render()
		}
	}
}

func render() {
	tui.Render(drawables...)
}

func sendToLog(m string) {
	if !debug {
		return
	}

	uiLoggerList.Rows = append(uiLoggerList.Rows, fmt.Sprintf("%s", m))
	uiLoggerList.ScrollBottom()
	tui.Render(uiLoggerList)
}
