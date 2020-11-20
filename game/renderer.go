package game

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var colorTable map[string]uint32 = map[string]uint32{
	"Background": 0x00fbf8ef,
	"GameBoard":  0x00bbada0,
	"Tile0":      0xc4eee4da,
	"Tile2":      0x00eee4da,
	"Tile4":      0x00ede0c8,
	"Tile8":      0x00f2b179,
	"Tile16":     0x00f59563,
	"Tile32":     0x00f67c5f,
	"Tile64":     0x00f65e3b,
	"Tile128":    0x00edcf72,
	"Tile256":    0x00edcc61,
	"Tile512":    0x00edc850,
	"Tile1024":   0x00edc53f,
	"Tile2048":   0x00edc22e,
}

type Renderer struct {
	window   *sdl.Window
	surface  *sdl.Surface
	buffer   *sdl.Surface
	renderer *sdl.Renderer
	font     *ttf.Font
}

func NewRenderer(width int32, height int32, fontPath string, fontSize int, windowTitle string) *Renderer {
	var err error

	renderer := Renderer{
		window:   nil,
		surface:  nil,
		buffer:   nil,
		renderer: nil,
		font:     nil,
	}

	// Initialize sdl2
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Printf("Failed to init sdl2: %s\n", err)
		panic(err)
	}

	// Initialize font
	if err = ttf.Init(); err != nil {
		fmt.Printf("Failed to init font: %s\n", err)
		panic(err)
	}

	// Load the font for our text
	if renderer.font, err = ttf.OpenFont(fontPath, fontSize); err != nil {
		fmt.Printf("Failed to load font: %s\n", err)
		panic(err)
	}

	// Create window
	renderer.window, err = sdl.CreateWindow(windowTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("Failed to create window: %s\n", err)
		panic(err)
	}

	// Create draw surface and draw buffer
	if renderer.surface, err = renderer.window.GetSurface(); err != nil {
		fmt.Printf("Failed to get window surface: %s\n", err)
		panic(err)
	}

	if renderer.buffer, err = renderer.surface.Convert(renderer.surface.Format, renderer.window.GetFlags()); err != nil {
		fmt.Printf("Failed to create buffer: %s\n", err)
	}

	// Create renderer
	renderer.renderer, err = sdl.CreateRenderer(renderer.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("Failed to create renderer: %s\n", err)
		panic(err)
	}

	return &renderer
}

// Private draw functions
func (renderer *Renderer) drawString(x int, y int, str string, color *sdl.Color) {
	if len(str) <= 0 {
		return
	}

	// Create text
	text, err := renderer.font.RenderUTF8Blended(str, *color)
	if err != nil {
		fmt.Printf("Failed to create text: %s\n", err)
		return
	}
	defer text.Free()

	// Draw text, noted that you should always draw on buffer instead of directly draw on screen and blow yourself up.
	if err = text.Blit(nil, renderer.buffer, &sdl.Rect{X: int32(x), Y: int32(y)}); err != nil {
		fmt.Printf("Failed to draw text: %s\n", err)
		return
	}
}

// Draw a full sprite.
func (renderer *Renderer) drawSprite(x int, y int, sprite *sdl.Surface) {
	sprite.Blit(nil, renderer.buffer, &sdl.Rect{X: int32(x), Y: int32(y)})
}

// Draw a part of sprite.
func (renderer *Renderer) drawPartialSprite(dstX int, dstY int, sprite *sdl.Surface, srcX int, srcY int, w int, h int) {
	dstRect := sdl.Rect{X: int32(dstX), Y: int32(dstY), W: int32(w), H: int32(h)}
	srcRect := sdl.Rect{X: int32(srcX), Y: int32(srcY), W: int32(w), H: int32(h)}
	sprite.Blit(&srcRect, renderer.buffer, &dstRect)
}

func (renderer *Renderer) Update(gameBoard GameBoard) {
	// Render game.
	// 0xAARRGGBB
	// turningSnake := make([]([2]int), 0)
	// turningSnake = append(turningSnake, gameBoard.Snake[0])

	// lastDirection := [2]int{int(math.Abs(float64(gameBoard.Snake[0][0] - gameBoard.Snake[1][0]))), int(math.Abs(float64(gameBoard.Snake[0][1] - gameBoard.Snake[1][1])))}
	// lastSnake := gameBoard.Snake[1]

	// for _, i := range gameBoard.Snake[2:] {
	// 	direction := [2]int{int(math.Abs(float64(i[0] - lastSnake[0]))), int(math.Abs(float64(i[1] - lastSnake[1])))}
	// 	if lastDirection != direction {
	// 		turningSnake = append(turningSnake, i)
	// 	}
	// 	lastDirection = direction
	// 	lastSnake = i
	// }
	scale := renderer.buffer.W / int32(gameBoard.Width)

	for _, i := range gameBoard.Snake {
		renderer.buffer.FillRect(&sdl.Rect{X: int32(i[0]) * scale, Y: int32(i[1]) * scale, W: scale, H: scale}, 0xFFFFFFFF)
	}

	renderer.buffer.FillRect(&sdl.Rect{X: int32(gameBoard.Food[0]) * scale, Y: int32(gameBoard.Food[1]) * scale, W: scale, H: scale}, 0xFFFFFFFF)

	// Swap buffer and present our rendered content.
	renderer.window.UpdateSurface()
	renderer.buffer.Blit(nil, renderer.surface, nil)

	// Clear out buffer for next render round.
	renderer.buffer.FillRect(nil, 0xFF000000)
	renderer.renderer.Clear()
}
