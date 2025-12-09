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
	// can we refactor to use polymorphism?
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

func (m model) View() string {
	if len(m.target) == 0 {
		return "No text loaded.\n"
	}

	var b strings.Builder

	// title
	b.WriteString(titleStyle.Render("Typing Practice") + "\n")
	if m.passageName != "" {
		b.WriteString(infoStyle.Render("Passage: "+m.passageName) + "\n")
	}
	b.WriteString("\n")

	// render text with per-character colouring:
	for i, r := range m.target {
		display := string(r)
		if r == ' ' {
			display = "." // show spaces so you can see them
		}

		if i < len(m.typed) {
			if m.typed[i] == r {
				b.WriteString(correctStyle.Render(display))
			} else {
				b.WriteString(wrongStyle.Render(display))
			}
		} else if i == len(m.typed) && !m.done {
			b.WriteString(cursorStyle.Render(display))
		} else {
			b.WriteString(display)
		}
	}
	b.WriteString("\n\n")

	// stats:
	wpm, acc := m.stats()
	elapsed := m.elapsed().Seconds()

	b.WriteString(fmt.Sprintf(
		"Time: %.1fs WPM: %.1f Accuracy: %.1f%%\n",
		elapsed, wpm, acc,
	))
	b.WriteString(fmt.Sprintf(
		"Chars: %d / %d\n",
		len(m.typed), len(m.target),
	))

	// instructions
	b.WriteString("\n")
	if m.done {
		b.WriteString(infoStyle.Render("Finished! Press 'r' for a new passage, or Ctrl+C to"))
	} else if m.startedAt.IsZero() {
		b.WriteString(infoStyle.Render("Start typing to begin. Backspace to correct, Esc/Ctrl to ..."))
	} else {
		b.WriteString(infoStyle.Render("keep going! Backspace to correct, Esc/Ctrl+C to quit"))
	}

	return b.String()

}

// stats helpers
func (m model) elapsed() time.Duration {
	if m.startedAt.IsZero() {
		return 0
	}
	if m.done {
		return m.finishedAt.Sub(m.startedAt)
	}
	return time.Since(m.startedAt)
}

func (m model) stats() (wpm float64, accuracy float64) {
	if len(m.typed) == 0 {
		return 0, 0
	}

	correct := 0
	for i, r := range m.typed {
		if i < len(m.target) && r == m.target[i] {
			correct++
		}
	}
	total := len(m.typed)

	elapsedMinutes := m.elapsed().Minutes()
	if elapsedMinutes <= 0 {
		// avoid division by zero when you finish super quickly
		elapsedMinutes = 1.0 / 60.0
	}

	// standard typing metric : 5 characters = 1 'word'
	// change: TODO:
	wpm = (float64(correct) / 5.0) / elapsedMinutes
	accuracy = float64(correct) / float64(total) * 100.0

	return wpm, accuracy
}

// main cli
func main() {
	filePath := flag.String("file", "", "path to a text file to use as the passage")
	rawText := flag.String("text", "", "inline text to use as the passage (overrides -file)")
	flag.Parse() // provided from top import

	//decide where packages come from:
	if *rawText != "" {
		passages = []passage{
			{Name: "Custom text", Text: *rawText},
		}
	} else if *filePath != "" {
		// load from file
		data, err := os.ReadFile(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read file %q: %v\n", *filePath, err)
			passages = defaultPassages()
		} else {
			// entire file as one passage; we can later split by lines.
			passages = []passage{
				{Name: *filePath, Text: string(data)},
			}
		}
	} else {
		// fall back to built-in
		passages = defaultPassages()
	}

	p := tea.NewProgram(newRandomModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
