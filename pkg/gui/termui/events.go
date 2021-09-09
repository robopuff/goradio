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
	colorSelected = `(fg:black,bg:white)`
)

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

func (em *eventsManager) run() {
	uiEvents := tui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if r := em.manageKeyboardEvent(e); r != 0 {
				return
			}
			em.ui.render()
		}
	}
}

func (em *eventsManager) manageKeyboardEvent(e tui.Event) int {
	if e.Type == tui.ResizeEvent {
		em.ui.windowResize(e)
	}

	switch e.ID {
	case "<Enter>":
		selected := em.ui.uiStationsList.SelectedRow
		selectedStation := em.ui.stationsList.GetSelected(selected)
		currentStation := em.ui.stationsList.GetSelected(em.current)

		if selected == em.current {
			em.driver.Pause()
			return 0
		}

		if selectedStation == nil {
			return 0
		}

		if currentStation != nil {
			em.removeStationSelection()
		}

		em.driver.Play(selectedStation.URL)
		em.current = selected
		em.setVolumeGauge("25")
		em.addStationSelection()
	case "s":
		if em.current < 0 {
			return 0
		}

		em.driver.Close()
		em.removeStationSelection()
		em.setVolumeGauge("")
		em.current = -1
	case "R":
		em.driver.Close()
		em.current = -1
		em.ui.stationsList.Reload()
		em.ui.uiStationsList.Rows = em.ui.stationsList.GetRows(em.ui.uiStationsList.Size().X)
	case "m":
		em.driver.Mute()
	case "p":
		em.driver.Pause()
	case "k", "<Up>":
		em.ui.uiStationsList.ScrollUp()
	case "j", "<Down>":
		em.ui.uiStationsList.ScrollDown()
	case "K", "<PageUp>":
		em.ui.uiStationsList.ScrollPageUp()
	case "J", "<PageDown>":
		em.ui.uiStationsList.ScrollPageDown()
	case "h", "<Left>", "<MouseWheelUp>":
		em.ui.uiLoggerList.ScrollPageUp()
	case "l", "<Right>", "<MouseWheelDown>":
		em.ui.uiLoggerList.ScrollPageDown()
	case "+", "=":
		em.driver.IncVolume()
	case "-":
		em.driver.DecVolume()
	case "q", "<C-c>", "<Esc>":
		em.driver.Close()
		return 1
	}

	return 0
}

func (em *eventsManager) manageDriverLogs() {
	titleRegex := regexp.MustCompile(regexTitle)
	volumeRegex := regexp.MustCompile(regexVolume)

	for {
		select {
		case outPipe := <-em.driver.PipeChan():
			reader := bufio.NewReader(outPipe)
			for {
				data, err := reader.ReadString('\n')
				if err != nil {
					em.setCurrentlyPlaying("")
					em.ui.log(fmt.Sprintf("Pipe closed: %v", err.Error()))
					break
				}

				match := titleRegex.FindStringSubmatch(data)
				if len(match) > 0 {
					em.setCurrentlyPlaying(match[1])
				}

				match = volumeRegex.FindStringSubmatch(data)
				if len(match) > 0 {
					em.setVolumeGauge(match[1])
					continue
				}

				if data != "" {
					em.ui.log(strings.Trim(data, "\n"))
				}
			}
		}
	}
}

func (em *eventsManager) setVolumeGauge(value string) {
	volume, _ := strconv.Atoi(value)
	em.ui.uiVolumeGauge.Percent = volume
	em.ui.uiVolumeGauge.Visible = true

	if volume == 0 {
		em.ui.uiVolumeGauge.Visible = false
	}

	em.ui.render()
}

func (em *eventsManager) setCurrentlyPlaying(currently string) {
	format := fmt.Sprintf(em.ui.fullLineFormatter, gui.CurrentlyPlayingFormat)
	if currently == "" {
		format = "%s"
		currently = ""
	}
	em.ui.uiPlayingParagraph.Text = fmt.Sprintf(format, currently)
	em.ui.render()
}

func (em *eventsManager) addStationSelection() {
	em.ui.uiStationsList.Rows[em.current] = fmt.Sprintf("[%v]%s", em.ui.uiStationsList.Rows[em.current], colorSelected)
}

func (em *eventsManager) removeStationSelection() {
	rowString := string(em.ui.uiStationsList.Rows[em.current])
	em.ui.uiStationsList.Rows[em.current] = rowString[1 : len(rowString)-1-len(colorSelected)]
}
