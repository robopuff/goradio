package ui

import (
	"bufio"
	"fmt"

	tui "github.com/gizak/termui/v3"
	"github.com/robopuff/goradio/pkg/driver"
)

var colorSelected = "(fg:black,bg:white)"

func manageKeyboardEvent(e tui.Event, d driver.Driver) int {
	if e.Type == tui.ResizeEvent {
		tui.Clear()
		Init(stationsList, debug)
	}

	switch e.ID {
	case "<Enter>":
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
		addStationSelection()
	case "s":
		if current < 0 {
			return 0
		}

		d.Close()
		removeStationSelection()
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
	case "k", "<Up>":
		uiStationsList.ScrollUp()
	case "j", "<Down>":
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
	if !debug {
		return
	}

	for {
		select {
		case outPipe := <-d.PipeChan():
			reader := bufio.NewReader(outPipe)
			for {
				data, err := reader.ReadString('\n')
				if err != nil {
					log(fmt.Sprintf("Pipe closed: %v", err.Error()))
					break
				}
				log(data[:len(data)-1])
			}
		}
	}
}

func addStationSelection() {
	uiStationsList.Rows[current] = fmt.Sprintf("[%v]%s", uiStationsList.Rows[current], colorSelected)
}

func removeStationSelection() {
	rowString := string(uiStationsList.Rows[current])
	uiStationsList.Rows[current] = rowString[1 : len(rowString)-1-len(colorSelected)]
}
