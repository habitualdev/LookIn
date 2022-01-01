package ui

import (
	"LookIn/utils"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/patrickmn/go-cache"
	"golang.org/x/term"
	"os"
	"strconv"
	"time"
)

type model struct {
	cache         *cache.Cache
	subjectLines  []string
	cursor        int
	listStart     int
	listEnd       int
	pageSize      int
	search        string
	width, height int
	exitCode      int
	editMode	  bool
	mailView	  bool
}

var banner =`
******************************************
*                 LookIn                 *
*                                        *
*  An IMAP client for the cli inclined   *
******************************************
`

func (m *model) Init() tea.Cmd {
	m.mailView = false
	m.editMode = false
	m.width, m.height, _ = term.GetSize(int(os.Stdin.Fd()))

	m.width = m.width - 5
	m.height = m.height - 5

	return nil
}

func (m *model) View() string {

	var outputString string
	outputString += banner
	if m.mailView{
		emailBody, found := m.cache.Get(strconv.Itoa(m.cursor - 1))
		if found{
			byteContent := emailBody.(utils.CacheEntry).Body
			for _, byteLine := range byteContent{
				line := string(byteLine)
				outputString += line
			}
		}
	} else {
		n := m.listEnd
		for n > m.listStart {
			var line string
			email, found := m.cache.Get(strconv.Itoa(n - 1))
			if found {
				header := email.(utils.CacheEntry).Header
				from, _ := header.AddressList("From")
				subject, _ := header.Subject()
				datetime, _ := header.Date()

				if n == m.cursor {
					line = cursor.Render(fmt.Sprintf("%d.)Recieved: %s -- From %s -- %s", m.listEnd-n+1, datetime.Format(time.RFC822Z), from, subject)) + "\n"
				} else {
					line = fmt.Sprintf("%d.)Recieved: %s -- From %s -- %s \n", m.listEnd-n+1, datetime.Format(time.RFC822Z), from, subject)
				}
				outputString += line
				n -= 1
			} else {
				line = ""
				outputString += line
			}
		}
	}
	return outputString
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editMode {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible after prompt.
			m.exitCode = 2
			return m, tea.Quit

		case "esc":
			_, _ = fmt.Fprintln(os.Stderr) // Keep last item visible after prompt.
			return m, tea.Quit

		case "up":
			if m.cursor == m.listEnd {
				m.cursor = m.listStart + 1
			} else if m.cursor > m.listStart{
				m.cursor = m.cursor + 1
			}

		case "down":
			if m.cursor == m.listStart + 1{
				m.cursor = m.listEnd
			} else if m.cursor <= m.listEnd{
				m.cursor = m.cursor - 1
			}

		case "n":
			if !m.mailView{
			if m.listStart > 0 {
				m.listStart = m.listStart - 25
				m.listEnd = m.listEnd - 25
				m.cursor = m.cursor - 25
			}
			}

		case "p":
			if !m.mailView{
			if m.listEnd != m.cache.ItemCount() {
				m.listStart = m.listStart + 25
				m.listEnd = m.listEnd + 25
				m.cursor = m.cursor + 25
			}
			}

		case "right":
			if !m.mailView{
				if m.listStart > 0 {
					m.listStart = m.listStart - 25
					m.listEnd = m.listEnd - 25
					m.cursor = m.cursor - 25
				}
			}

		case "left":
			if !m.mailView{
				if m.listEnd != m.cache.ItemCount() {
					m.listStart = m.listStart + 25
					m.listEnd = m.listEnd + 25
					m.cursor = m.cursor + 25
				}
			}

		case "enter":
			if m.mailView{
				m.mailView = false
			} else {
				m.mailView = true
			}
		}
	}

	return m, nil
}