# Visweeper

Minesweepr TUI with forced vi-like cursor movement (`hjkl`).
![Screenshot-L](https://github.com/user-attachments/assets/417a3d77-44cc-4e50-8f6d-e5427052d2ff)

## Install and play
Assuming you an appropriate version of go installed:
```sh
go install github.com/vminoz/visweeper@latest
visweeper
```

## CLI
```
visweeper <flags>
  -L    Large  (16x30)
  -M    Medium (16x16)
  -S    Small  (10x10) (default)
  -X    XL     (36x36)
  -arrow-keys
        allows using arrow keys (no score)
  -cheat
        allows showing mines (no score)
  -leaderboard
        Look at leaderboard
  -mine-percent int
        override percentage of mines (no score) (default 16)
```

## Key Bindings
`ctrl+c` and `Esc` always exit the program, so does `q` when not typing a text.

#### Game View
Hit `?` to show controls.

#### Leaderboard View
| Key      | Action                                           | Condition        |
| ---      | ---                                              | ---              |
| `Tab`    | Cycle focus between table and name input         | name box visible |
| `space`  | Cycle through board sizes                        | name box hidden  |
| `r`      | Restart/Run a game with the current size         | table focused    |
| `ctrl+d` | Clear leaderboard for current size (permanently) | table focused    |

## Details
- Uses [Bubbletea](https://github.com/charmbracelet/bubbletea) and related libraries for TUI stuff!
- The Leaderboards are stored in a sqlite database at `$HOME/.visweeper/leaderboard.db`
