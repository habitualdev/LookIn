package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jwalton/go-supportscolor"
	"github.com/muesli/termenv"
	"github.com/patrickmn/go-cache"
	"time"

	//"github.com/patrickmn/go-cache"
	"os"
)

var (
	modified  = lipgloss.NewStyle().Foreground(lipgloss.Color("#588FE6"))
	added     = lipgloss.NewStyle().Foreground(lipgloss.Color("#6ECC8E"))
	untracked = lipgloss.NewStyle().Foreground(lipgloss.Color("#D95C50"))
	cursor    = lipgloss.NewStyle().Background(lipgloss.Color("#825DF2")).Foreground(lipgloss.Color("#FFFFFF"))
	bar       = lipgloss.NewStyle().Background(lipgloss.Color("#5C5C5C")).Foreground(lipgloss.Color("#FFFFFF"))
)

func StartUi(ca *cache.Cache) {
	term := supportscolor.Stderr()
	if term.Has16m {
		lipgloss.SetColorProfile(termenv.TrueColor)
	} else if term.Has256 {
		lipgloss.SetColorProfile(termenv.ANSI256)
	} else {
		lipgloss.SetColorProfile(termenv.ANSI)
	}

	fmt.Println("Retrieving Emails...")
	time.Sleep(2 * time.Second)
	fmt.Println("Initializing local email database...")
	time.Sleep(1 * time.Second)
	termenv.ClearScreen()

	m := &model{cache: ca, listStart: ca.ItemCount() - 25, listEnd: ca.ItemCount(), cursor: ca.ItemCount()}

	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	p.Start()
	os.Exit(m.exitCode)

}
