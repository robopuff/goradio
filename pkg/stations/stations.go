package stations

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const stationsCsvURL = "https://raw.githubusercontent.com/coderholic/pyradio/master/pyradio/stations.csv"

// Station station struct
type Station struct {
	Name string
	URL  string
}

// StationsList list of available stations
type StationsList struct {
	stations []*Station
}

// GetRows get list of stations for ui list
func (l *StationsList) GetRows() []string {
	list := []string{}
	for _, s := range l.stations {
		list = append(list, s.Name)
	}
	return list
}

// GetSelected get selected station by it's index
func (l *StationsList) GetSelected(selected int) *Station {
	if selected < 0 || selected > len(l.stations)-1 {
		return nil
	}

	return l.stations[selected]
}

// Load load stations from file
func Load(path string) *StationsList {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = downloadStations(path)
			if err != nil {
				log.Fatalf("Cannot find or install stations ")
			}
		}
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("cannot open stations file: %v", err.Error())
	}
	defer f.Close()

	s := []*Station{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\n\r")
		pair := strings.Split(line, ",")
		if len(pair) == 2 {
			s = append(s, &Station{
				Name: strings.TrimSpace(pair[0]),
				URL:  strings.TrimSpace(pair[1]),
			})
		}
	}

	return &StationsList{
		stations: s,
	}
}

func downloadStations(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Chdir(dir)
		if err != nil {
			return err
		}
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(stationsCsvURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
