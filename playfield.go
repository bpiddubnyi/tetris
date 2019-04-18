package main

import (
	"time"

	"github.com/pkg/errors"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	playfieldHeight          = 24
	playfieldWidth           = 10
	playfieldVisibleHeight   = playfieldHeight - invisiblePlayfieldHeight
	invisiblePlayfieldHeight = 3

	tick = time.Millisecond * 100
)

const (
	dirDown = iota
	dirLeft
	dirRight
)

type position struct {
	x, y int32
}

type game struct {
	playfield  [playfieldHeight][playfieldWidth]*sdl.Texture
	pos        position
	tet        *tetrimino
	latestMove time.Time
	latestFall time.Time
	speed      float32
}

func newGame(p position) game {
	return game{
		pos:   p,
		speed: 2,
		tet:   newRandomTet(),
	}
}

func (g *game) tetMoveIsPossible(pos position, rot int8) bool {
	shape := g.tet.shapes[rot]
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {
			if shape[i][j] == 1 && (int(pos.y)+i < 0 ||
				int(pos.y)+i >= playfieldHeight ||
				int(pos.x)+j < 0 ||
				int(pos.x)+j >= playfieldWidth ||
				g.playfield[int(pos.y)+i][int(pos.x)+j] != nil) {
				return false
			}
		}
	}
	return true
}

func (g *game) moveTet(dir int) bool {
	switch dir {
	case dirLeft:
		if g.tetMoveIsPossible(position{x: g.tet.pos.x - 1, y: g.tet.pos.y}, g.tet.rotation) {
			g.tet.pos.x--
		}
	case dirRight:
		if g.tetMoveIsPossible(position{x: g.tet.pos.x + 1, y: g.tet.pos.y}, g.tet.rotation) {
			g.tet.pos.x++
		}
	case dirDown:
		if g.tetMoveIsPossible(position{x: g.tet.pos.x, y: g.tet.pos.y + 1}, g.tet.rotation) {
			g.tet.pos.y++
		} else {
			return false
		}
	}
	return true
}

func (g *game) rotateTet() bool {
	nR := g.tet.rotation + 1
	if nR == 4 {
		nR = 0
	}
	if g.tetMoveIsPossible(g.tet.pos, nR) {
		g.tet.rotation = nR
		return true
	}
	return false
}

func (g *game) update() {
	if time.Since(g.latestMove) >= tick {
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_UP] == 1 {
			g.rotateTet()
			g.latestMove = time.Now()
		}
		if keys[sdl.SCANCODE_LEFT] == 1 {
			g.moveTet(dirLeft)
			g.latestMove = time.Now()
		}
		if keys[sdl.SCANCODE_RIGHT] == 1 {
			g.moveTet(dirRight)
			g.latestMove = time.Now()
		}
		if keys[sdl.SCANCODE_DOWN] == 1 {
			g.moveTet(dirDown)
			g.latestMove = time.Now()
		}
		if keys[sdl.SCANCODE_SPACE] == 1 {
			for g.moveTet(dirDown) {
			}
			g.latestMove = time.Now()
		}
	}

	if time.Since(g.latestFall) >= time.Duration(float32(tick)*g.speed) {
		if !g.moveTet(dirDown) {
			g.mergeTet()
			g.tet = newRandomTet()
			g.clearLines()
		}
		g.latestFall = time.Now()
	}
}

func (g *game) moveLinesDown(start, num int) {
	copy(g.playfield[num:start+1], g.playfield[0:start-num+1])
	for i := 0; i < num; i++ {
		for j := 0; j < len(g.playfield[i]); j++ {
			g.playfield[i][j] = nil
		}
	}
}

func (g *game) clearLines() {
	startLine := 0
	lines := 0
	// Scan from the bottom
lineLoop:
	for i := len(g.playfield) - 1; i >= 0; i-- {
		for j := 0; j < len(g.playfield[i]); j++ {
			if g.playfield[i][j] == nil {
				// This line is not complete
				if lines > 0 {
					g.moveLinesDown(startLine, lines)
					startLine, lines = 0, 0
				}
				continue lineLoop
			}
			if j+1 == len(g.playfield[i]) {
				// This is the last block in the line and it's not empty
				// the line is complete
				if lines == 0 {
					startLine = i
				}
				lines++
			}
		}
	}
	if lines > 0 {
		g.moveLinesDown(startLine, lines)
		startLine, lines = 0, 0
	}
}

func (g *game) mergeTet() {
	shape := g.tet.shapes[g.tet.rotation]
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {
			if shape[i][j] == 1 {
				g.playfield[int(g.tet.pos.y)+i][int(g.tet.pos.x)+j] = g.tet.tex
			}
		}
	}
}

func (g *game) draw(r *sdl.Renderer) {
	// Draw background
	back := getTex("white.png")
	_, _, tW, tH, err := back.Query()
	if err != nil {
		panic(errors.Wrap(err, "failed to query texture properties"))
	}
	if err = r.SetDrawColor(255, 255, 255, 255); err != nil {
		panic(errors.Wrap(err, "failed to set draw color"))
	}
	r.FillRect(&sdl.Rect{
		X: g.pos.x,
		Y: g.pos.y,
		W: tW * playfieldWidth,
		H: tH * playfieldVisibleHeight,
	})

	// Draw playfield
	// Start from the first visible row
	for i := invisiblePlayfieldHeight; i < playfieldHeight; i++ {
		for j := 0; j < playfieldWidth; j++ {
			var t *sdl.Texture
			shape := g.tet.shape()
			if g.playfield[i][j] != nil {
				// Draw  block from the matrix
				t = g.playfield[i][j]
			} else if j >= int(g.tet.pos.x) &&
				j < int(g.tet.pos.x)+len(shape[0]) &&
				i >= int(g.tet.pos.y) &&
				i < int(g.tet.pos.y)+len(shape) &&
				shape[i-int(g.tet.pos.y)][j-int(g.tet.pos.x)] == 1 {
				// or from the mapped player's tetrimino
				t = g.tet.tex
			} else {
				t = back
			}
			if t != nil {
				r.Copy(t,
					&sdl.Rect{X: 0, Y: 0, W: tW, H: tH},
					&sdl.Rect{
						X: g.pos.x + int32(j)*tW,
						Y: g.pos.y + int32(i-invisiblePlayfieldHeight)*tH,
						W: tW, H: tH,
					})
			}
		}
	}
}
