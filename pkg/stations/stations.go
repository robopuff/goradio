package stations

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const stationsCsvURL = "https://raw.githubusercontent.com/coderholic/pyradio/master/pyradio/stations.csv"

type station struct {
	Name string
	URL  string
}

// List list of available stations
type List struct {
	path     string
	stations []*station
}

// GetRows get list of stations for gui list
func (l *List) GetRows(width int) []string {
	formatter := fmt.Sprintf("%%-%ds", width-2)
	var list []string
	for _, s := range l.stations {
		list = append(list, fmt.Sprintf(formatter, s.Name))
	}
	return list
}

// GetSelected get selected station by it's index
func (l *List) GetSelected(selected int) *station {
	if selected < 0 || selected > len(l.stations)-1 {
		return nil
	}

	return l.stations[selected]
}

// Reload reloads list
func (l *List) Reload() {
	load, err := Load(l.path)
	if err != nil {
		log.Fatalf("cannot refresh stations list: %v", err)
	}

	l.stations = load.stations
}

// Load load stations from file
func Load(path string) (*List, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = downloadStations(path)
			if err != nil {
				return nil, err
			}
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var s []*station
	reader := csv.NewReader(f)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		s = append(s, &station{
			Name: strings.TrimSpace(line[0]),
			URL:  strings.TrimSpace(line[1]),
		})
	}

	return &List{
		path:     path,
		stations: s,
	}, nil
}

func downloadStations(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("pyradio config dir not found, trying to create")
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		log.Printf("created directory: %v", dir)
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
