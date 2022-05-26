package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type GameBoard struct {
	Snake  []([2]int)
	Food   [2]int
	Width  int
	Height int

	GameScore    int
	GameWin      bool
	GameOverFlag bool
}

func NewGameBoard(width int, height int) *GameBoard {
	fmt.Println("Creating new board.")
	gameBoard := GameBoard{}

	gameBoard.Snake = make([]([2]int), 0)
	gameBoard.Width = width
	gameBoard.Height = height
	gameBoard.GameOverFlag = false
	gameBoard.newSnake()
	gameBoard.generateNewFood()
	return &gameBoard
}

func (gameBoard *GameBoard) newSnake() {
	for i := 0; i < 3; i++ {
		gameBoard.Snake = append(gameBoard.Snake, [2]int{gameBoard.Width/4 + i, gameBoard.Height / 2})
	}
}

func Contains(slices [][2]int, item [2]int) int {
	for idx, slice := range slices {
		if slice == item {
			return idx
		}
	}

	return -1
}

func (gameBoard *GameBoard) generateNewFood() {
	rand.Seed(time.Now().Unix())
	newFood := [2]int{rand.Intn(gameBoard.Width), rand.Intn(gameBoard.Height)}
	for Contains(gameBoard.Snake, newFood) != -1 {
		newFood = [2]int{rand.Intn(gameBoard.Width), rand.Intn(gameBoard.Height)}
	}
	gameBoard.Food = newFood
}

func (gameBoard *GameBoard) move(tile [2]int) {
	hadFood := false
	nextTile := [2]int{(gameBoard.Snake[len(gameBoard.Snake)-1][0] + tile[0]) % gameBoard.Width, (gameBoard.Snake[len(gameBoard.Snake)-1][1] + tile[1]) % gameBoard.Height}

	if nextTile == gameBoard.Snake[len(gameBoard.Snake)-2] {
		tile[0] = -tile[0]
		tile[1] = -tile[1]
		nextTile = [2]int{(gameBoard.Snake[len(gameBoard.Snake)-1][0] + tile[0]) % gameBoard.Width, (gameBoard.Snake[len(gameBoard.Snake)-1][1] + tile[1]) % gameBoard.Height}
	}

	if nextTile[0] < 0 {
		nextTile[0] = gameBoard.Width - 1
	}
	if nextTile[1] < 0 {
		nextTile[1] = gameBoard.Height - 1
	}

	// Next tile is snake
	idx := Contains(gameBoard.Snake, nextTile)
	if idx != -1 {
		// Cut the snake in half at where the head and the body collided.
		gameBoard.Snake = gameBoard.Snake[idx+1:]
	}

	// Next tile is food
	if nextTile == gameBoard.Food {
		gameBoard.Snake = append([][2]int{nextTile}, gameBoard.Snake...)
		gameBoard.GameScore++
		gameBoard.generateNewFood()
		hadFood = true
	}

	gameBoard.Snake = append(gameBoard.Snake, nextTile)
	_ = hadFood
	// Pop the tail if no food is consumed.
	// Otherwise the length of the snake will automatically increase 1.
	gameBoard.Snake = gameBoard.Snake[1:]
}

func (gameBoard *GameBoard) Update(keyState sdl.Keycode) {
	if gameBoard.GameWin {
		return
	}

	var input [2]int

	switch keyState {
	case sdl.K_RIGHT:
		fmt.Println("Right pressed.")
		input = [2]int{1, 0}
	case sdl.K_UP:
		fmt.Println("Up pressed.")
		input = [2]int{0, -1}
	case sdl.K_LEFT:
		fmt.Println("Left pressed.")
		input = [2]int{-1, 0}
	case sdl.K_DOWN:
		fmt.Println("Down pressed.")
		input = [2]int{0, 1}
	}

	if len(gameBoard.Snake) == gameBoard.Height*gameBoard.Width {
		gameBoard.GameWin = true
		fmt.Println("You Win! There is no more tiles for you to go!")
	}

	if !gameBoard.GameOverFlag {
		gameBoard.move(input)
	} else {
		fmt.Println("Game over!")
	}
}
