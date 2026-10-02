package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/amanallah-jendoubi/Textio/tui/app"
	"os"
)

func main() {
	p := tea.NewProgram(app.NewApp())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
