package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ---------- Theme ----------

var (
	colBg     = lipgloss.Color("#0D0208")
	colDim    = lipgloss.Color("#003B00")
	colMid    = lipgloss.Color("#008F11")
	colBright = lipgloss.Color("#00FF41")
	colWhite  = lipgloss.Color("#C8FFD4")
	colError  = lipgloss.Color("#FF4D4D")
)

var (
	rainHead   = lipgloss.NewStyle().Foreground(colWhite).Background(colBg).Bold(true)
	rainBright = lipgloss.NewStyle().Foreground(colBright).Background(colBg)
	rainMid    = lipgloss.NewStyle().Foreground(colMid).Background(colBg)
	rainDim    = lipgloss.NewStyle().Foreground(colDim).Background(colBg)
	blankCell  = lipgloss.NewStyle().Background(colBg).Render(" ")

	bannerStyle   = lipgloss.NewStyle().Foreground(colBright).Background(colBg).Bold(true)
	subtitleStyle = lipgloss.NewStyle().Foreground(colMid).Background(colBg).Width(30).Align(lipgloss.Center)
	itemStyle     = lipgloss.NewStyle().Foreground(colMid).Background(colBg).Width(18).Align(lipgloss.Center)
	itemSelected  = lipgloss.NewStyle().Foreground(colBg).Background(colBright).Bold(true).Width(18).Align(lipgloss.Center)
	helpStyle     = lipgloss.NewStyle().Foreground(colDim).Background(colBg)
	labelStyle    = lipgloss.NewStyle().Foreground(colMid).Background(colBg).Width(fieldW)
	inputActive   = lipgloss.NewStyle().Foreground(colWhite).Background(colBg).Bold(true).Width(fieldW)
	inputIdle     = lipgloss.NewStyle().Foreground(colMid).Background(colBg).Width(fieldW)
	errStyle      = lipgloss.NewStyle().Foreground(colError).Background(colBg)
	textStyle     = lipgloss.NewStyle().Foreground(colMid).Background(colBg)
	boxStyle      = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colBright).
			BorderBackground(colBg).
			Background(colBg).
			Padding(1, 4)
)

var glyphs = []rune("ｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ0123456789:.=*+-<>¦")

const banner = "╔╦╗╔═╗╔╦╗╦═╗╦═╗ ╦\n║║║╠═╣ ║ ╠╦╝║╔╩╦╝\n╩ ╩╩ ╩ ╩ ╩╚═╩╩ ╚═"

const tagline = "Wake up, Neo..."

const (
	innerW = 36 // width of the box content
	fieldW = 32 // width of a form field
)

// ---------- Backend hooks (TODO: wire to your server) ----------

func login(username, password string) error {
	// TODO: POST to your backend. Return an error to show it in the form.
	return nil
}

func signup(username, password string) error {
	// TODO: POST to your backend. Return an error to show it in the form.
	return nil
}

// ---------- Screens & menu ----------

type screen int

const (
	screenWelcome screen = iota
	screenLogin
	screenSignup
	screenHome
)

type menuItem struct {
	id    string
	label string
}

var items = []menuItem{
	{"login", "Login"},
	{"signup", "Sign up"},
	{"quit", "Quit"},
}

type field struct {
	label  string
	value  []rune
	secret bool
}

func loginFields() []field {
	return []field{
		{label: "Username"},
		{label: "Password", secret: true},
	}
}

func signupFields() []field {
	return []field{
		{label: "Username"},
		{label: "Password", secret: true},
		{label: "Confirm password", secret: true},
	}
}

// ---------- Rain ----------

type column struct {
	head   float64
	speed  float64
	length int
}

func newColumn(h int, initial bool) column {
	c := column{speed: 0.3 + rand.Float64()*0.9, length: 6 + rand.Intn(16)}
	if initial {
		c.head = rand.Float64()*float64(h+c.length) - float64(c.length)
	} else {
		c.head = -rand.Float64() * float64(h)
	}
	return c
}

