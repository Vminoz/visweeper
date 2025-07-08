package game

import (
	"math/rand"
	"time"
)

type Cell struct {
	IsMine        bool
	IsRevealed    bool
	IsFlagged     bool
	NeighborMines int
}

type Game struct {
	Rows, Cols  int
	Board       [][]Cell
	Mines       int
	Flags       int
	GameOver    bool
	GameWon     bool
	IsFirstMove bool
	StartTime   time.Time
	ElapsedTime time.Duration
}

func New(rows, cols, mines int) *Game {
	g := &Game{
		Rows:        rows,
		Cols:        cols,
		Mines:       mines,
		Flags:       0,
		GameOver:    false,
		GameWon:     false,
		IsFirstMove: true,
		StartTime:   time.Now(),
		ElapsedTime: 0,
	}
	g.initializeBoard()
	return g
}

func (g *Game) StartTimer() {
	g.StartTime = time.Now()
}

func (g *Game) StopTimer() {
	g.ElapsedTime = time.Since(g.StartTime)
}

func (g *Game) initializeBoard() {
	g.Board = make([][]Cell, g.Rows)
	for r := range g.Board {
		g.Board[r] = make([]Cell, g.Cols)
	}
}

func (g *Game) placeMines(excludeRow, excludeCol int) {
	minesPlaced := 0
	for minesPlaced < g.Mines {
		r := rand.Intn(g.Rows)
		c := rand.Intn(g.Cols)
		if g.Board[r][c].IsMine {
			continue
		}

		isExcluded := false
		if r == excludeRow && c == excludeCol {
			isExcluded = true
		} else {
			g.forEachNeighbor(r, c, func(nr, nc int) {
				if nr == excludeRow && nc == excludeCol {
					isExcluded = true
				}
			})
		}

		if isExcluded {
			continue
		}

		g.Board[r][c].IsMine = true
		minesPlaced++
		g.forEachNeighborCell(r, c, func(c *Cell) { c.NeighborMines++ })
	}
}

func (g *Game) isValid(row, col int) bool {
	return row >= 0 && row < g.Rows && col >= 0 && col < g.Cols
}

func (g *Game) start(row, col int) {
	g.placeMines(row, col)
	g.IsFirstMove = false
	g.StartTimer()
}

func (g *Game) Reveal(row, col int) {
	if g.IsFirstMove {
		g.start(row, col)
	}
	if !g.isValid(row, col) || g.Board[row][col].IsFlagged {
		return
	}

	cell := &g.Board[row][col]
	if cell.IsRevealed {
		g.RevealNeighbors(row, col)
		return
	}

	cell.IsRevealed = true

	if cell.IsMine {
		g.GameOver = true
		g.StopTimer()
		return
	}

	if cell.NeighborMines == 0 {
		g.RevealNeighbors(row, col)
	}
	g.checkWinCondition()
}

func (g *Game) RevealNeighbors(row, col int) {
	nm := g.Board[row][col].NeighborMines
	if nm > 0 {
		n := 0
		g.forEachNeighborCell(row, col, func(c *Cell) {
			if c.IsFlagged {
				n++
			}
		})
		if n != nm {
			return
		}
	}

	g.forEachNeighbor(row, col, func(nr, nc int) {
		if !g.Board[nr][nc].IsRevealed {
			g.Reveal(nr, nc)
		}
	})
}

func (g *Game) Flag(row, col int) {
	if !g.isValid(row, col) {
		return
	}
	if g.Board[row][col].IsRevealed {
		g.FlagNeighbors(row, col)
		return
	}
	g.Board[row][col].IsFlagged = !g.Board[row][col].IsFlagged
	if g.Board[row][col].IsFlagged {
		g.Flags++
	} else {
		g.Flags--
	}
}

func (g *Game) FlagNeighbors(row, col int) {
	nm := g.Board[row][col].NeighborMines
	if nm == 0 {
		return
	}

	n := 0
	g.forEachNeighborCell(row, col, func(c *Cell) {
		if !(c.IsRevealed) {
			n++
		}
	})
	if n != nm {
		return
	}

	g.forEachNeighborCell(row, col, func(c *Cell) {
		c.IsFlagged = !c.IsRevealed
	})
}

func (g *Game) RevealedCells() int {
	revealedCount := 0
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			if g.Board[r][c].IsRevealed {
				revealedCount++
			}
		}
	}
	return revealedCount
}

func (g *Game) forEachNeighbor(row, col int, callback func(r, c int)) {
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr, nc := row+dr, col+dc
			if g.isValid(nr, nc) {
				callback(nr, nc)
			}
		}
	}
}

func (g *Game) forEachNeighborCell(row, col int, callback func(c *Cell)) {
	g.forEachNeighbor(row, col, func(r, c int) { callback(&g.Board[r][c]) })
}

func (g *Game) checkWinCondition() {
	if g.RevealedCells() == g.Rows*g.Cols-g.Mines {
		g.GameWon = true
		g.GameOver = true
		g.StopTimer()
	}
}
