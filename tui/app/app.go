package app

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/amanallah-jendoubi/Textio/tui/nav"
	"github.com/amanallah-jendoubi/Textio/tui/pages/login"
	"github.com/amanallah-jendoubi/Textio/tui/pages/menu"
	"github.com/amanallah-jendoubi/Textio/tui/pages/signup"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
	"strings"
)

// screen builder
func build(r nav.Route) nav.Screen {
	switch r {
	case nav.Login:
		return login.NewLogin()
	case nav.Signup:
		return signup.NewSignup()
	case nav.Menu:
		return menu.NewMenu()
	default:
		return nil
	}
}

type app struct {
	w, h    int
	frame   int
	rain    styles.Rain
	current nav.Screen
}

func NewApp() app { return app{current: menu.NewMenu()} }

func (a app) Init() tea.Cmd { return styles.Tick() }

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.w, a.h = msg.Width, msg.Height
		a.rain.Cols = make([]styles.Column, a.w)
		for i := range a.rain.Cols {
			a.rain.Cols[i] = styles.NewColumn(a.h, true)
		}
		return a, nil

	case styles.TickMsg:
		a.frame++
		for i := range a.rain.Cols {
			c := &a.rain.Cols[i]
			c.Head += c.Speed
			if int(c.Head)-c.Length > a.h {
				*c = styles.NewColumn(a.h, false)
			}
		}
		return a, styles.Tick()

	case nav.NavigateMsg:
		a.current = build(msg.To)
		return a, nil

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" { // global quit only
			return a, tea.Quit
		}
	}

	// everything else goes to the active screen
	var cmd tea.Cmd
	a.current, cmd = a.current.Update(msg)
	return a, cmd
}

func (a app) View() tea.View {
	v := tea.NewView(a.render())
	v.AltScreen = true
	return v
}

func (a app) render() string {
	if a.w == 0 || a.h == 0 {
		return ""
	}
	box := a.current.Box(a.frame)
	boxLines := strings.Split(box, "\n")
	bw, bh := lipgloss.Width(box), len(boxLines)
	if bw > a.w || bh > a.h {
		return box
	}
	x0, y0 := (a.w-bw)/2, (a.h-bh)/2
	lines := make([]string, a.h)
	for r := 0; r < a.h; r++ {
		if r >= y0 && r < y0+bh {
			lines[r] = a.rain.RainRow(r, 0, x0, a.frame) + boxLines[r-y0] + a.rain.RainRow(r, x0+bw, a.w, a.frame)
		} else {
			lines[r] = a.rain.RainRow(r, 0, a.w, a.frame)
		}
	}
	return strings.Join(lines, "\n")
}
