package tui

import (
	"fmt"
	"gh-user-activity/client"
	"gh-user-activity/formatter"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateInput state = iota
	stateLoading
	stateResults
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Foreground(lipgloss.Color("205"))

	eventTypeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("220"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	detailStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			PaddingLeft(4)
)

func eventMarker(eventType string) string {
	switch eventType {
	case "PushEvent":
		return "[Push]"
	case "IssuesEvent":
		return "[Issue]"
	case "WatchEvent":
		return "[Star]"
	case "ForkEvent":
		return "[Fork]"
	case "CreateEvent":
		return "[Create]"
	case "DeleteEvent":
		return "[Delete]"
	case "PullRequestEvent":
		return "[PR]"
	case "IssueCommentEvent":
		return "[Comment]"
	case "ReleaseEvent":
		return "[Release]"
	case "MemberEvent":
		return "[Member]"
	case "PublicEvent":
		return "[Public]"
	case "GollumEvent":
		return "[Wiki]"
	case "PullRequestReviewEvent":
		return "[Review]"
	case "PullRequestReviewCommentEvent":
		return "[Comment]"
	default:
		return "[Event]"
	}
}

type eventDetail struct {
	lines []string
}

func buildDetail(e client.Event) eventDetail {
	var d eventDetail
	switch e.Type {
	case "PushEvent":
		for _, c := range e.Payload.Commits {
			msg := strings.Split(c.Message, "\n")[0]
			d.lines = append(d.lines, fmt.Sprintf("commit %s: %s", c.SHA[:7], msg))
		}
	case "IssuesEvent":
		if e.Payload.Issue != nil {
			d.lines = append(d.lines, fmt.Sprintf("#%d: %s", e.Payload.Issue.Number, e.Payload.Issue.Title))
		}
	case "PullRequestEvent":
		if e.Payload.PullRequest != nil {
			d.lines = append(d.lines, fmt.Sprintf("#%d: %s", e.Payload.PullRequest.Number, e.Payload.PullRequest.Title))
		}
	case "PullRequestReviewEvent", "PullRequestReviewCommentEvent":
		if e.Payload.PullRequest != nil {
			d.lines = append(d.lines, fmt.Sprintf("PR #%d: %s", e.Payload.PullRequest.Number, e.Payload.PullRequest.Title))
		}
	case "CreateEvent":
		if e.Payload.Ref != "" {
			d.lines = append(d.lines, fmt.Sprintf("ref: %s (%s)", e.Payload.Ref, e.Payload.RefType))
		}
	case "DeleteEvent":
		d.lines = append(d.lines, fmt.Sprintf("ref: %s (%s)", e.Payload.Ref, e.Payload.RefType))
	}
	return d
}

type fetchEventsMsg struct {
	events []client.Event
	err    error
}

type model struct {
	state       state
	username    string
	textInput   textinput.Model
	spinner     spinner.Model
	viewport    viewport.Model
	events      []client.Event
	err         error
	ready       bool
	expanded    map[int]bool
	windowWidth int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "github username"
	ti.Focus()
	ti.CharLimit = 64
	ti.Width = 40

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	s.Spinner = spinner.Dot

	return model{
		state:     stateInput,
		textInput: ti,
		spinner:   s,
		expanded:  make(map[int]bool),
	}
}

func RunTUI() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func fetchEventsCmd(username string) tea.Cmd {
	return func() tea.Msg {
		events, err := client.GetUserActivity(username)
		return fetchEventsMsg{events: events, err: err}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		headerHeight := 3
		footerHeight := 2
		m.viewport.Width = msg.Width - 2
		m.viewport.Height = msg.Height - headerHeight - footerHeight
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case stateInput:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				username := strings.TrimSpace(m.textInput.Value())
				if username == "" {
					return m, nil
				}
				m.username = username
				m.state = stateLoading
				return m, fetchEventsCmd(username)
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd

		case stateLoading:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			}

		case stateResults:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "esc":
				m.state = stateInput
				m.textInput.Reset()
				m.textInput.Focus()
				m.events = nil
				m.err = nil
				m.expanded = make(map[int]bool)
				return m, nil
			case "ctrl+r":
				m.state = stateInput
				m.textInput.Reset()
				m.textInput.Focus()
				m.events = nil
				m.err = nil
				m.expanded = make(map[int]bool)
				return m, nil
			case "enter":
				top := m.viewport.YOffset
				if top >= 0 && top < len(m.events) {
					if m.expanded[top] {
						delete(m.expanded, top)
					} else {
						m.expanded[top] = true
					}
					m.viewport.SetContent(m.renderEvents())
				}
				return m, nil
			}
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}

	case fetchEventsMsg:
		m.events = msg.events
		m.err = msg.err
		m.state = stateResults
		m.viewport.SetContent(m.renderEvents())
		return m, nil

	case spinner.TickMsg:
		if m.state == stateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case stateInput:
		return m.viewInput()
	case stateLoading:
		return m.viewLoading()
	case stateResults:
		return m.viewResults()
	default:
		return ""
	}
}

func (m model) viewInput() string {
	var b strings.Builder
	b.WriteString("\n\n\n")
	b.WriteString(titleStyle.Render(" GitHub User Activity "))
	b.WriteString("\n\n")
	b.WriteString("Enter a GitHub username:\n\n")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n  [enter] submit  •  [q] quit")
	return lipgloss.NewStyle().Align(lipgloss.Center).Width(m.windowWidth).Render(b.String())
}

func (m model) viewLoading() string {
	return fmt.Sprintf("\n\n%s Fetching activity for %s...", m.spinner.View(), m.username)
}

func (m model) viewResults() string {
	var b strings.Builder
	title := fmt.Sprintf(" GitHub Activity: %s ", m.username)
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.Repeat("─", min(m.windowWidth-1, 60))))
	b.WriteString("\n")
	b.WriteString(m.viewport.View())
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(" [↑/↓] scroll  [enter] toggle detail  [esc] search  [q] quit "))
	return b.String()
}

func (m model) renderEvents() string {
	if m.err != nil {
		return fmt.Sprintf("\n%s\n\n%s",
			dimStyle.Render("Error fetching activity:"),
			lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(m.err.Error()),
		)
	}
	if len(m.events) == 0 {
		return dimStyle.Render("No recent events found.")
	}

	var b strings.Builder
	for i, e := range m.events {
		marker := eventMarker(e.Type)
		desc := formatter.FormatEvent(e)
		ts := e.CreatedAt.Format("2006-01-02 15:04")

		line := fmt.Sprintf("%s %s",
			eventTypeStyle.Render(marker),
			desc,
		)
		b.WriteString(line)
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("  " + ts))
		b.WriteString("\n")

		if m.expanded[i] {
			d := buildDetail(e)
			for _, l := range d.lines {
				b.WriteString(detailStyle.Render(l))
				b.WriteString("\n")
			}
		}

		if i < len(m.events)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}
