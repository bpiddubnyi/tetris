package main

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	playfieldHeight          = 24
	playfieldWidth           = 10
	playfieldVisibleHeight   = playfieldHeight - invisiblePlayfieldHeight
	invisiblePlayfieldHeight = 4

	moveCooldown = time.Millisecond * 150
)

const (
	dirDown = iota
	dirLeft
	dirRight
)

type position struct {
	x, y int32
}

type block struct {
	tex *sdl.Texture
}

type playfield struct {
	matrix     [playfieldHeight][playfieldWidth]block
	pos        position
	tet        *tetrimino
	latestMove time.Time
}

func (p *playfield) placeTetrimino(t *tetrimino, pos position) bool {
	shape := t.currentShape()
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {

		}
	}
	return false
}

func (p *playfield) moveTet(dir int) bool {
	switch dir {
	case dirLeft:
		if p.tet.pos.x-1 >= 0 {
			p.tet.pos.x--
		}
	case dirRight:
		if p.tet.pos.x+int32(len(p.tet.currentShape()[0]))+1 <= playfieldWidth {
			p.tet.pos.x++
		}
	case dirDown:
		if p.tet.pos.y+int32(len(p.tet.currentShape()))+1 <= playfieldHeight {
			p.tet.pos.y++
		}
	}
	log.Println("tet x:", p.tet.pos.x, "y:", p.tet.pos.y)
	return true
}

func (p *playfield) rotateTet() bool {
	p.tet.rotate()
	return true
}

func (p *playfield) update() {
	if time.Since(p.latestMove) >= moveCooldown {
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_SPACE] == 1 {
			p.rotateTet()
			p.latestMove = time.Now()
		}
		if keys[sdl.SCANCODE_LEFT] == 1 {
			p.moveTet(dirLeft)
			p.latestMove = time.Now()
		}
		if keys[sdl.SCANCODE_RIGHT] == 1 {
			p.moveTet(dirRight)
			p.latestMove = time.Now()
		}
		if keys[sdl.SCANCODE_DOWN] == 1 {
			p.moveTet(dirDown)
			p.latestMove = time.Now()
		}
	}
}

func (p *playfield) draw(r *sdl.Renderer) {
	// Draw background
	_, _, tW, tH, err := getTex("blue.png").Query()
	if err != nil {
		panic(errors.Wrap(err, "failed to query texture properties"))
	}
	if err = r.SetDrawColor(0, 0, 0, 255); err != nil {
		panic(errors.Wrap(err, "failed to set draw color"))
	}
	r.FillRect(&sdl.Rect{
		X: p.pos.x,
		Y: p.pos.y,
		W: tW * playfieldWidth,
		H: tH * playfieldVisibleHeight,
	})

	// Draw playfield
	// Start from the first visible row
	for i := invisiblePlayfieldHeight; i < playfieldHeight; i++ {
		for j := 0; j < playfieldWidth; j++ {
			var t *sdl.Texture
			if p.matrix[i][j].tex != nil {
				// Draw  block from the matrix
				t = p.matrix[i][j].tex
			} else if j >= int(p.tet.pos.x) &&
				j < int(p.tet.pos.x)+len(p.tet.currentShape()[0]) &&
				i >= int(p.tet.pos.y) &&
				i < int(p.tet.pos.y)+len(p.tet.currentShape()) &&
				p.tet.currentShape()[i-int(p.tet.pos.y)][j-int(p.tet.pos.x)] == 1 {
				// or from the mapped player's tetrimino
				t = p.tet.tex
			}
			if t != nil {
				r.Copy(t,
					&sdl.Rect{X: 0, Y: 0, W: tW, H: tH},
					&sdl.Rect{
						X: p.pos.x + int32(j)*tW,
						Y: p.pos.y + int32(i-invisiblePlayfieldHeight)*tH,
						W: tW, H: tH,
					})
			}
		}
	}
}
