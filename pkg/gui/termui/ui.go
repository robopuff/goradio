package termui

import (
	"fmt"

	"github.com/robopuff/goradio/pkg/drivers"
	"github.com/robopuff/goradio/pkg/gui"
	"github.com/robopuff/goradio/pkg/stations"

	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	colorGray tui.Color = 8
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

	drawables []tui.Drawable
}

func NewTermUI(stationsList *stations.List, debug bool) *ui {
	return &ui{
		debug:        debug,
		stationsList: stationsList,
	}
}

// Init initialize UI
func (self *ui) Init() error {
	if err := tui.Init(); nil != err {
		return err
	}

	self.w, self.h = tui.TerminalDimensions()
	self.fullLineFormatter = fmt.Sprintf("%%-%ds", self.w)

	self.uiPlayingParagraph = widgets.NewParagraph()
	self.uiPlayingParagraph.Border = false
	self.uiPlayingParagraph.TextStyle.Fg = tui.ColorRed

	self.uiFooterParagraph = widgets.NewParagraph()
	self.uiFooterParagraph.Text = fmt.Sprintf(self.fullLineFormatter, gui.HelpFooter)
	self.uiFooterParagraph.WrapText = false
	self.uiFooterParagraph.PaddingLeft = -1
	self.uiFooterParagraph.PaddingRight = -1
	self.uiFooterParagraph.Border = false
	self.uiFooterParagraph.TextStyle.Fg = tui.ColorBlack
	self.uiFooterParagraph.TextStyle.Bg = colorGray

	self.uiLoggerList = widgets.NewList()
	self.uiLoggerList.Title = "[ log ]"
	self.uiLoggerList.TextStyle.Fg = tui.ColorBlue
	self.uiLoggerList.BorderStyle.Fg = colorGray
	self.uiLoggerList.SelectedRowStyle.Fg = self.uiLoggerList.TextStyle.Fg
	self.uiLoggerList.Rows = []string{""}

	self.uiStationsList = widgets.NewList()
	self.uiStationsList.Title = "[ stations ]"
	self.uiStationsList.TextStyle.Fg = 8
	self.uiStationsList.TextStyle.Modifier = tui.ModifierBold
	self.uiStationsList.SelectedRowStyle.Modifier = tui.ModifierBold
	self.uiStationsList.SelectedRowStyle.Fg = tui.ColorWhite
	self.uiStationsList.SelectedRowStyle.Bg = colorGray
	self.uiStationsList.BorderStyle.Fg = colorGray
	self.uiStationsList.WrapText = false

	self.uiVolumeGauge = NewVolume()
	self.uiVolumeGauge.Border = false
	self.uiVolumeGauge.Percent = 0

	self.windowResize(tui.Event{
		Payload: tui.Resize{
			Width:  self.w,
			Height: self.h,
		},
	})
	self.drawables = []tui.Drawable{
		self.uiPlayingParagraph,
		self.uiVolumeGauge,
		self.uiFooterParagraph,
		self.uiStationsList,
	}

	if self.debug {
		self.drawables = append(self.drawables, self.uiLoggerList)
	}

	return nil
}

// Run run UI and events
func (self *ui) Run(d drivers.Driver) {
	self.render()

	events := newEventsManager(self, d)
	go events.manageDriverLogs()
	events.run()
}

func (self *ui) Close() {
	tui.Close()
	if r := recover(); r != nil {
		fmt.Printf("panic: %v\n", r)
	}
}

func (self *ui) windowResize(e tui.Event) {
	payload := e.Payload.(tui.Resize)
	tui.Clear()

	self.w = payload.Width
	self.h = payload.Height
	self.fullLineFormatter = fmt.Sprintf("%%-%ds", self.w)

	self.uiPlayingParagraph.SetRect(0, -1, self.w, 3)
	self.uiFooterParagraph.SetRect(0, self.h-3, self.w, self.h)
	self.uiLoggerList.SetRect(self.w/2, 1, self.w-1, self.h-2)
	self.uiVolumeGauge.SetRect(self.w-21, -1, self.w-1, 2)

	if self.debug {
		self.uiStationsList.SetRect(0, 1, (self.w/2)-1, self.h-2)
	} else {
		self.uiStationsList.SetRect(0, 1, self.w, self.h-2)
	}

	self.uiStationsList.Rows = self.stationsList.GetRows(self.uiStationsList.Size().X)
}

func (self *ui) render() {
	tui.Render(self.drawables...)
}

func (self *ui) log(m string) {
	self.uiLoggerList.Rows = append(self.uiLoggerList.Rows, fmt.Sprintf("%s", m))
	if len(self.uiLoggerList.Rows) > 1000 {
		self.uiLoggerList.Rows = self.uiLoggerList.Rows[500:]
	}

	self.uiLoggerList.ScrollBottom()
	if self.debug {
		tui.Render(self.uiLoggerList)
	}
}
