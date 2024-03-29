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

// NewVolume create new volume custom widget
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

func (v *volume) Draw(buf *tui.Buffer) {
	v.Block.Draw(buf)

	if !v.Visible {
		return
	}

	label := strconv.Itoa(v.Percent)

	barXCoordinate := v.Inner.Min.X
	barYCoordinate := v.Inner.Min.Y + ((v.Inner.Dy() - 1) / 2)
	barDxCoordinate := v.Inner.Max.X - barXCoordinate - (len(label) + 1)

	// plot bar
	barWidth := int((float64(v.Percent) / 100) * float64(barDxCoordinate))
	barStyle := tui.NewStyle(v.BarColors[0], tui.ColorClear)
	for i, char := range v.characters {
		if i > barWidth {
			break
		}

		if c, ok := v.BarColors[i]; ok {
			barStyle = tui.NewStyle(c, tui.ColorClear)
		}

		buf.SetCell(tui.NewCell(char, barStyle), image.Pt(barXCoordinate+i, barYCoordinate))
	}

	// plot label
	labelXCoordinate := v.Inner.Max.X - len(label)
	labelYCoordinate := v.Inner.Min.Y + ((v.Inner.Dy() - 1) / 2)
	if labelYCoordinate < v.Inner.Max.Y {
		style := v.LabelStyle
		for i, char := range label {
			buf.SetCell(tui.NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
		}
	}
}
