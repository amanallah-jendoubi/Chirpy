package menu

import (
	tea "charm.land/bubbletea/v2"
	"github.com/amanallah-jendoubi/matrix-chat/tui/nav"
	"github.com/amanallah-jendoubi/matrix-chat/tui/styles"
	"strings"
)

type menuItem struct{ id, label string }

var items = []menuItem{
	{"login", "Login"},
	{"signup", "Sign up"},
	{"quit", "Quit"},
}

type menu struct {
	cursor int
	items  []menuItem
}

func NewMenu() *menu {
	return &menu{
		items: items,
	}
}

func (m *menu) Init() tea.Cmd { return nil }

func (m *menu) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	if k, ok := msg.(tea.KeyPressMsg); ok {
		switch k.String() {
		case "q", "esc":
			return m, tea.Quit
		case "up", "k":
			m.cursor = (m.cursor + len(m.items) - 1) % len(m.items)
		case "down", "j", "tab":
			m.cursor = (m.cursor + 1) % len(m.items)
		case "enter":
			switch m.items[m.cursor].id {
			case "login":
				return m, nav.Navigate(nav.Login)
			case "signup":
				return m, nav.Navigate(nav.Signup)
			case "quit":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m menu) Box(frame int) string {
	// typewriter subtitle with blinking cursor
	runes := []rune(tagline)
	n := frame / 3
	if n > len(runes) {
		n = len(runes)
	}
	cur := " "
	if (frame/6)%2 == 0 {
		cur = "█"
	}
	sub := styles.SubtitleStyle.Render(string(runes[:n]) + cur)

	var rows []string
	for _, l := range strings.Split(banner, "\n") {
		rows = append(rows, styles.Center(styles.BannerStyle.Render(l)))
	}
	rows = append(rows, styles.Center(""), styles.Center(sub), styles.Center(""))
	for i, it := range m.items {
		if i == m.cursor {
			rows = append(rows, styles.Center(styles.ItemSelected.Render(it.label)))
		} else {
			rows = append(rows, styles.Center(styles.ItemStyle.Render(it.label)))
		}
	}
	rows = append(rows, styles.Center(""), styles.Center(styles.HelpStyle.Render("↑/↓ navigate • enter select • q quit")))

	return styles.BoxStyle.Render(strings.Join(rows, "\n"))
}
