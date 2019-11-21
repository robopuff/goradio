package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"

	"github.com/robopuff/goradio/pkg/driver"
	"github.com/robopuff/goradio/pkg/stations"
	"github.com/robopuff/goradio/pkg/ui"
)

var d driver.Driver

func main() {
	usr, _ := user.Current()
	flagStations := flag.String("s", fmt.Sprintf("%s/.config/pyradio/stations.csv", usr.HomeDir), "Stations file path")
	flagMplayer := flag.String("m", "mplayer", "MPlayer executable")
	flagDebugMode := flag.Bool("d", false, "Debug mode (shows logger window)")
	flag.Parse()

	s := stations.Load(*flagStations)

	if err := ui.Init(s, *flagDebugMode); err != nil {
		log.Fatalf("failed to initialize ui: %v", err)
	}

	d = driver.NewMPlayer(*flagMplayer)
	ui.Run(d)
}
