package login

import (
	"bytes"
	"encoding/json"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/amanallah-jendoubi/Textio/tui/auth"
	"github.com/amanallah-jendoubi/Textio/tui/client"
	apperr "github.com/amanallah-jendoubi/Textio/tui/errors"
	"github.com/amanallah-jendoubi/Textio/tui/nav"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
)

// errMsg is private: only this screen sends and handles it.

type field struct {
	label  string
	value  []rune
	secret bool
}

type login struct {
	fields  []field
	focus   int // 0..len(fields)-1 = inputs, len(fields) = button
	err     string
	loading bool
}

func NewLogin() *login {
	return &login{
		fields: []field{
			{label: "Username"},
			{label: "Password", secret: true},
		},
	}
}

func (l *login) Init() tea.Cmd { return nil }

func (l *login) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case apperr.LoginErrMsg:
		l.loading, l.err = false, msg.UserErr
	case apperr.AuthErrMsg:
		l.err = msg.UserErr
	case tea.KeyPressMsg:
		if l.loading {
			return l, nil // ignore input while waiting for the server
		}
		n := len(l.fields)
		switch msg.String() {
		case "esc":
			return l, nav.Navigate(nav.Menu)
		case "tab", "down":
			l.focus = (l.focus + 1) % (n + 1)
		case "shift+tab", "up":
			l.focus = (l.focus + n) % (n + 1)
		case "enter":
			if l.focus < n {
				l.focus++
				return l, nil
			}
			return l, l.submit()
		case "backspace":
			if l.focus < n {
				v := l.fields[l.focus].value
				if len(v) > 0 {
					l.fields[l.focus].value = v[:len(v)-1]
				}
			}
		default:
			// printable text (letters, digits, symbols, space)
			if l.focus < n && msg.Text != "" {
				l.fields[l.focus].value = append(l.fields[l.focus].value, []rune(msg.Text)...)
			}
		}
	}
	return l, nil
}

// submit validates the form and returns the HTTP request as a Cmd.
func (l *login) submit() tea.Cmd {
	user := strings.TrimSpace(string(l.fields[0].value))
	pass := string(l.fields[1].value)
	if user == "" || pass == "" {
		l.err = "All fields are required"
		return nil
	}
	l.err, l.loading = "", true
	return loginCmd(user, pass)
}

func loginCmd(username, password string) tea.Cmd {
	return func() tea.Msg {
		body, err := json.Marshal(map[string]string{"name": username, "password": password})
		if err != nil {
			return apperr.HandleLocalError(apperr.AsLogin, err)
		}
		res, err := client.Post("/login", bytes.NewReader(body))
		if err != nil {
			return apperr.HandleLocalError(apperr.AsLogin, err)
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return apperr.HandleAPIError(apperr.AsLogin, res, "")
		}
		err = auth.StoreTokens(res)
		if err != nil {
			return apperr.HandleLocalError(apperr.AsAuth, err)
		}
		return nav.Navigate(nav.Home)()
	}
}

func (l *login) input(i int, frame int) string {
	f := l.fields[i]
	text := string(f.value)
	if f.secret {
		text = strings.Repeat("•", len(f.value))
	}
	r := []rune(text)
	const maxText = styles.FieldW - 3 // 2 prefix + 1 cursor
	if len(r) > maxText {
		r = r[len(r)-maxText:]
	}
	if i == l.focus {
		cur := " "
		if (frame/6)%2 == 0 {
			cur = "█"
		}
		return styles.InputActive.Render("▌ " + string(r) + cur)
	}
	return styles.InputIdle.Render("  " + string(r))
}

func (l *login) Box(frame int) string {
	btn := "Login"
	if l.loading {
		btn = "Connecting..."
	}

	rows := []string{styles.Center(styles.BannerStyle.Render("LOGIN")), styles.Center("")}
	for i, f := range l.fields {
		rows = append(rows,
			styles.Center(styles.LabelStyle.Render(f.label)),
			styles.Center(l.input(i, frame)),
			styles.Center(""),
		)
	}

	if l.focus == len(l.fields) || l.loading {
		rows = append(rows, styles.Center(styles.ItemSelected.Render(btn)))
	} else {
		rows = append(rows, styles.Center(styles.ItemStyle.Render(btn)))
	}

	if l.err != "" {
		rows = append(rows, styles.Center(styles.ErrStyle.Render(styles.Clip(l.err, styles.FieldW))))
	} else {
		rows = append(rows, styles.Center(""))
	}
	rows = append(rows, styles.Center(styles.HelpStyle.Render("↑/↓ move • enter next • esc back")))

	return styles.BoxStyle.Render(strings.Join(rows, "\n"))
}