// ---------- Model ----------

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(70*time.Millisecond, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type model struct {
	w, h   int
	cols   []column
	frame  int
	screen screen

	cursor int // welcome menu

	fields []field // login / signup form
	focus  int     // 0..len(fields)-1 = inputs, len(fields) = submit button
	errMsg string

	user string // logged-in username
}

func (m model) Init() tea.Cmd { return tick() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.cols = make([]column, m.w)
		for i := range m.cols {
			m.cols[i] = newColumn(m.h, true)
		}

	case tickMsg:
		m.frame++
		for i := range m.cols {
			c := &m.cols[i]
			c.head += c.speed
			if int(c.head)-c.length > m.h {
				*c = newColumn(m.h, false)
			}
		}
		return m, tick()

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch m.screen {
		case screenWelcome:
			return m.updateWelcome(msg)
		case screenLogin, screenSignup:
			return m.updateForm(msg)
		default:
			if s := msg.String(); s == "q" || s == "esc" {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) updateWelcome(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit
	case "up", "k":
		m.cursor = (m.cursor + len(items) - 1) % len(items)
	case "down", "j", "tab":
		m.cursor = (m.cursor + 1) % len(items)
	case "enter":
		switch items[m.cursor].id {
		case "login":
			return m.withForm(screenLogin), nil
		case "signup":
			return m.withForm(screenSignup), nil
		default:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) withForm(s screen) model {
	m.screen, m.focus, m.errMsg = s, 0, ""
	if s == screenLogin {
		m.fields = loginFields()
	} else {
		m.fields = signupFields()
	}
	return m
}

func (m model) updateForm(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	n := len(m.fields)
	switch msg.String() {
	case "esc":
		m.screen, m.fields, m.errMsg = screenWelcome, nil, ""
		return m, nil
	case "tab", "down":
		m.focus = (m.focus + 1) % (n + 1)
	case "shift+tab", "up":
		m.focus = (m.focus + n) % (n + 1)
	case "enter":
		if m.focus < n {
			m.focus++
			return m, nil
		}
		return m.submit(), nil
	case "backspace":
		if m.focus < n {
			v := m.fields[m.focus].value
			if len(v) > 0 {
				m.fields[m.focus].value = v[:len(v)-1]
			}
		}
	default:
		// printable text (letters, digits, symbols, space)
		if m.focus < n && msg.Text != "" {
			m.fields[m.focus].value = append(m.fields[m.focus].value, []rune(msg.Text)...)
		}
	}
	return m, nil
}

func (m model) submit() model {
	vals := make([]string, len(m.fields))
	for i, f := range m.fields {
		vals[i] = string(f.value)
		if !f.secret {
			vals[i] = strings.TrimSpace(vals[i])
		}
	}
	for _, v := range vals {
		if v == "" {
			m.errMsg = "All fields are required"
			return m
		}
	}

	var err error
	if m.screen == screenLogin {
		err = login(vals[0], vals[1])
	} else {
		switch {
		case len([]rune(vals[1])) < 8:
			m.errMsg = "Password must be 8+ characters"
			return m
		case vals[1] != vals[2]:
			m.errMsg = "Passwords do not match"
			return m
		}
		err = signup(vals[0], vals[1])
	}
	if err != nil {
		m.errMsg = err.Error()
		return m
	}

	m.user, m.screen, m.fields, m.errMsg = vals[0], screenHome, nil, ""
	return m
}

// ---------- View ----------

func (m model) cell(r, c int) string {
	if c >= len(m.cols) {
		return blankCell
	}
	col := m.cols[c]
	d := int(col.head) - r
	if d < 0 || d >= col.length {
		return blankCell
	}
	g := string(glyphs[(r*131+c*31+m.frame/(3+c%4))%len(glyphs)])
	switch {
	case d == 0:
		return rainHead.Render(g)
	case d < 3:
		return rainBright.Render(g)
	case d < col.length/2:
		return rainMid.Render(g)
	default:
		return rainDim.Render(g)
	}
}

func (m model) rainRow(r, from, to int) string {
	var sb strings.Builder
	for c := from; c < to; c++ {
		sb.WriteString(m.cell(r, c))
	}
	return sb.String()
}

// center pads s to innerW with background-colored spaces, so no unstyled
// gaps show the terminal's own background.
func center(s string) string {
	left := (innerW - lipgloss.Width(s)) / 2
	right := innerW - lipgloss.Width(s) - left
	return strings.Repeat(blankCell, left) + s + strings.Repeat(blankCell, right)
}

func (m model) welcomeRows() []string {
	// typewriter subtitle with blinking cursor
	runes := []rune(tagline)
	n := m.frame / 3
	if n > len(runes) {
		n = len(runes)
	}
	cur := " "
	if (m.frame/6)%2 == 0 {
		cur = "█"
	}
	sub := subtitleStyle.Render(string(runes[:n]) + cur)

	var rows []string
	for _, l := range strings.Split(banner, "\n") {
		rows = append(rows, center(bannerStyle.Render(l)))
	}
	rows = append(rows, center(""), center(sub), center(""))
	for i, it := range items {
		if i == m.cursor {
			rows = append(rows, center(itemSelected.Render(it.label)))
		} else {
			rows = append(rows, center(itemStyle.Render(it.label)))
		}
	}
	rows = append(rows, center(""), center(helpStyle.Render("↑/↓ navigate • enter select • q quit")))
	return rows
}

func (m model) inputView(i int, f field) string {
	text := string(f.value)
	if f.secret {
		text = strings.Repeat("•", len(f.value))
	}
	r := []rune(text)
	const maxText = fieldW - 3 // 2 for the prefix, 1 for the cursor
	if len(r) > maxText {
		r = r[len(r)-maxText:]
	}
	if i == m.focus {
		cur := " "
		if (m.frame/6)%2 == 0 {
			cur = "█"
		}
		return inputActive.Render("▌ " + string(r) + cur)
	}
	return inputIdle.Render("  " + string(r))
}

func (m model) formRows() []string {
	title, btn := "LOGIN", "Login"
	if m.screen == screenSignup {
		title, btn = "SIGN UP", "Create account"
	}

	rows := []string{center(bannerStyle.Render(title)), center("")}
	for i, f := range m.fields {
		rows = append(rows, center(labelStyle.Render(f.label)), center(m.inputView(i, f)), center(""))
	}

	if m.focus == len(m.fields) {
		rows = append(rows, center(itemSelected.Render(btn)))
	} else {
		rows = append(rows, center(itemStyle.Render(btn)))
	}

	if m.errMsg != "" {
		rows = append(rows, center(errStyle.Render(m.errMsg)))
	} else {
		rows = append(rows, center(""))
	}
	rows = append(rows, center(helpStyle.Render("↑/↓ move • enter next • esc back")))
	return rows
}

func (m model) homeRows() []string {
	return []string{
		center(bannerStyle.Render("ACCESS GRANTED")),
		center(""),
		center(textStyle.Render("Welcome, " + m.user)),
		center(textStyle.Render("Conversations screen: next step")),
		center(""),
		center(helpStyle.Render("q quit")),
	}
}

func (m model) box() string {
	var rows []string
	switch m.screen {
	case screenLogin, screenSignup:
		rows = m.formRows()
	case screenHome:
		rows = m.homeRows()
	default:
		rows = m.welcomeRows()
	}
	return boxStyle.Render(strings.Join(rows, "\n"))
}

func (m model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

func (m model) render() string {
	if m.w == 0 || m.h == 0 {
		return ""
	}
	box := m.box()
	boxLines := strings.Split(box, "\n")
	bw, bh := lipgloss.Width(box), len(boxLines)

	// terminal too small: show the box only
	if bw > m.w || bh > m.h {
		return box
	}

	x0, y0 := (m.w-bw)/2, (m.h-bh)/2
	lines := make([]string, m.h)
	for r := 0; r < m.h; r++ {
		if r >= y0 && r < y0+bh {
			lines[r] = m.rainRow(r, 0, x0) + boxLines[r-y0] + m.rainRow(r, x0+bw, m.w)
		} else {
			lines[r] = m.rainRow(r, 0, m.w)
		}
	}
	return strings.Join(lines, "\n")
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
