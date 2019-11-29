package gui

import "github.com/robopuff/goradio/pkg/drivers"

const (
	CurrentlyPlayingFormat = "Currently playing: %s"
	HelpFooter             = "k/↑ : Up | j/↓: Down | Enter: Select | p: Pause | m: Mute | s: Stop | +: Louder | -: Quieter | R: Refresh | q: Quit"
)

// GUI gui interface
type GUI interface {
	Init() error
	Run(d drivers.Driver)
	Close()
}
