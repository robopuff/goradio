# Goradio

An golang kind of implementation of pyradio.
To properly use this application you will need a _**Golang >= 1.13**_ and _**mplayer**_ installed

## Usage

`$ goradio`

```
$ goradio --help
Usage of goradio:
  -d	Debug mode (shows logger window)
  -m string
    	MPlayer executable (default "mplayer")
  -s string
    	Stations file path (default "~/.config/pyradio/stations.csv")
```


## Installation

1. `git clone https://github.com/robopuff/goradio`
2. `go get`
3. Use one of:
    * `make run` - run an application without really building an executable (always in debug mode)
    * `make build` - build an application (saved into _build/goradio_)
    * `make install` - installs an application in go bin directory

## External libraries

* [term-ui](https://github.com/gizak/termui)

## License

No. Do wathever.
