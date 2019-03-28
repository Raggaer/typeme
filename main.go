package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
)

var (
	gameIsAlive  = true
	gameStatus   = 0
	gamePoints   = 0
	gameRows     = 10
	currentInput = ""
)

func main() {
	// Load game words
	word_list, err := loadWords("words.json")
	if err != nil {
		fmt.Println("Unable to load word list (words.json):")
		fmt.Println(err.Error())
		return
	}
	words = word_list

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("Unable to create new screen:")
		fmt.Println(err.Error())
		return
	}
	if err := screen.Init(); err != nil {
		fmt.Println("Unable to initialize screen:")
		fmt.Println(err.Error())
		return
	}
	screen.Show()
	quit := make(chan struct{})

	// Handle input keys
	go handleInput(screen, quit)

	// Start application
	go startApplication(screen)

	<-quit
	screen.Fini()
}

func startApplication(screen tcell.Screen) {
	writeln(screen, "Welcome to TypeMe. Press Enter to start!")
}

func handleInput(screen tcell.Screen, quit chan struct{}) {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEnter:
				if gameStatus == 0 {
					gameIsAlive = true
					gameStatus = 1
					row = 0
					screen.Clear()
					screen.Sync()
					go startWordColumns(screen, gameRows)
					go handleGameInput(screen, gameRows)
					go handleGameScore(screen, gameRows)
					go handleGameOver(screen)
				} else {
					clearRow(screen, 13)
					handleWordInput(screen, currentInput)
					currentInput = ""
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(currentInput) > 0 {
					clearRow(screen, 13)
					currentInput = currentInput[0 : len(currentInput)-1]
				}
			case tcell.KeyEscape:
				close(quit)
				return
			case tcell.KeyRune:
				if gameStatus == 1 {
					currentInput += string(ev.Rune())
				}
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}

func handleGameOver(screen tcell.Screen) {
	for {
		if !gameIsAlive {
			gameStatus = 0
			screen.Clear()
			screen.Sync()
			writeln(screen, fmt.Sprintf("Game Over!. You achieved an score of %d. Press Enter to try again", gamePoints))
			gamePoints = 0
			speed = int(time.Second)
			return
		}
		time.Sleep(time.Duration(speed))
	}
}

func handleGameInput(screen tcell.Screen, max int) {
	for {
		if gameIsAlive {
			writexy(screen, 1, max+3, fmt.Sprintf("Input: %s", currentInput))
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func handleGameScore(screen tcell.Screen, max int) {
	for {
		if gameIsAlive {
			writexy(screen, 1, max+5, fmt.Sprintf("Game Score: %d", gamePoints))
		}
		time.Sleep(time.Duration(speed))

	}
}
