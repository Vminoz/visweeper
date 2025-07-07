package tui

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"visweeper/game"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

const (
	// Frame colors
	defaultFrameColor = GREY
	winFrameColor     = GREEN
	lossFrameColor    = RED
)

var (
	// Styles
	style           = lip.NewStyle()
	cellStyle       = style.MarginRight(1)
	hiddenCellStyle = cellStyle.Foreground(GREY)
	cursorStyle     = cellStyle.Background(WHITE).Foreground(BLACK)
)

type GameOptions struct {
	Cheat        bool
	UseArrowKeys bool
}

type gameState struct {
	CursorX   int
	CursorY   int
	showMines bool
	showHelp  bool
	ticker    *time.Ticker
	message   string
	numBuff   string
}

type model struct {
	Game  *game.Game
	state gameState
	opts  GameOptions
	done  bool
}

func New(g *game.Game) *model {
	state := gameState{
		CursorX: g.Cols/2 - 1,
		CursorY: g.Rows/2 - 1,
	}

	return &model{
		Game:  g,
		state: state,
	}
}

// Lifecycle ------------------------------------------------------------------
func (m *model) Init() tea.Cmd {
	m.state.ticker = time.NewTicker(time.Second)
	return m.tick()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s := &m.state
	g := m.Game
	movementKeys := map[string]struct{ dx, dy int }{
		"h": {-1, 0},
		"j": {0, 1},
		"k": {0, -1},
		"l": {1, 0},
	}
	arrowKeys := map[string]struct{ dx, dy int }{
		"left":  {-1, 0},
		"down":  {0, 1},
		"up":    {0, -1},
		"right": {1, 0},
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		k := msg.String()
		if move, ok := movementKeys[k]; ok {
			m.moveCursor(move.dx, move.dy)
		} else if arrow, ok := arrowKeys[k]; ok {
			if !m.opts.UseArrowKeys {
				m.state.message += "h← j↓ k↑ l→"
				return m, nil
			}
			m.moveCursor(arrow.dx, arrow.dy)
		} else {
			switch k {
			case "f":
				m.flag()
			case "r":
				m.Game = game.New(g.Rows, g.Cols, g.Mines)
				s.showMines = false
			case "?":
				s.showHelp = !s.showHelp
			case "q", "ctrl+c", "esc":
				m.done = true
				return m, tea.Quit
			case " ", "d":
				m.detonate()
			case "s":
				if m.opts.Cheat {
					s.showMines = !s.showMines
				}
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				s.numBuff += k
			case "0":
				if s.numBuff != "" {
					s.numBuff += "0"
				}
			}
		}
	case tickMsg:
		if !g.GameOver && !g.IsFirstMove {
			g.ElapsedTime = time.Since(g.StartTime)
		}
		return m, m.tick()
	}
	return m, nil
}

func (m *model) View() string {
	if m.done {
		return ""
	}

	frameColor := getFrameColor(m.Game)
	frameStyle := style.Border(lip.NormalBorder()).BorderForeground(frameColor)
	bottomBordered := frameStyle.Border(lip.NormalBorder(), false, false, true, false)

	banner := m.drawBanner()
	board := m.drawBoard()

	mainContent := lip.JoinVertical(
		lip.Left,
		bottomBordered.Render(banner),
		board,
	)

	view := frameStyle.MarginBottom(1).Render(mainContent)

	if m.state.message != "" {
		view = lip.JoinVertical(lip.Left, view, m.state.message)
		m.state.message = ""
	}

	if m.state.showHelp || m.Game.GameOver {
		help := m.drawHelp()
		view = lip.JoinVertical(lip.Left, view, help)
	}

	return view
}

type tickMsg time.Time

func (m *model) tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Actions --------------------------------------------------------------------
func (m *model) flag() {
	if !m.Game.GameOver {
		m.Game.Flag(m.state.CursorY, m.state.CursorX)
	}
}

func (m *model) detonate() {
	if !m.Game.GameOver {
		m.Game.Reveal(m.state.CursorY, m.state.CursorX)
	}
}

func (m *model) moveCursor(dx, dy int) {
	s := &m.state
	if s.numBuff != "" {
		repeat, err := strconv.Atoi(s.numBuff)
		if err == nil {
			dx *= repeat
			dy *= repeat
			s.numBuff = ""
		} else {
			s.message += "failed to use numBuff: " + err.Error()
		}
	}
	s.CursorX = clamp(s.CursorX+dx, 0, m.Game.Cols)
	s.CursorY = clamp(s.CursorY+dy, 0, m.Game.Rows)
}

// Drawing --------------------------------------------------------------------
func (m *model) drawBoard() string {
	var rows []string
	for r := 0; r < m.Game.Rows; r++ {
		row := " "
		for c := 0; c < m.Game.Cols; c++ {
			row += m.renderCell(r, c)
		}
		rows = append(rows, row)
	}
	return lip.JoinVertical(lip.Left, rows...)
}

