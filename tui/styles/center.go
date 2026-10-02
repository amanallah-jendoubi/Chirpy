package styles

import (
	"charm.land/lipgloss/v2"
	"strings"
)

// innerW is the width of the box content (the help line is the widest).
const innerW = 36

func Center(s string) string {
	left := (innerW - lipgloss.Width(s)) / 2
	right := innerW - lipgloss.Width(s) - left
	return strings.Repeat(BlankCell, left) + s + strings.Repeat(BlankCell, right)
}
