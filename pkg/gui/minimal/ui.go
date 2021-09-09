package minimal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/robopuff/goradio/pkg/drivers"
	"github.com/robopuff/goradio/pkg/stations"
)

const (
	regexTitle  = `(?m)^ICY Info: StreamTitle='(.*?)';`
)

type ui struct {
	stationsList *stations.List
	debug        bool
}

func NewMinimal(list *stations.List, debug bool) *ui {
	return &ui{
		stationsList: list,
		debug:        debug,
	}
}

func (u *ui) Init() error {
	return nil
}

func (u *ui) Run(d drivers.Driver) {
	fmt.Println("List of stations:")
	for i, v := range u.stationsList.GetRows(0) {
		fmt.Printf(" [%v] %v\n", i+1, v)
	}

	var err error
	var selected *stations.Station
	for {
		selected, err = u.selectRadio()
		if err != nil {
			fmt.Printf("Invalid selection: %s\n", err)
			continue
		}
		break
	}

	fmt.Printf("Selected %v\n", selected.Name)

	titleRegex := regexp.MustCompile(regexTitle)
	d.Play(selected.URL)
	for {
		select {
		case outPipe := <-d.PipeChan():
			reader := bufio.NewReader(outPipe)
			for {
				data, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				if u.debug {
					fmt.Print(data)
					continue
				}

				match := titleRegex.FindStringSubmatch(data)
				if len(match) > 0 {
					u.setCurrentlyPlaying(match[1])
				}
			}
		}
	}
}

func (u *ui) Close() {
	return
}

func (u *ui) setCurrentlyPlaying(msg string) {
	fmt.Printf("Currently playing: %v\n", msg)
}

func (u *ui) selectRadio() (*stations.Station, error) {
	fmt.Printf("Your selection [1-%d]: ", u.stationsList.Count())

	reader := bufio.NewReader(os.Stdin)
	r, _ := reader.ReadString('\n')
	r = strings.Trim(r, "\n")

	ri, err := strconv.Atoi(r)
	if err != nil {
		return nil, err
	}

	selected := u.stationsList.GetSelected(ri-1)
	if selected == nil {
		return nil, errors.New("invalid station selected")
	}
	return selected, nil
}
