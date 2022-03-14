package main

import (
	. "github.com/charmbracelet/lipgloss"
)

var (
	E      = NewStyle().Foreground(Color("#ae2012")).Render
	W      = NewStyle().Foreground(Color("#ee9b00")).Render
	OK     = NewStyle().Foreground(AdaptiveColor{Dark: "#52b788", Light: "#2d6a4f"}).Render
	Num    = NewStyle().Foreground(Color("#b48ead")).Render
	border = NewStyle().Border(NormalBorder(), true).Padding(0, 1).Render
	center = NewStyle().Align(Center).Render
	H      = NewStyle().Foreground(Color("#84a98c")).Render
)

func ternary(c bool, s1 func(string) string, s2 func(string) string) func(string) string {
	if c {
		return s1
	}
	return s2
}
