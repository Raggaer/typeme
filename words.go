package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/gdamore/tcell"
)

var (
	speed        = int(time.Second)
	words        = []string{}
	columnsLocks = map[int]*sync.Mutex{}
)

func handleWordInput(screen tcell.Screen, input string) {
	width := maxWordsWidth(screen)
	for row, c := range columnsLocks {
		c.Lock()
	widthRange:
		for i := width; i >= 0; i-- {
			mainc, _, _, _ := screen.GetContent(i, row)
			if mainc == 32 {
				continue
			}
			// End of the possible word
			if mainc == rune(input[len(input)-1]) {
				// Check if the character is the end
				check, _, _, _ := screen.GetContent(i+1, row)
				if check != 32 {
					continue widthRange
				}
				k := 0
				for ip := len(input) - 1; ip >= 0; ip-- {
					cs, _, _, _ := screen.GetContent(i-k, row)
					if rune(input[ip]) != cs {
						continue widthRange
					}
					// If we end the search check if the word is complete
					if ip == 0 {
						check, _, _, _ := screen.GetContent(i-k-1, row)
						if check != 32 {
							continue widthRange
						}
					}
					k++
				}
				// We find the word
				for ip := len(input) - 1; ip >= 0; ip-- {
					screen.SetContent(i-ip, row, 32, []rune{}, style)
				}
				gamePoints += len(input)
				c.Unlock()
				return
			}
		}
		c.Unlock()
	}
}

func loadWords(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	var output []string
	if err := dec.Decode(&output); err != nil {
		return nil, err
	}
	return output, nil
}

func startWordColumns(screen tcell.Screen, n int) {
	width := maxWordsWidth(screen)
	go handleColumnsSpeed(screen)
	go renderColumnsSpeed(screen, n)
	for i := 0; i < n; i++ {
		columnsLocks[i] = &sync.Mutex{}
		writexy(screen, width+2, i, "|")
		go executeColumn(screen, i)
	}
}

func renderColumnsSpeed(screen tcell.Screen, max int) {
	for {
		if gameIsAlive {
			writexy(screen, 1, max+4, fmt.Sprintf("Game Speed: %d", speed))
		}
		time.Sleep(time.Second * 2)
	}
}

func handleColumnNewWords(screen tcell.Screen, row int) {
	// Check if we can add new words
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	word := words[r.Intn(len(words))]
	if spaceForWord(screen, row, len(word)) {
		if r.Intn(50) >= 3 {
			return
		}
		writey(screen, row, word)
	}
}

func spaceForWord(screen tcell.Screen, row, needed int) bool {
	for i := 1; i <= needed+1; i++ {
		mainc, _, _, _ := screen.GetContent(i, row)
		if mainc != 32 {
			return false
		}
	}
	return true
}

func executeColumn(screen tcell.Screen, row int) {
	for {
		time.Sleep(time.Duration(speed))
		if gameIsAlive {
			columnsLocks[row].Lock()
			moverow(screen, row)
			handleColumnNewWords(screen, row)
			columnsLocks[row].Unlock()
			screen.Show()
		}
	}
}

func handleColumnsSpeed(screen tcell.Screen) {
	for {
		if speed < 190000000 {
			return
		}
		speed -= int(time.Millisecond * 22)
		time.Sleep(time.Second * 2)
	}
}

func maxWordsWidth(screen tcell.Screen) int {
	width, _ := screen.Size()
	return width - 3
}
