package termui

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/robopuff/goradio/pkg/drivers"
	"github.com/robopuff/goradio/pkg/gui"

	tui "github.com/gizak/termui/v3"
)

const (
	regexTitle  = `(?m)^ICY Info: StreamTitle='(.*?)';`
	regexVolume = `(?m)Volume: (\d+)`
)

var colorSelected = "(fg:black,bg:white)"

type eventsManager struct {
	ui      *ui
	driver  drivers.Driver
	current int
}

func newEventsManager(ui *ui, driver drivers.Driver) *eventsManager {
	return &eventsManager{
		ui:      ui,
		current: -1,
		driver:  driver,
	}
}

func (self *eventsManager) run() {
	uiEvents := tui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if r := self.manageKeyboardEvent(e); r != 0 {
				return
			}
			self.ui.render()
		}
	}
}

func (self *eventsManager) manageKeyboardEvent(e tui.Event) int {
	if e.Type == tui.ResizeEvent {
		self.ui.windowResize(e)
	}

	switch e.ID {
	case "<Enter>":
		selected := self.ui.uiStationsList.SelectedRow
		selectedStation := self.ui.stationsList.GetSelected(selected)
		currentStation := self.ui.stationsList.GetSelected(self.current)

		if selected == self.current {
			self.driver.Pause()
			return 0
		}

		if selectedStation == nil {
			return 0
		}

		if currentStation != nil {
			self.removeStationSelection()
		}

		self.driver.Play(selectedStation.URL)
		self.current = selected
		self.setVolumeGauge("25")
		self.addStationSelection()
	case "s":
		if self.current < 0 {
			return 0
		}

		self.driver.Close()
		self.removeStationSelection()
		self.setVolumeGauge("")
		self.current = -1
	case "R":
		self.driver.Close()
		self.current = -1
		self.ui.stationsList.Reload()
		self.ui.uiStationsList.Rows = self.ui.stationsList.GetRows(self.ui.uiStationsList.Size().X)
	case "m":
		self.driver.Mute()
	case "p":
		self.driver.Pause()
	case "k", "<Up>":
		self.ui.uiStationsList.ScrollUp()
	case "j", "<Down>":
		self.ui.uiStationsList.ScrollDown()
	case "K", "<PageUp>":
		self.ui.uiStationsList.ScrollPageUp()
	case "J", "<PageDown>":
		self.ui.uiStationsList.ScrollPageDown()
	case "h", "<Left>", "<MouseWheelUp>":
		self.ui.uiLoggerList.ScrollPageUp()
	case "l", "<Right>", "<MouseWheelDown>":
		self.ui.uiLoggerList.ScrollPageDown()
	case "+", "=":
		self.driver.IncVolume()
	case "-":
		self.driver.DecVolume()
	case "q", "<C-c>", "<Esc>":
		self.driver.Close()
		return 1
	}

	return 0
}

func (self *eventsManager) manageDriverLogs() {
	titleRegex := regexp.MustCompile(regexTitle)
	volumeRegex := regexp.MustCompile(regexVolume)

	for {
		select {
		case outPipe := <-self.driver.PipeChan():
			reader := bufio.NewReader(outPipe)
			for {
				data, err := reader.ReadString('\n')
				if err != nil {
					self.setCurrentlyPlaying("")
					self.ui.log(fmt.Sprintf("Pipe closed: %v", err.Error()))
					break
				}

				match := titleRegex.FindStringSubmatch(data)
				if len(match) > 0 {
					self.setCurrentlyPlaying(match[1])
				}

				match = volumeRegex.FindStringSubmatch(data)
				if len(match) > 0 {
					self.setVolumeGauge(match[1])
					continue
				}

				if data != "" {
					self.ui.log(strings.Trim(data, "\n"))
				}
			}
		}
	}
}

func (self *eventsManager) setVolumeGauge(value string) {
	volume, _ := strconv.Atoi(value)
	self.ui.volumeGauge.Percent = volume
	self.ui.volumeGauge.Visible = true

	if volume == 0 {
		self.ui.volumeGauge.Visible = false
	}

	self.ui.render()
}

func (self *eventsManager) setCurrentlyPlaying(currently string) {
	format := fmt.Sprintf(self.ui.fullLineFormatter, gui.CurrentlyPlayingFormat)
	if currently == "" {
		format = "%s"
		currently = ""
	}
	self.ui.uiPlayingParagraph.Text = fmt.Sprintf(format, currently)
	self.ui.render()
}

func (self *eventsManager) addStationSelection() {
	self.ui.uiStationsList.Rows[self.current] = fmt.Sprintf("[%v]%s", self.ui.uiStationsList.Rows[self.current], colorSelected)
}

func (self *eventsManager) removeStationSelection() {
	rowString := string(self.ui.uiStationsList.Rows[self.current])
	self.ui.uiStationsList.Rows[self.current] = rowString[1 : len(rowString)-1-len(colorSelected)]
}
