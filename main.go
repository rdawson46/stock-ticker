package main

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rdawson46/stock-ticker/ui"
)


func main() {
    p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())

    if _, err := p.Run(); err != nil {
        fmt.Println("broke:", err)
        os.Exit(1)
    }
}