func (m *model) renderCell(r, c int) string {
	cell := m.Game.Board[r][c]
	var content string
	var style lip.Style

	if (m.state.showMines || m.Game.GameOver) && cell.IsMine {
		content = "⁘"
		style = cellStyle.Foreground(RED)
		if m.Game.GameWon {
			content = "◎"
			style = style.Foreground(WHITE)
		}
	} else if cell.IsRevealed {
		if cell.NeighborMines > 0 {
			content = fmt.Sprintf("%d", cell.NeighborMines)
			style = cellStyle.Foreground(getNumberColor(cell.NeighborMines))
		} else {
			content = " "
			style = cellStyle
		}
	} else if cell.IsFlagged {
		content = "F"
		if m.Game.GameOver && cell.IsMine {
			style = cellStyle.Foreground(GREEN)
		} else {
			style = cellStyle.Foreground(YELLOW)
		}
	} else {
		content = "·"
		style = hiddenCellStyle
	}

	if r == m.state.CursorY && c == m.state.CursorX {
		style = cursorStyle
	}

	return style.Render(content)
}

func (m *model) drawBanner() string {
	w := m.Game.Cols*2 + 1
	if m.Game.GameOver {
		return m.drawGameOverBanner(w)
	} else {
		return m.drawNormalBanner(w)
	}
}

func (m *model) drawNormalBanner(w int) string {
	// L1
	fs := fmt.Sprintf("%d/%d", m.Game.Flags, m.Game.Mines)
	ts := formatDuration(m.Game.ElapsedTime)
	lsty := style.Width(w / 2).Align(lip.Left)
	rsty := style.Width(w/2 + w%2).Align(lip.Right)
	flagsPart := lsty.Align(lip.Left).Render(fs)
	timerPart := rsty.Align(lip.Right).Render(ts)
	line1 := lip.JoinHorizontal(lip.Top, flagsPart, timerPart)

	// L2
	line2 := ""
	if m.state.message != "" && len(m.state.message) < w {
		// Flash message
		line2 = style.Width(w).AlignHorizontal(lip.Center).Foreground(RED).Background(WHITE).Render(m.state.message)
		m.state.message = ""
	} else {
		// Progress bar
		revealed := m.Game.RevealedCells()
		totalClear := m.Game.Rows*m.Game.Cols - m.Game.Mines
		progress := 0.0
		if totalClear > 0 {
			progress = float64(revealed) / float64(totalClear)
		}
		progressChars := int(progress * float64(w))

		barForeground := style.Width(w).Render(m.state.numBuff)
		p := barForeground[:progressChars]
		r := barForeground[progressChars:]
		fg := style.Foreground(BLACK)
		p = fg.Background(YELLOW).Render(p)
		r = fg.Background(GREY).Render(r)

		line2 = p + r
	}
	return lip.JoinVertical(lip.Left, line1, line2)
}

func (m *model) drawGameOverBanner(w int) string {
	bannerStyle := style.Width(w)
	msg := ""
	fg := WHITE
	if m.Game.GameWon {
		msg = "Win!"
		fg = winFrameColor
	} else {
		msg = "Game Over"
		fg = lossFrameColor
	}
	msgStyle := style.Foreground(fg)
	timerStr := formatDuration(m.Game.ElapsedTime)

	line1 := bannerStyle.Align(lip.Center).Render(msgStyle.Render(msg))
	line2 := bannerStyle.Align(lip.Right).Render(timerStr)
	return lip.JoinVertical(lip.Left, line1, line2)
}

func (m *model) drawHelp() string {
	var s []string
	if !m.Game.GameOver {
		s = append(s, style.Foreground(YELLOW).Render("F")+"lag")
	}
	s = append(s, style.Foreground(YELLOW).Render("R")+"eplay")
	s = append(s, style.Foreground(YELLOW).Render("Q")+"uit")
	return lip.JoinVertical(lip.Left, s...)
}

func getNumberColor(n int) color {
	switch n {
	case 1:
		return BLUE
	case 2:
		return GREEN
	case 3:
		return RED
	case 4:
		return PURPLE
	case 5:
		return CYAN
	case 6:
		return YELLOW
	case 7:
		return WHITE
	case 8:
		return RED
	default:
		return GREY
	}
}

func getFrameColor(g *game.Game) color {
	if g.GameOver {
		if g.GameWon {
			return winFrameColor
		} else {
			return lossFrameColor
		}
	}
	return defaultFrameColor
}

// Main entrypoint ------------------------------------------------------------
func Run(g *game.Game, options GameOptions) {
	m := New(g)
	m.opts = options
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
