package main

import (
	"flag"
	"gosweeper/game"
	"gosweeper/tui"
)

func main() {
	sizeS := flag.Bool("S", false, "Small  (10x10) (default)")
	sizeM := flag.Bool("M", false, "Medium (16x16)")
	sizeL := flag.Bool("L", false, "Large  (16x30)")
	sizeXL := flag.Bool("X", false, "XL    (36x36)")
	cheat := flag.Bool("cheat", false, "allows showing mines (no score)")
	arrowKeys := flag.Bool("arrow-keys", false, "allows using arrow keys (no score)")
	minePercent := flag.Int("mine-percent", 16, "override percentage of mines (no score)")

	flag.Parse()

	rows, cols := 10, 10

	if *sizeS {
		rows, cols = 10, 10
	} else if *sizeM {
		rows, cols = 16, 16
	} else if *sizeL {
		rows, cols = 16, 30
	} else if *sizeXL {
		rows, cols = 36, 36
	}

	mines := int(float32(rows*cols) * float32(*minePercent) / 100.0)

	g := game.New(rows, cols, mines)
	opts := tui.GameOptions{
		Cheat:        *cheat,
		UseArrowKeys: *arrowKeys,
	}
	tui.Run(g, opts)
}
