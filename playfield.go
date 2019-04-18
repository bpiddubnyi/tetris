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

func (p *playfield) tetMoveIsPossible(pos position, rot int8) bool {
	shape := p.tet.shapes[rot]
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {
			if shape[i][j] == 1 && (int(pos.y)+i < 0 ||
				int(pos.y)+i >= playfieldHeight ||
				int(pos.x)+j < 0 ||
				int(pos.x)+j >= playfieldWidth ||
				p.matrix[int(pos.y)+i][int(pos.x)+j].tex != nil) {
				return false
			}
		}
	}
	return true
}

func (p *playfield) moveTet(dir int) bool {
	switch dir {
	case dirLeft:
		if p.tetMoveIsPossible(position{x: p.tet.pos.x - 1, y: p.tet.pos.y}, p.tet.rotation) {
			p.tet.pos.x--
		}
	case dirRight:
		if p.tetMoveIsPossible(position{x: p.tet.pos.x + 1, y: p.tet.pos.y}, p.tet.rotation) {
			p.tet.pos.x++
		}
	case dirDown:
		if p.tetMoveIsPossible(position{x: p.tet.pos.x, y: p.tet.pos.y + 1}, p.tet.rotation) {
			p.tet.pos.y++
		}
	}
	log.Println("tet x:", p.tet.pos.x, "y:", p.tet.pos.y)
	return true
}

func (p *playfield) rotateTet() bool {
	nR := p.tet.rotation + 1
	if nR == 4 {
		nR = 0
	}
	if p.tetMoveIsPossible(p.tet.pos, nR) {
		p.tet.rotation = nR
		return true
	}
	return false
}

func (p *playfield) update() {
	if time.Since(p.latestMove) >= moveCooldown {
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_SPACE] == 1 || keys[sdl.SCANCODE_UP] == 1 {
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
