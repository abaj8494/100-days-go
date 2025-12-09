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
