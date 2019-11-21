package ui

import (
	"bufio"
	"fmt"

	tui "github.com/gizak/termui/v3"
	"github.com/robopuff/goradio/pkg/driver"
)

func manageKeyboardEvent(e tui.Event, d driver.Driver) int {
	if e.Type == tui.ResizeEvent {
		tui.Clear()
		Init(stationsConf, debug)
	}

	switch e.ID {
	case "q", "<C-c>", "<Esc>":
		d.Close()
		return 1
	case "m":
		d.Mute()
	case "<Enter>":
		selected := stationsList.SelectedRow
		selectedStation := stationsConf.GetSelected(selected)

		currentStation := stationsConf.GetSelected(current)

		if selected == current {
			d.Pause()
			return 0
		}

		if currentStation != nil {
			stationsList.Rows[current] = currentStation.Name
		}

		if selectedStation == nil {
			log("Invalid station selected")
			return 0
		}

		stationsList.Rows[selected] = fmt.Sprintf("* %s", selectedStation.Name)

		d.Close()

		d.Play(selectedStation.URL)
		current = selected
	case "p":
		d.Pause()
	case "s":
		d.Close()
		current = -1
	case "k", "<Up>":
		stationsList.ScrollUp()
	case "j", "<Down>":
		stationsList.ScrollDown()
	case "K", "<PageUp>":
		stationsList.ScrollPageUp()
	case "J", "<PageDown>":
		stationsList.ScrollPageDown()
	case "h", "<Left>":
		loggerList.ScrollPageUp()
	case "l", "<Right>":
		loggerList.ScrollPageDown()
	case "+":
		d.IncVolume()
	case "-":
		d.DecVolume()
	}

	render()
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
