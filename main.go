package main

import (
	"flag"
	"fmt"
	"github.com/robopuff/goradio/pkg/gui/minimal"
	"log"
	"os/user"

	"github.com/robopuff/goradio/pkg/drivers"
	"github.com/robopuff/goradio/pkg/gui"
	"github.com/robopuff/goradio/pkg/gui/termui"
	"github.com/robopuff/goradio/pkg/stations"
)

func main() {
	usr, _ := user.Current()
	flagStations := flag.String("s", fmt.Sprintf("%s/.config/pyradio/stations.csv", usr.HomeDir), "Stations file path")
	flagMplayer := flag.String("m", "mplayer", "MPlayer executable")
	flagDebugMode := flag.Bool("d", false, "Debug mode (shows logger window)")
	flagMinimal := flag.Bool("minimal", false, "Use minimal ui")
	flag.Parse()

	d := drivers.NewMPlayer(*flagMplayer)
	if err := d.CheckPrerequisites(); err != nil {
		log.Fatalf("system failed drivers prerequisites check: %v", err)
	}

	s, err := stations.Load(*flagStations)
	if err != nil {
		log.Fatalf("cannot load or download stations list: %v", err)
	}

	var userInterface gui.GUI
	userInterface = termui.NewTermUI(s, *flagDebugMode)
	if *flagMinimal {
		userInterface = minimal.NewMinimal(s, *flagDebugMode)
	}

	defer (func() {
		d.Close()
		userInterface.Close()
	})()

	if err := userInterface.Init(); err != nil {
		log.Fatalf("failed to initialize gui: %v", err)
	}
	userInterface.Run(d)
}
