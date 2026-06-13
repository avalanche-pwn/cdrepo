package view

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/avalanche-pwn/cdrepo/internal/core"
	"github.com/avalanche-pwn/cdrepo/internal/searchif"
)

func Run() {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	m := initialModel()
	m.search_meta, err = core.InitSearch()

	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))
	tea_m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	m = tea_m.(model)
	core.FinSearch(m.search_meta)
	tty.Close()
	fmt.Fprint(os.Stdout, m.repos[m.cursor].Value)
}

type model struct {
	textInput   textinput.Model
	err         error
	quitting    bool
	repos       []*searchif.ViewSearchResult
	cursor      int
	search_meta core.SearchMeta
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.SetVirtualCursor(false)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	return model{textInput: ti}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, cmd
		case "down":
			if m.cursor < len(m.repos)-1 {
				m.cursor++
			}
			return m, cmd
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	m.repos = core.Search(m.search_meta, m.textInput.Value())

	return m, cmd
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y += lipgloss.Height(m.headerView())
	}

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), m.textInput.View(), m.footerView())
	if m.quitting {
		str = ""
	}

	v := tea.NewView(str)
	v.Cursor = c
	return v
}

func (m model) headerView() string {
	res := "Repos:\n"
	search_res := make([]string, len(m.repos))
	for i, repo := range m.repos {
		if i == m.cursor {
			search_res[i] = "> " + repo.Value
			continue
		}
		search_res[i] = repo.Value
	}
	res += strings.Join(search_res, "\n")
	return res
}
func (m model) footerView() string { return "\n(esc to quit)" }
