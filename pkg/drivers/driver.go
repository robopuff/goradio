package drivers

import "io"

// Driver music playing drivers interface
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
