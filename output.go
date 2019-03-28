package main

import (
	"github.com/gdamore/tcell"
	runewidth "github.com/mattn/go-runewidth"
)

var (
	row   = 0
	style = tcell.StyleDefault
)

func clearRow(screen tcell.Screen, row int) {
	width, _ := screen.Size()
	for i := 0; i <= width; i++ {
		screen.SetContent(i, row, 32, []rune{}, style)
	}
}

func moverow(screen tcell.Screen, row int) {
	var future_mainc rune
	var future_st tcell.Style
	width := maxWordsWidth(screen)
	for i := 1; i <= width; i += 2 {
		mainc, _, st, _ := screen.GetContent(i, row)
		if i > 0 {
			if i == width && mainc != 32 {
				// At this point the game is lost
			}
			screen.SetContent(i, row, future_mainc, []rune{}, future_st)
		} else {
			screen.SetContent(i, row, rune(32), []rune{}, style)
		}
		future_mainc, _, future_st, _ = screen.GetContent(i+1, row)
		screen.SetContent(i+1, row, mainc, []rune{}, st)
	}
}

func writexy(screen tcell.Screen, x, y int, content string) {
	puts(screen, style, x, y, content)
	screen.Show()
}

func writex(screen tcell.Screen, x int, content string) {
	puts(screen, style, x, row, content)
	screen.Show()
}

func writey(screen tcell.Screen, y int, content string) {
	puts(screen, style, 1, y, content)
	screen.Show()
}

func writeln(screen tcell.Screen, content string) {
	puts(screen, style, 1, row, content)
	screen.Show()
	row++
}

func write(screen tcell.Screen, content string) {
	puts(screen, style, 1, row, content)
	screen.Show()
}

func getRowContent(screen tcell.Screen, row int) []string {
	width, _ := screen.Size()
	now_character := false
	list := []string{}
	current_content := ""
	for i := 0; i <= width; i++ {
		mainc, _, _, _ := screen.GetContent(i, row)
		if mainc != 32 {
			if !now_character {
				now_character = true
			}
			current_content = current_content + string(mainc)
		} else {
			if now_character {
				now_character = false
				list = append(list, current_content)
				current_content = ""
			}
		}
	}
	row++
	return list
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	zwj := false
	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}
