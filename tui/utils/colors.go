package utils

import (
	lip "github.com/charmbracelet/lipgloss"
)

type Color lip.TerminalColor

const (
	// Generic
	BLACK = lip.Color("0")
	// ANSI intense
	GREY   = lip.Color("8")
	RED    = lip.Color("9")
	GREEN  = lip.Color("10")
	YELLOW = lip.Color("11")
	BLUE   = lip.Color("12")
	PURPLE = lip.Color("13")
	CYAN   = lip.Color("14")
	WHITE  = lip.Color("15")
)

func GetNumberColor(n int) Color {
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

func GetSizeColor(sizeName string) Color {
	switch sizeName {
	case "S":
		return GREEN
	case "M":
		return YELLOW
	case "L":
		return RED
	case "XL":
		return PURPLE
	default:
		return WHITE
	}
}
