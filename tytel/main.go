package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type passage struct {
	Name string
	Text string
}

var passages []passage

// ////////////////////////////////////////////////////////////////////////////
// styles:
// ////////////////////////////////////////////////////////////////////////////
var (
	titleStyle = lipgloss.NewStyle().Bold(true).Underline(true)
	// seems like booleans are lowercase; recall in python they are uppercase
	infoStyle    = lipgloss.NewStyle().Faint(true)
	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")) // greenish
	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	cursorStyle  = lipgloss.NewStyle().Underline(true)
)

// ////////////////////////////////////////////////////////////////////////////
// model:
// ////////////////////////////////////////////////////////////////////////////
type model struct {
	target []rune // the passage to type
	typed  []rune // what the user has typed

	passageName string

	startedAt  time.Time
	finishedAt time.Time
	done       bool
}

func newRandomModel() model {
	if len(passages) == 0 {
		passages = defaultPassages()
	}
	rand.Seed(time.Now().UnixNano())
	p := passages[rand.Intn(len(passages))]

	return model{
		target:      []rune(strings.TrimSpace(p.Text)),
		passageName: p.Name,
	}
}

func defaultPassages() []passage {
	return []passage{
		{
			Name: "Shakespeare - Sonnet 18 (opening)",
			Text: "Shall I compare thee to a summer's day? Thou art more lovely and more...",
		},
		{
			Name: "Shakespeare - The Tempest",
			Text: "We are such stuff as dreams are made on, and our little life is rounded...",
		},
	}
}

// ////////////////////////////////////////////////////////////////////////////
// bubbletea interface:
// ////////////////////////////////////////////////////////////////////////////
func (m model) Init() tea.Cmd {
	// no initial async commands needed for now
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyBackspace:
			if !m.done && len(m.typed) > 0 {
				m.typed = m.typed[:len(m.typed)-1]
			}
			return m, nil
		}
		// if finished, allow restart with 'r'
		if m.done {
			if msg.String() == "r" {
				return newRandomModel(), nil
			}
		}
		// convert keypress into a rune we care about
		var r rune
		switch msg.Type {
		case tea.KeyEnter:
			r = '\n'
		case tea.KeySpace:
			r = ' '
		default:
			if len(msg.Runes) == 0 {
				//ignore non-character keys
				return m, nil
			}
			r = msg.Runes[0]
		}
		if m.startedAt.IsZero() {
			m.startedAt = time.Now()
		}
		m.typed = append(m.typed, r)
		if len(m.typed) >= len(m.target) {
			m.done = true
			m.finishedAt = time.Now()
		}
	}
	return m, nil
}
