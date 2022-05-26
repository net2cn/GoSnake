package main

import (
	"fmt"

	"github.com/net2cn/GoSnake/game"
)

// SDL2 variables
var windowTitle string = "GoSnake SDL2"
var fontPath string = "./assets/UbuntuMono-R.ttf" // Man, do not use a variable-width font! It looks too ugly with that!
var fontSize int = 20
var windowWidth, windowHeight int32 = 640, 640

func main() {
	fmt.Println(windowTitle)
	fmt.Println("Yet another snake game written in golang.")

	game := game.NewController(windowWidth, windowHeight, fontPath, fontSize, windowTitle)
	game.Start()

	fmt.Println("Game exited. Bye!")
}
