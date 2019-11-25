package driver

import "io"

// Driver music playing driver interface
type Driver interface {
	PipeChan() chan io.ReadCloser
	CheckPrerequisites() error
	Play(url string)
	Mute()
	Pause()
	IncVolume()
	DecVolume()
	Close()
}
