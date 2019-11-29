package termui

import (
	"image"
	"strconv"

	tui "github.com/gizak/termui/v3"
)

type volume struct {
	tui.Block
	Percent    int
	Visible    bool
	BarColors  map[int]tui.Color
	LabelStyle tui.Style
	characters []rune
}

//NewVolume create new volume custom widget
func NewVolume() *volume {
	return &volume{
		Block: *tui.NewBlock(),
		BarColors: map[int]tui.Color{
			0: tui.ColorGreen,
			4: tui.ColorYellow,
			9: tui.ColorRed,
		},
		LabelStyle: tui.Theme.Gauge.Label,
		characters: []rune{
			'▁', '▁', '▂', '▂', '▃',
			'▃', '▅', '▅', '▆', '▆',
			'▇', '▇', '█', '█',
		},
	}
}

func (self *volume) Draw(buf *tui.Buffer) {
	self.Block.Draw(buf)

	if !self.Visible {
		return
	}

	label := strconv.Itoa(self.Percent)

	barXCoordinate := self.Inner.Min.X
	barYCoordinate := self.Inner.Min.Y + ((self.Inner.Dy() - 1) / 2)
	barDxCoordinate := self.Inner.Max.X - barXCoordinate - (len(label) + 1)

	// plot bar
	barWidth := int((float64(self.Percent) / 100) * float64(barDxCoordinate))
	barStyle := tui.NewStyle(self.BarColors[0], tui.ColorClear)
	for i, char := range self.characters {
		if i > barWidth {
			break
		}

		if c, ok := self.BarColors[i]; ok {
			barStyle = tui.NewStyle(c, tui.ColorClear)
		}

		buf.SetCell(tui.NewCell(char, barStyle), image.Pt(barXCoordinate+i, barYCoordinate))
	}

	// plot label
	labelXCoordinate := self.Inner.Max.X - len(label)
	labelYCoordinate := self.Inner.Min.Y + ((self.Inner.Dy() - 1) / 2)
	if labelYCoordinate < self.Inner.Max.Y {
		style := self.LabelStyle
		for i, char := range label {
			buf.SetCell(tui.NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
		}
	}
}
