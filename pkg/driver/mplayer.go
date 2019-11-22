package driver

import (
	"io"
	"log"
	"os/exec"
	"strings"
)

// MPlayer mplayer driver
type MPlayer struct {
	name      string
	isPlaying bool
	streamURL string
	command   *exec.Cmd
	in        io.WriteCloser
	out       io.ReadCloser
	pipeChan  chan io.ReadCloser
}

// PipeChan gets pipe channel
func (driver *MPlayer) PipeChan() chan io.ReadCloser {
	return driver.pipeChan
}

// Play play provided url
func (driver *MPlayer) Play(url string) {
	if driver.isPlaying {
		if driver.streamURL == url {
			return
		}

		driver.Close()
	}

	var err error
	if strings.HasSuffix(url, ".m3u") || strings.HasSuffix(url, ".pls") {
		driver.command = exec.Command(driver.name, "-quiet", "-playlist", url)
	} else {
		driver.command = exec.Command(driver.name, "-quiet", url)
	}

	driver.in, err = driver.command.StdinPipe()
	if nil != err {
		log.Fatalf("cannot map mplayer stdin: %v", err)
	}

	driver.out, err = driver.command.StdoutPipe()
	if nil != err {
		log.Fatalf("cannot map mplayer stdout: %v", err)
	}

	if err = driver.command.Start(); nil != err {
		log.Fatalf("cannot start mplayer: %v", err)
	}

	driver.isPlaying = true
	driver.streamURL = url

	go func() {
		driver.pipeChan <- driver.out
	}()
}

// Close close mplayer
func (driver *MPlayer) Close() {
	if !driver.isPlaying {
		return
	}

	driver.isPlaying = false

	driver.sendKey("q")
	driver.in.Close()
	driver.out.Close()
	driver.command.Process.Kill()

	driver.command = nil
	driver.streamURL = ""
}

// Mute send mute command
func (driver *MPlayer) Mute() {
	driver.sendKey("m")
}

// Pause send pause command
func (driver *MPlayer) Pause() {
	driver.sendKey("p")
}

// IncVolume send increment volume command
func (driver *MPlayer) IncVolume() {
	driver.sendKey("*")
}

// DecVolume send decrease volume command
func (driver *MPlayer) DecVolume() {
	driver.sendKey("/")
}

// sendKey send key directly to mplayer
func (driver *MPlayer) sendKey(key string) {
	if !driver.isPlaying {
		return
	}

	driver.in.Write([]byte(key))
}

// NewMPlayer create new MPlayer driver instance
func NewMPlayer(executablePath string) *MPlayer {
	return &MPlayer{
		name:      executablePath,
		isPlaying: false,
		pipeChan:  make(chan io.ReadCloser),
	}
}
