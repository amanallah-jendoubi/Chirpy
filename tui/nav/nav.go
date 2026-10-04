package nav

import (
	tea "charm.land/bubbletea/v2"
	"log"
)

// Every screen implements this.
type Screen interface {
	Update(msg tea.Msg) (Screen, tea.Cmd)
	Box(frame int) string // the foreground content drawn over the rain
	Init() tea.Cmd
}

type Route string

const (
	Menu   Route = "menu"
	Login  Route = "login"
	Signup Route = "signup"
	Home   Route = "home"
)

// Screens send this to switch to another screen.
type NavigateMsg struct{ To Route }

func Navigate(to Route) tea.Cmd {
	return func() tea.Msg {
		log.Println("navigate to", to)
		return NavigateMsg{To: to}
	}
}
