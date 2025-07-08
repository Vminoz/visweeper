package main

import (
	"flag"

	"visweeper/internal/game"
	"visweeper/tui/gameview"
)

func main() {
	flag.Bool("S", false, "Small  (10x10) (default)")
	sizeM := flag.Bool("M", false, "Medium (16x16)")
	sizeL := flag.Bool("L", false, "Large  (16x30)")
	sizeXL := flag.Bool("X", false, "XL    (36x36)")
	cheat := flag.Bool("cheat", false, "allows showing mines (no score)")
	arrowKeys := flag.Bool("arrow-keys", false, "allows using arrow keys (no score)")
	minePercent := flag.Int("mine-percent", 16, "override percentage of mines (no score)")

	flag.Parse()

	rows, cols := 10, 10

	if *sizeM {
		rows, cols = 16, 16
	} else if *sizeL {
		rows, cols = 16, 30
	} else if *sizeXL {
		rows, cols = 36, 36
	}

	mines := int(float32(rows*cols) * float32(*minePercent) / 100.0)

	g := game.New(rows, cols, mines)
	opts := gameview.GameOptions{
		Cheat:        *cheat,
		UseArrowKeys: *arrowKeys,
	}

	gameview.Run(g, opts)
}
