package leaderboardview

import (
	"fmt"
	"os"
	"strconv"
	"visweeper/internal/game"
	"visweeper/internal/leaderboard"
	"visweeper/tui/utils"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

var (
	focusedColor utils.Color = utils.WHITE
	// Styles
	style            = lip.NewStyle()
	unfocusedBorders = style.BorderStyle(lip.ThickBorder()).BorderForeground(utils.GREY)
	focusedBorders   = unfocusedBorders.BorderForeground(focusedColor)
)

type state struct {
	needName    bool
	wantRestart bool
}

type model struct {
	size        string
	game        *game.Game
	leaderboard *leaderboard.Leaderboard
	state       *state
	table       table.Model
	nameInput   textinput.Model
}

func New(size string, game *game.Game, leaderboard *leaderboard.Leaderboard) model {
	viewState := state{}
	t := loadTable(size, leaderboard)
	ni := textinput.New()
	if game != nil {
		t.Blur()
		ni.Focus()
		ni.Placeholder = "Player name"
		ni.Width = 20
		ni.CharLimit = 20
		viewState.needName = true
		updateTableStyles(&t)
	}
	return model{
		size:        size,
		game:        game,
		state:       &viewState,
		table:       t,
		nameInput:   ni,
		leaderboard: leaderboard,
	}
}

// Lifecycle ------------------------------------------------------------------
func (m model) Init() tea.Cmd {
	if m.state.needName {
		return textinput.Blink
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	tableFocused := m.table.Focused()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			if m.state.needName {
				if tableFocused {
					m.table.Blur()
				} else {
					m.table.Focus()
				}
				updateTableStyles(&m.table)
				return m, textinput.Blink
			}
		}
	}
	if tableFocused {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "ctrl+d":
				m.clearLeaderboard()
				return m, nil
			case " ":
				if !m.state.needName {
					m.cycleSize()
				}
			case "r":
				m.state.wantRestart = true
				return m, tea.Quit
			}
		}
		m.table, cmd = m.table.Update(msg)
	} else if m.state.needName {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.nameInput.Value() == "" {
					m.nameInput.PlaceholderStyle = style.Foreground(utils.RED)
					m.nameInput.Placeholder = "Please enter a name"
				} else {
					m.submitScore()
				}
			}

			m.nameInput, cmd = m.nameInput.Update(msg)
		}
	}

	return m, cmd
}

func (m model) View() string {
	tableOuterStyle := focusedBorders
	tableFocused := m.table.Focused()
	if !tableFocused {
		tableOuterStyle = unfocusedBorders
	}
	sz := style.Foreground(focusedColor).Render(m.size)
	tableView := tableOuterStyle.Render("High scores: "+sz+"\n", m.table.View())

	if m.state.needName {
		message := fmt.Sprintf("Your time: %s\nType your name and press enter to submit!\n", utils.FormatDuration(m.game.ElapsedTime))
		nameBoxStyle := focusedBorders
		if tableFocused {
			nameBoxStyle = unfocusedBorders
		}
		return lip.JoinVertical(lip.Left, tableView, nameBoxStyle.Render(message, m.nameInput.View()))
	}
	return tableView
}

// Actions ---------------------------------------------------------------------
func (m *model) submitScore() {
	m.leaderboard.InsertEntry(m.size, m.nameInput.Value(), m.game.ElapsedTime)
	m.state.needName = false
	m.table = loadTable(m.size, m.leaderboard)
}

func (m *model) clearLeaderboard() {
	m.leaderboard.Clear(m.size)
	m.table = loadTable(m.size, m.leaderboard)
}

func (m *model) cycleSize() {
	switch m.size {
	case "S":
		m.size = "M"
	case "M":
		m.size = "L"
	case "L":
		m.size = "XL"
	case "XL":
		m.size = "S"
	}
	m.table = loadTable(m.size, m.leaderboard)
}

// Helpers --------------------------------------------------------------------
func loadTable(size string, leaderboard *leaderboard.Leaderboard) table.Model {
	rows, err := leaderboard.GetAll(size)
	if err != nil {
		return table.New(
			table.WithColumns([]table.Column{
				{Title: "Error", Width: 20},
			}),
			table.WithRows([]table.Row{{"Error getting leaderboard: " + err.Error()}}),
		)
	}

	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Rank", Width: 8},
			{Title: "Player", Width: 20},
			{Title: "Time", Width: 10},
			{Title: "Date", Width: 20},
		}),
		table.WithHeight(20),
		table.WithRows(toTableRows(rows)),
		table.WithFocused(true),
	)

	focusedColor = utils.GetSizeColor(size)
	focusedBorders = unfocusedBorders.BorderForeground(focusedColor)

	updateTableStyles(&t)

	return t
}

func toTableRows(entries []leaderboard.Entry) []table.Row {
	var rows []table.Row
	for _, entry := range entries {
		rows = append(rows, table.Row{
			style.Foreground(utils.GetNumberColor(entry.Rank)).Render(strconv.Itoa(entry.Rank)),
			entry.Name,
			utils.FormatDuration(entry.Time),
			entry.Timestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return rows
}

// Rendering ------------------------------------------------------------------
func updateTableStyles(t *table.Model) {
	s := table.DefaultStyles()
	s.Header = s.Header.BorderForeground(utils.GREY).
		BorderStyle(lip.ThickBorder()).BorderBottom(true)
	s.Selected = s.Selected.Background(utils.GREY)
	if t.Focused() {
		s.Header = s.Header.BorderForeground(focusedColor)
		s.Selected = s.Selected.Background(utils.WHITE)
	}
	t.SetStyles(s)
}

// Main entrypoint ------------------------------------------------------------
func Run(size string, game *game.Game, lb *leaderboard.Leaderboard) string {
	m := New(size, game, lb)
	p := tea.NewProgram(m, tea.WithAltScreen())
	fm, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
	m = fm.(model)
	if m.state.wantRestart {
		return m.size
	}
	return ""
}
