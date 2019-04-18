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

type block struct {
	tex *sdl.Texture
}

type playfield struct {
	matrix     [playfieldHeight][playfieldWidth]block
	pos        position
	tet        *tetrimino
	latestMove time.Time
	latestFall time.Time
	speed      float32
}

func newPlayfield(p position) playfield {
	return playfield{
		pos:   p,
		speed: 2,
		tet:   newRandomTet(),
	}
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
		} else {
			return false
		}
	}
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
	if time.Since(p.latestMove) >= tick {
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_UP] == 1 {
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
		if keys[sdl.SCANCODE_SPACE] == 1 {
			for p.moveTet(dirDown) {
			}
			p.latestMove = time.Now()
		}
	}

	if time.Since(p.latestFall) >= time.Duration(float32(tick)*p.speed) {
		if !p.moveTet(dirDown) {
			p.mergeTet()
			p.tet = newRandomTet()
			p.clearLines()
		}
		p.latestFall = time.Now()
	}
}

func (p *playfield) moveLinesDown(start, num int) {
	copy(p.matrix[num:start+1], p.matrix[0:start-num+1])
	for i := 0; i < num; i++ {
		for j := 0; j < len(p.matrix[i]); j++ {
			p.matrix[i][j].tex = nil
		}
	}
}

func (p *playfield) clearLines() {
	startLine := 0
	lines := 0
	// Scan from the bottom
lineLoop:
	for i := len(p.matrix) - 1; i >= 0; i-- {
		for j := 0; j < len(p.matrix[i]); j++ {
			if p.matrix[i][j].tex == nil {
				// This line is not complete
				if lines > 0 {
					p.moveLinesDown(startLine, lines)
					startLine, lines = 0, 0
				}
				continue lineLoop
			}
			if j+1 == len(p.matrix[i]) {
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
		p.moveLinesDown(startLine, lines)
		startLine, lines = 0, 0
	}
}

func (p *playfield) mergeTet() {
	shape := p.tet.shapes[p.tet.rotation]
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {
			if shape[i][j] == 1 {
				p.matrix[int(p.tet.pos.y)+i][int(p.tet.pos.x)+j].tex = p.tet.tex
			}
		}
	}
}

func (p *playfield) draw(r *sdl.Renderer) {
	// Draw background, any Tet texture would do to calculate the size
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
