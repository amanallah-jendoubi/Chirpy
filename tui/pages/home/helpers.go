package home

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
)

func formatMessageTime(at time.Time) string {
	months := [...]string{
		"janvier", "février", "mars", "avril", "mai", "juin",
		"juillet", "août", "septembre", "octobre", "novembre", "décembre",
	}
	return fmt.Sprintf("%d %s %d, %02d:%02d", at.Day(), months[int(at.Month())-1], at.Year(), at.Hour(), at.Minute())
}

// wrap breaks s into lines of at most w runes, splitting on spaces.
func wrap(s string, w int) []string {
	var lines []string
	line := ""
	for _, word := range strings.Fields(s) {
		r := []rune(word)
		for len(r) > w { // a single word longer than a line
			if line != "" {
				lines = append(lines, line)
				line = ""
			}
			lines = append(lines, string(r[:w]))
			r = r[w:]
		}
		if len(r) == 0 {
			continue
		}
		word = string(r)
		switch {
		case line == "":
			line = word
		case len([]rune(line))+1+len(r) <= w:
			line += " " + word
		default:
			lines = append(lines, line)
			line = word
		}
	}
	if line != "" || len(lines) == 0 {
		lines = append(lines, line)
	}
	return lines
}

// pad fills s with background-colored spaces up to w columns.
func pad(s string, w int) string {
	if gap := w - lipgloss.Width(s); gap > 0 {
		return s + strings.Repeat(styles.BlankCell, gap)
	}
	return s
}

// padPlain pads with plain spaces, for text that is styled afterwards.
func padPlain(s string, w int) string {
	if n := len([]rune(s)); n < w {
		return s + strings.Repeat(" ", w-n)
	}
	return s
}

// renderMsg wraps one message to width w and styles it: time, sender, text.
func renderMsg(m message, w int) []string {
	who := styles.RainMid
	if m.from == "me" {
		who = styles.RainHead
	}
	prefix := m.at + " " + m.from + ": "
	pl, tl := len([]rune(prefix)), len([]rune(m.at))

	lines := wrap(prefix+m.body, w)
	out := make([]string, len(lines))
	for i, ln := range lines {
		r := []rune(ln)
		if i == 0 && len(r) >= pl {
			out[i] = styles.RainDim.Render(string(r[:tl])) + styles.RainMid.Render(" ") +
				who.Render(string(r[tl+1:pl])) + styles.RainMid.Render(string(r[pl:]))
		} else {
			out[i] = styles.RainMid.Render(ln)
		}
	}
	return out
}

// help panel
func (h *home) help() string {
	switch h.mode {
	case modeChat:
		return "enter send • esc back"
	case modeGroup:
		return "tab/↑↓ move • enter/space toggle • esc cancel"
	case modeDM:
		return "enter start • esc cancel"
	default:
		return "↑/↓ select • c new chat • n new group • q quit"
	}
}
