package main

import (
	"flag"
	"log"

	"github.com/Vminoz/visweeper/internal/game"
	"github.com/Vminoz/visweeper/internal/leaderboard"
	"github.com/Vminoz/visweeper/tui/gameview"
	"github.com/Vminoz/visweeper/tui/leaderboardview"
)

const (
	// Frame colors
	defaultMinePercent = 16
)

var sizes = map[string](struct{ rows, cols int }){
	"S":  {10, 10},
	"M":  {16, 16},
	"L":  {16, 30},
	"XL": {36, 36},
}

func main() {
	flag.Bool("S", false, "Small  (10x10) (default)")
	sizeM := flag.Bool("M", false, "Medium (16x16)")
	sizeL := flag.Bool("L", false, "Large  (16x30)")
	sizeXL := flag.Bool("X", false, "XL     (36x36)")
	cheat := flag.Bool("cheat", false, "allows showing mines (no score)")
	arrowKeys := flag.Bool("arrow-keys", false, "allows using arrow keys (no score)")
	minePercent := flag.Int("mine-percent", defaultMinePercent, "override percentage of mines (no score)")
	showLeaderboard := flag.Bool("leaderboard", false, "Look at leaderboard")

	flag.Parse()

	lb, err := leaderboard.New()
	if err != nil {
		log.Fatal(err)
	}
	defer lb.Close()

	var size string
	if *sizeM {
		size = "M"
	} else if *sizeL {
		size = "L"
	} else if *sizeXL {
		size = "XL"
	} else {
		size = "S"
	}

	if *minePercent != defaultMinePercent {
		if *minePercent < 0 || *minePercent > 100 {
			log.Fatal("Mine percentage must be between 0 and 100")
		}
	}

	gameLoop(
		size,
		*minePercent,
		*showLeaderboard,
		gameview.GameOptions{
			Cheat:        *cheat,
			UseArrowKeys: *arrowKeys,
		},
		lb,
	)
}

func gameLoop(
	size string,
	minePercent int,
	startAtLeaderboard bool,
	gameOpts gameview.GameOptions,
	lb *leaderboard.Leaderboard,
) {
	for size != "" {
		if startAtLeaderboard {
			size = leaderboardview.Run(size, nil, lb)
			startAtLeaderboard = false
			continue
		}

		rows, cols := sizes[size].rows, sizes[size].cols
		mines := int(float32(rows*cols) * float32(minePercent) / 100.0)

		g := game.New(rows, cols, mines)

		endGame := gameview.Run(g, gameOpts)

		keepScore := endGame.GameWon &&
			minePercent == defaultMinePercent &&
			!gameOpts.Cheat &&
			!gameOpts.UseArrowKeys

		if !keepScore {
			return
		}
		size = leaderboardview.Run(size, endGame, lb)
	}
}
