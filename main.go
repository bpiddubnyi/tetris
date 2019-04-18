package main

import (
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 600
	windowHeight = 800
	windowTitle  = "Tetris"
	assetDir     = "assets"

	targetFPS = 60
)

var (
	texStor   map[string]*sdl.Texture
	lastFrame time.Time
)

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	window, err := sdl.CreateWindow(
		windowTitle,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		log.Fatalln("error: failed to create window:", err)
	}
	defer window.Destroy()

	rndr, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalln("error: failed to create renderer:", err)
	}
	defer rndr.Destroy()

	if err = rndr.SetScale(1, 1); err != nil {
		log.Fatalln("error: failed to set scale:", err)
	}

	res, err := loadGameResources(rndr)
	if err != nil {
		log.Fatalf("error: failed to load game resources: %s", err)
	}
	kbd := kbd{}
	game := newGame()

theLoop:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break theLoop
			}
		}

		if err = rndr.SetDrawColor(0, 0, 0, 255); err != nil {
			log.Fatalln("error: failed to set draw color:", err)
		}
		rndr.Clear()

		game.update(&kbd)
		game.draw(rndr, position{x: 20, y: 30}, res)

		// Reduce CPU usage
		if time.Since(lastFrame) < time.Second/targetFPS {
			time.Sleep(time.Until(lastFrame.Add(time.Second / targetFPS)))
		}

		rndr.Present()
		lastFrame = time.Now()
	}
}
