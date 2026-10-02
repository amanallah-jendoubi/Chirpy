package signup

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/amanallah-jendoubi/Textio/tui/client"
	"github.com/amanallah-jendoubi/Textio/tui/nav"
	"github.com/amanallah-jendoubi/Textio/tui/styles"
)

type field struct {
	label  string
	value  []rune
	secret bool
}

type signup struct {
	fields  []field
	focus   int // 0..len(fields)-1 = inputs, len(fields) = button
	err     string
	loading bool
}

func NewSignup() *signup {
	return &signup{
		fields: []field{
			{label: "Username"},
			{label: "Password", secret: true},
			{label: "Confirm password", secret: true},
		},
	}
}

func (s *signup) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case client.ErrMsg:
		s.loading, s.err = false, msg.Error()

	case tea.KeyPressMsg:
		if s.loading {
			return s, nil // ignore input while waiting for the server
		}
		n := len(s.fields)
		switch msg.String() {
		case "esc":
			return s, nav.Navigate(nav.Menu)
		case "tab", "down":
			s.focus = (s.focus + 1) % (n + 1)
		case "shift+tab", "up":
			s.focus = (s.focus + n) % (n + 1)
		case "enter":
			if s.focus < n {
				s.focus++
				return s, nil
			}
			return s, s.submit()
		case "backspace":
			if s.focus < n {
				v := s.fields[s.focus].value
				if len(v) > 0 {
					s.fields[s.focus].value = v[:len(v)-1]
				}
			}
		default:
			// printable text (letters, digits, symbols, space)
			if s.focus < n && msg.Text != "" {
				s.fields[s.focus].value = append(s.fields[s.focus].value, []rune(msg.Text)...)
			}
		}
	}
	return s, nil
}

// submit validates the form and returns the HTTP request as a Cmd.
func (s *signup) submit() tea.Cmd {
	user := strings.TrimSpace(string(s.fields[0].value))
	pass := string(s.fields[1].value)
	confirm := string(s.fields[2].value)

	switch {
	case user == "" || pass == "" || confirm == "":
		s.err = "All fields are required"
	case len([]rune(pass)) < 8:
		s.err = "Password must be 8+ characters"
	case pass != confirm:
		s.err = "Passwords do not match"
	default:
		s.err, s.loading = "", true
		return signupCmd(user, pass)
	}
	return nil
}

func signupCmd(username, password string) tea.Cmd {
	return func() tea.Msg {
		body, err := json.Marshal(map[string]string{"name": username, "password": password})
		if err != nil {
			return client.ErrMsg{Err: err}
		}
		res, err := client.Post("/register", bytes.NewReader(body))
		if err != nil {
			return client.ErrMsg{Err: errors.New("cannot reach server")}
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode > 299 {
			var payload struct {
				Error string `json:"error"`
			}
			// a body that isn't JSON just falls back to the status text
			_ = json.NewDecoder(io.LimitReader(res.Body, 1024)).Decode(&payload)
			text := strings.TrimSpace(payload.Error)
			if text == "" {
				text = http.StatusText(res.StatusCode)
			}
			return client.ErrMsg{Err: errors.New(text)}
		}
		return nav.Navigate(nav.Home)()
	}
}

func (s *signup) input(i int, frame int) string {
	f := s.fields[i]
	text := string(f.value)
	if f.secret {
		text = strings.Repeat("•", len(f.value))
	}
	r := []rune(text)
	const maxText = styles.FieldW - 3 // 2 prefix + 1 cursor
	if len(r) > maxText {
		r = r[len(r)-maxText:]
	}
	if i == s.focus {
		cur := " "
		if (frame/6)%2 == 0 {
			cur = "█"
		}
		return styles.InputActive.Render("▌ " + string(r) + cur)
	}
	return styles.InputIdle.Render("  " + string(r))
}

func (s *signup) Box(frame int) string {
	btn := "Create account"
	if s.loading {
		btn = "Creating..."
	}

	rows := []string{styles.Center(styles.BannerStyle.Render("SIGN UP")), styles.Center("")}
	for i, f := range s.fields {
		rows = append(rows,
			styles.Center(styles.LabelStyle.Render(f.label)),
			styles.Center(s.input(i, frame)),
			styles.Center(""),
		)
	}

	if s.focus == len(s.fields) || s.loading {
		rows = append(rows, styles.Center(styles.ItemSelected.Render(btn)))
	} else {
		rows = append(rows, styles.Center(styles.ItemStyle.Render(btn)))
	}

	if s.err != "" {
		rows = append(rows, styles.Center(styles.ErrStyle.Render(styles.Clip(s.err, styles.FieldW))))
	} else {
		rows = append(rows, styles.Center(""))
	}
	rows = append(rows, styles.Center(styles.HelpStyle.Render("↑/↓ move • enter next • esc back")))

	return styles.BoxStyle.Render(strings.Join(rows, "\n"))
}
