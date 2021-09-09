package termui

import (
	"fmt"

	"github.com/robopuff/goradio/pkg/drivers"
	"github.com/robopuff/goradio/pkg/gui"
	"github.com/robopuff/goradio/pkg/stations"

	termui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	colorGray termui.Color = 8
)

type ui struct {
	debug             bool
	w, h              int
	fullLineFormatter string

	stationsList *stations.List

	uiPlayingParagraph *widgets.Paragraph
	uiFooterParagraph  *widgets.Paragraph
	uiStationsList     *widgets.List
	uiLoggerList       *widgets.List
	uiVolumeGauge      *volume

	drawables []termui.Drawable
}

func NewTermUI(stationsList *stations.List, debug bool) *ui {
	return &ui{
		debug:        debug,
		stationsList: stationsList,
	}
}

// Init initialize UI
func (u *ui) Init() error {
	if err := termui.Init(); nil != err {
		return err
	}

	u.w, u.h = termui.TerminalDimensions()
	u.fullLineFormatter = fmt.Sprintf("%%-%ds", u.w)

	u.uiPlayingParagraph = widgets.NewParagraph()
	u.uiPlayingParagraph.Border = false
	u.uiPlayingParagraph.TextStyle.Fg = termui.ColorRed

	u.uiFooterParagraph = widgets.NewParagraph()
	u.uiFooterParagraph.Text = fmt.Sprintf(u.fullLineFormatter, gui.HelpFooter)
	u.uiFooterParagraph.WrapText = false
	u.uiFooterParagraph.PaddingLeft = -1
	u.uiFooterParagraph.PaddingRight = -1
	u.uiFooterParagraph.Border = false
	u.uiFooterParagraph.TextStyle.Fg = termui.ColorBlack
	u.uiFooterParagraph.TextStyle.Bg = colorGray

	u.uiLoggerList = widgets.NewList()
	u.uiLoggerList.Title = "[ log ]"
	u.uiLoggerList.TextStyle.Fg = termui.ColorBlue
	u.uiLoggerList.BorderStyle.Fg = colorGray
	u.uiLoggerList.SelectedRowStyle.Fg = u.uiLoggerList.TextStyle.Fg
	u.uiLoggerList.Rows = []string{""}

	u.uiStationsList = widgets.NewList()
	u.uiStationsList.Title = "[ stations ]"
	u.uiStationsList.TextStyle.Fg = 8
	u.uiStationsList.TextStyle.Modifier = termui.ModifierBold
	u.uiStationsList.SelectedRowStyle.Modifier = termui.ModifierBold
	u.uiStationsList.SelectedRowStyle.Fg = termui.ColorWhite
	u.uiStationsList.SelectedRowStyle.Bg = colorGray
	u.uiStationsList.BorderStyle.Fg = colorGray
	u.uiStationsList.WrapText = false

	u.uiVolumeGauge = NewVolume()
	u.uiVolumeGauge.Border = false
	u.uiVolumeGauge.Percent = 0

	u.windowResize(termui.Event{
		Payload: termui.Resize{
			Width:  u.w,
			Height: u.h,
		},
	})
	u.drawables = []termui.Drawable{
		u.uiPlayingParagraph,
		u.uiVolumeGauge,
		u.uiFooterParagraph,
		u.uiStationsList,
	}

	if u.debug {
		u.drawables = append(u.drawables, u.uiLoggerList)
	}

	return nil
}

// Run run UI and events
func (u *ui) Run(d drivers.Driver) {
	u.render()

	events := newEventsManager(u, d)
	go events.manageDriverLogs()
	events.run()
}

func (u *ui) Close() {
	termui.Close()
	if r := recover(); r != nil {
		fmt.Printf("panic: %v\n", r)
	}
}

func (u *ui) windowResize(e termui.Event) {
	payload := e.Payload.(termui.Resize)
	termui.Clear()

	u.w = payload.Width
	u.h = payload.Height
	u.fullLineFormatter = fmt.Sprintf("%%-%ds", u.w)

	u.uiPlayingParagraph.SetRect(0, -1, u.w, 3)
	u.uiFooterParagraph.SetRect(0, u.h-3, u.w, u.h)
	u.uiLoggerList.SetRect(u.w/2, 1, u.w-1, u.h-2)
	u.uiVolumeGauge.SetRect(u.w-21, -1, u.w-1, 2)

	if u.debug {
		u.uiStationsList.SetRect(0, 1, (u.w/2)-1, u.h-2)
	} else {
		u.uiStationsList.SetRect(0, 1, u.w, u.h-2)
	}

	u.uiStationsList.Rows = u.stationsList.GetRows(u.uiStationsList.Size().X)
}

func (u *ui) render() {
	termui.Render(u.drawables...)
}

func (u *ui) log(m string) {
	u.uiLoggerList.Rows = append(u.uiLoggerList.Rows, fmt.Sprintf("%s", m))
	if len(u.uiLoggerList.Rows) > 1000 {
		u.uiLoggerList.Rows = u.uiLoggerList.Rows[500:]
	}

	u.uiLoggerList.ScrollBottom()
	if u.debug {
		termui.Render(u.uiLoggerList)
	}
}
