package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/amanallah-jendoubi/matrix-chat/tui/app"
	"os"
)

func main() {
	f, err := tea.LogToFile("debug.log", "matrix_chat")
	if err != nil {
		fmt.Println("cannot open log file:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(app.NewApp())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}
