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
