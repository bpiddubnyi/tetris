package main

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/veandco/go-sdl2/img"
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

	if err = rndr.SetScale(2, 2); err != nil {
		log.Fatalln("error: failed to set scale:", err)
	}

	preloadTex(rndr,
		"blue.png",
		"cyan.png",
		"green.png",
		"magenta.png",
		"orange.png",
		"red.png",
		"yellow.png",
	)
	defer destroyPreloadedTex()

	blue := getTex("blue.png")
	defer blue.Destroy()

	pf := playfield{
		pos: position{
			x: 7, y: 30,
		},
		tet: newTetrimino(tetI),
	}
	pf.tet.pos = position{
		x: 3,
		y: 0,
	}
	pf.tet.rotation = 1

theLoop:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break theLoop
			}
		}

		// Reduce CPU usage
		if time.Since(lastFrame) < time.Second/targetFPS {
			time.Sleep(time.Until(lastFrame.Add(time.Second / targetFPS)))
		}

		if err = rndr.SetDrawColor(255, 255, 255, 255); err != nil {
			log.Fatalln("error: failed to set draw color:", err)
		}
		rndr.Clear()

		pf.update()
		pf.draw(rndr)

		rndr.Present()
		lastFrame = time.Now()
	}
}

func assetPath(name string) string {
	return path.Join(assetDir, name)
}

func textureFromFile(rndr *sdl.Renderer, name string) *sdl.Texture {
	p := assetPath(name)
	t, err := img.LoadTexture(rndr, p)
	if err != nil {
		panic(fmt.Errorf("error: failed to load texture from %s: %v", p, err))
	}
	return t
}

func preloadTex(rndr *sdl.Renderer, names ...string) {
	if texStor == nil {
		texStor = map[string]*sdl.Texture{}
	}
	for _, f := range names {
		texStor[f] = textureFromFile(rndr, f)
	}
}

func destroyPreloadedTex() {
	for k, v := range texStor {
		v.Destroy()
		delete(texStor, k)
	}
}

func getTex(name string) *sdl.Texture {
	return texStor[name]
}
