package styles

import (
	tea "charm.land/bubbletea/v2"
	"math/rand"
	"strings"
	"time"
)

// rain styling

type Column struct {
	Head   float64
	Speed  float64
	Length int
}

type Rain struct {
	Cols []Column
}

func NewColumn(h int, initial bool) Column {
	c := Column{Speed: 0.3 + rand.Float64()*0.9, Length: 6 + rand.Intn(16)}
	if initial {
		c.Head = rand.Float64()*float64(h+c.Length) - float64(c.Length)
	} else {
		c.Head = -rand.Float64() * float64(h)
	}
	return c
}

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(70*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg(t) })
}

func (rain Rain) cell(r, c, frame int) string {
	if c >= len(rain.Cols) {
		return BlankCell
	}
	col := rain.Cols[c]
	d := int(col.Head) - r
	if d < 0 || d >= col.Length {
		return BlankCell
	}
	g := string(Glyphs[(r*131+c*31+frame/(3+c%4))%len(Glyphs)])
	switch {
	case d == 0:
		return RainHead.Render(g)
	case d < 3:
		return RainBright.Render(g)
	case d < col.Length/2:
		return RainMid.Render(g)
	default:
		return RainDim.Render(g)
	}
}

func (rain Rain) RainRow(r, from, to, frame int) string {
	var sb strings.Builder
	for c := from; c < to; c++ {
		sb.WriteString(rain.cell(r, c, frame))
	}
	return sb.String()
}
