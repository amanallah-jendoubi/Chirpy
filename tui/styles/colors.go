package styles

import (
	"charm.land/lipgloss/v2"
)

var (
	colBg     = lipgloss.Color("#0D0208")
	colDim    = lipgloss.Color("#003B00")
	colMid    = lipgloss.Color("#008F11")
	colBright = lipgloss.Color("#00FF41")
	colWhite  = lipgloss.Color("#C8FFD4")
)

var (
	RainHead   = lipgloss.NewStyle().Foreground(colWhite).Background(colBg).Bold(true)
	RainBright = lipgloss.NewStyle().Foreground(colBright).Background(colBg)
	RainMid    = lipgloss.NewStyle().Foreground(colMid).Background(colBg)
	RainDim    = lipgloss.NewStyle().Foreground(colDim).Background(colBg)
	BlankCell  = lipgloss.NewStyle().Background(colBg).Render(" ")

	BannerStyle   = lipgloss.NewStyle().Foreground(colBright).Background(colBg).Bold(true)
	SubtitleStyle = lipgloss.NewStyle().Foreground(colMid).Background(colBg).Width(30).Align(lipgloss.Center)
	ItemStyle     = lipgloss.NewStyle().Foreground(colMid).Background(colBg).Width(18).Align(lipgloss.Center)
	ItemSelected  = lipgloss.NewStyle().Foreground(colBg).Background(colBright).Bold(true).Width(18).Align(lipgloss.Center)
	HelpStyle     = lipgloss.NewStyle().Foreground(colDim).Background(colBg)
	BoxStyle      = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colBright).
			BorderBackground(colBg).
			Background(colBg).
			Padding(1, 4)
)

var Glyphs = []rune("ｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ0123456789:.=*+-<>¦")

// FieldW is the width of a form input (the box content is 36 wide).
const FieldW = 32

var (
	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#008F11")).
			Background(lipgloss.Color("#0D0208")).
			Width(FieldW)

	InputActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C8FFD4")).
			Background(lipgloss.Color("#0D0208")).
			Bold(true).
			Width(FieldW)

	InputIdle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#008F11")).
			Background(lipgloss.Color("#0D0208")).
			Width(FieldW)

	ErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4D4D")).
			Background(lipgloss.Color("#0D0208"))
)

// Clip shortens s to at most n runes, ending with "…".
func Clip(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}
