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

	difficulty int
	exitFlag   bool
	inputLock  bool
	lastKey    int
	keys       map[int]bool
	frames     int
}

func NewController(width int32, height int32, fontPath string, fontSize int, windowTitle string) *Controller {
	controller := Controller{
		Renderer:  NewRenderer(width, height, fontPath, fontSize, windowTitle),
		GameBoard: NewGameBoard(64, 48),

		difficulty: 32,
		exitFlag:   false,
		lastKey:    sdl.K_RIGHT,
		keys:       make(map[int]bool),
	}

	for _, i := range []int{sdl.K_LEFT, sdl.K_RIGHT, sdl.K_DOWN, sdl.K_UP} {
		controller.keys[i] = true
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
					controller.GameBoard = NewGameBoard(64, 48)
					controller.lastKey = sdl.K_RIGHT
					// Add game difficulty
				} else if t.Keysym.Sym == sdl.K_EQUALS {
					if controller.difficulty != 1 {
						controller.difficulty /= 2
					}
				} else if t.Keysym.Sym == sdl.K_MINUS {
					if controller.difficulty < 256 {
						controller.difficulty *= 2
					}
				} else {
					if controller.keys[int(t.Keysym.Sym)] {
						controller.lastKey = int(t.Keysym.Sym)
					}
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

	if controller.frames%controller.difficulty == 0 {
		controller.GameBoard.Update(sdl.Keycode(controller.lastKey))
	}

	controller.Renderer.Update(*controller.GameBoard)

	controller.exitFlag = false
}

func (controller *Controller) Start() {
	startTime := time.Now()
	passedTime := time.Since(startTime)

	controller.Update()
	// Update the game by difficulty...
	for !controller.exitFlag {
		startTime = time.Now()
		controller.Update()
		passedTime = time.Since(startTime)
		// lock 60fps
		if passedTime < time.Duration(time.Second.Nanoseconds()/60) {
			time.Sleep(time.Duration(time.Second.Nanoseconds()/60) - passedTime)
		}
		controller.frames = (controller.frames + 1) % 360
	}
}
