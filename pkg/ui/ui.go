package ui

import (
	"fmt"
	"time"

	"github.com/robopuff/goradio/pkg/driver"
	"github.com/robopuff/goradio/pkg/stations"

	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	debug            bool
	w, h, current    int
	stationsConf     *stations.StationsList
	playingParagraph *widgets.Paragraph
	footerParagraph  *widgets.Paragraph
	stationsList     *widgets.List
	loggerList       *widgets.List
	drawables        []tui.Drawable
)

// Init initialize UI
func Init(stationsConfList *stations.StationsList, debugFlag bool) error {
	if err := tui.Init(); nil != err {
		return err
	}

	debug = debugFlag
	current = -1

	w, h = tui.TerminalDimensions()

	stationsConf = stationsConfList

	formatter := fmt.Sprintf("%%-%ds", w)

	playingParagraph = widgets.NewParagraph()
	playingParagraph.Text = fmt.Sprintf(formatter, "Currently playing: None")
	playingParagraph.SetRect(0, -1, w, 3)
	playingParagraph.Border = false
	playingParagraph.TextStyle.Fg = tui.ColorRed

	footerParagraph = widgets.NewParagraph()
	footerParagraph.Text = fmt.Sprintf(formatter, "k/↑ : Up | j/↓: Down | Enter: Select | p: Pause | m: Mute | s: Stop | +: Louder | -: Quieter | q: Quit")
	footerParagraph.WrapText = false
	footerParagraph.PaddingLeft = -1
	footerParagraph.PaddingRight = -1
	footerParagraph.SetRect(0, h-3, w, h)
	footerParagraph.Border = false
	footerParagraph.TextStyle.Fg = tui.ColorBlack
	footerParagraph.TextStyle.Bg = 8

	loggerList = widgets.NewList()
	loggerList.Title = "[ log ]"
	loggerList.SetRect(w/2, 1, w-1, h-2)
	loggerList.TextStyle.Fg = tui.ColorBlue
	loggerList.BorderStyle.Fg = 8
	loggerList.SelectedRowStyle.Fg = loggerList.TextStyle.Fg
	loggerList.Rows = []string{""}

	stationsList = widgets.NewList()
	stationsList.Title = "[ stations ]"
	stationsList.Rows = stationsConf.GetRows()
	stationsList.SelectedRowStyle.Fg = tui.ColorBlack
	stationsList.SelectedRowStyle.Bg = tui.ColorWhite
	stationsList.BorderStyle.Fg = 8

	if debug {
		stationsList.SetRect(0, 1, (w/2)-1, h-2)
	} else {
		stationsList.SetRect(0, 1, w, h-2)
	}

	drawables = []tui.Drawable{
		playingParagraph,
		footerParagraph,
		stationsList,
	}

	if debug {
		drawables = append(drawables, loggerList)
	}

	return nil
}

// Run run UI and events
func Run(d driver.Driver) {
	render()

	uiEvents := tui.PollEvents()
	ticker := time.NewTicker(time.Second).C

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
		case <-ticker:
			render()
		}
	}
}

func render() {
	tui.Render(drawables...)
}

func log(m string) {
	if !debug {
		return
	}

	loggerList.Rows = append(loggerList.Rows, fmt.Sprintf("%s", m))
	loggerList.ScrollBottom()
	tui.Render(loggerList)
}
