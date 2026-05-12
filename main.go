package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
