package home

import (
	"github.com/amanallah-jendoubi/Textio/tui/styles"
)

// ---------- Layout ----------
const (
	leftW  = 24                 // conversation list
	rightW = 47                 // chat / group form
	innerW = leftW + 1 + rightW // +1 for the divider
	bodyH  = 14                 // rows of the two panes
)

var (
	// highlighted row: same colors as the shared selected item, without its fixed width
	selectSt = styles.ItemSelected.UnsetWidth().UnsetAlign()

	// same box as the shared one, with tighter padding to fit the 80x24 terminal
	boxSt = styles.BoxStyle.Padding(0, 1)
)
