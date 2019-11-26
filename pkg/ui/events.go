package ui

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/robopuff/goradio/pkg/driver"

	tui "github.com/gizak/termui/v3"
)

const (
	regexTitle  = `(?m)^ICY Info: StreamTitle='(.*?)';`
	regexVolume = `(?m)Volume: (\d+)`
)

var colorSelected = "(fg:black,bg:white)"

func manageKeyboardEvent(e tui.Event, d driver.Driver) int {
	if e.Type == tui.ResizeEvent {
		windowResize(e)
	}

	switch e.ID {
	case "<Enter>", "<MouseLeft>":
		selected := uiStationsList.SelectedRow
		selectedStation := stationsList.GetSelected(selected)
		currentStation := stationsList.GetSelected(current)

		if selected == current {
			d.Pause()
			return 0
		}

		if selectedStation == nil {
			return 0
		}

		if currentStation != nil {
			removeStationSelection()
		}

		d.Play(selectedStation.URL)
		current = selected
		setVolumeGauge("25")
		addStationSelection()
	case "s":
		if current < 0 {
			return 0
		}

		d.Close()
		removeStationSelection()
		setVolumeGauge("")
		current = -1
	case "R":
		d.Close()
		current = -1
		stationsList.Reload()
		uiStationsList.Rows = stationsList.GetRows(uiStationsList.Size().X)
	case "m":
		d.Mute()
	case "p":
		d.Pause()
	case "k", "<Up>", "<MouseWheelUp>":
		uiStationsList.ScrollUp()
	case "j", "<Down>", "<MouseWheelDown>":
		uiStationsList.ScrollDown()
	case "K", "<PageUp>":
		uiStationsList.ScrollPageUp()
	case "J", "<PageDown>":
		uiStationsList.ScrollPageDown()
	case "h", "<Left>":
		uiLoggerList.ScrollPageUp()
	case "l", "<Right>":
		uiLoggerList.ScrollPageDown()
	case "+", "=":
		d.IncVolume()
	case "-":
		d.DecVolume()
	case "q", "<C-c>", "<Esc>":
		d.Close()
		return 1
	}

	return 0
}

func manageDriverLogs(d driver.Driver) {
	titleRegex := regexp.MustCompile(regexTitle)
	volumeRegex := regexp.MustCompile(regexVolume)

	for {
		select {
		case outPipe := <-d.PipeChan():
			reader := bufio.NewReader(outPipe)
			for {
				data, err := reader.ReadString('\n')
				if err != nil {
					setCurrentlyPlaying("")
					sendToLog(fmt.Sprintf("Pipe closed: %v", err.Error()))
					break
				}

				match := titleRegex.FindStringSubmatch(data)
				if len(match) > 0 {
					setCurrentlyPlaying(match[1])
				}

				match = volumeRegex.FindStringSubmatch(data)
				if len(match) > 0 {
					setVolumeGauge(match[1])
					continue
				}

				if data != "" {
					sendToLog(strings.Trim(data, "\n"))
				}
			}
		}
	}
}

func setVolumeGauge(value string) {
	volume, _ := strconv.Atoi(value)
	volumeGauge.Percent = volume
	volumeGauge.Visible = true

	if volume == 0 {
		volumeGauge.Visible = false
	}

	render()
}

func setCurrentlyPlaying(currently string) {
	format := fmt.Sprintf(fullLineFormatter, currentlyPlayingFormat)
	if currently == "" {
		format = "%s"
		currently = ""
	}
	uiPlayingParagraph.Text = fmt.Sprintf(format, currently)
	render()
}

func addStationSelection() {
	uiStationsList.Rows[current] = fmt.Sprintf("[%v]%s", uiStationsList.Rows[current], colorSelected)
}

func removeStationSelection() {
	rowString := string(uiStationsList.Rows[current])
	uiStationsList.Rows[current] = rowString[1 : len(rowString)-1-len(colorSelected)]
}
