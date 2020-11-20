package game

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var gameTitle string

type Controller struct {
	Renderer  *Renderer
	GameBoard *GameBoard

	exitFlag  bool
	inputLock bool
}

func NewController(width int32, height int32, fontPath string, fontSize int, windowTitle string) *Controller {
	controller := Controller{
		Renderer:  NewRenderer(width, height, fontPath, fontSize, windowTitle),
		GameBoard: NewGameBoard(64, 48),

		exitFlag: false,
	}

	gameTitle = windowTitle

	return &controller
}

func (controller *Controller) Update() {
	keyState := sdl.GetKeyboardState()
	// Exit game.
	if keyState[sdl.SCANCODE_ESCAPE] != 0 {
		controller.exitFlag = true
		return
	}

	// Get inputs.
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			controller.exitFlag = true
			return
		case *sdl.KeyboardEvent:
			if !controller.inputLock {
				if t.Keysym.Sym == sdl.K_r { // Reset game.
					fmt.Println("Reseting game...")
					controller.GameBoard = NewGameBoard(4*10, 3*10)
				} else {
					controller.GameBoard.Update(t.Keysym.Sym)
				}
			}

			// Anti-jittering
			if t.Repeat > 0 {
				// Held.
				controller.inputLock = false
			} else {
				// Pressed once.
				if t.State == sdl.RELEASED {
					controller.inputLock = false
				} else if t.State == sdl.PRESSED {
					controller.inputLock = true
				}
			}
		}
	}

	controller.Renderer.Update(*controller.GameBoard)

	controller.exitFlag = false
}

func (controller *Controller) Start() {
	startTime := time.Now()
	passedTime := time.Now().Sub(startTime).Microseconds()

	controller.Update()
	// Update the game by difficulty...
	for !controller.exitFlag {
		startTime = time.Now()
		if passedTime > time.Microsecond.Microseconds() {
			// fmt.Println(passedTime)
			controller.Update()
			startTime = time.Now()
			passedTime = 0
		}
		passedTime += time.Now().Sub(startTime).Microseconds()
	}
}
