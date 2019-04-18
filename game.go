package main

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	playfieldHeight          = 26
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
	playfield  [playfieldHeight][playfieldWidth]int8
	tet        *tetrimino
	tetPos     position
	tetRot     int8
	latestMove time.Time
	latestFall time.Time
	speed      float32
}

func newGame() game {
	return game{
		speed: 2,
		tet:   newRandomTetrimino(),
	}
}

func (g *game) tetMoveIsPossible(pos position, rot int8) bool {
	shape := g.tet[rot]
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {
			if shape[i][j] > 0 && (int(pos.y)+i < 0 ||
				int(pos.y)+i >= playfieldHeight ||
				int(pos.x)+j < 0 ||
				int(pos.x)+j >= playfieldWidth ||
				g.playfield[int(pos.y)+i][int(pos.x)+j] > 0) {
				return false
			}
		}
	}
	return true
}

func (g *game) moveTet(dir int) bool {
	switch dir {
	case dirLeft:
		if g.tetMoveIsPossible(position{x: g.tetPos.x - 1, y: g.tetPos.y}, g.tetRot) {
			g.tetPos.x--
		}
	case dirRight:
		if g.tetMoveIsPossible(position{x: g.tetPos.x + 1, y: g.tetPos.y}, g.tetRot) {
			g.tetPos.x++
		}
	case dirDown:
		if g.tetMoveIsPossible(position{x: g.tetPos.x, y: g.tetPos.y + 1}, g.tetRot) {
			g.tetPos.y++
		} else {
			return false
		}
	}
	return true
}

func (g *game) rotateTet() bool {
	nR := g.tetRot + 1
	if nR == 4 {
		nR = 0
	}
	if g.tetMoveIsPossible(g.tetPos, nR) {
		g.tetRot = nR
		return true
	}
	return false
}

func (g *game) moveLinesDown(start, num int) {
	copy(g.playfield[num:start+1], g.playfield[0:start-num+1])
	for i := 0; i < num; i++ {
		for j := 0; j < len(g.playfield[i]); j++ {
			g.playfield[i][j] = 0
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
			if g.playfield[i][j] == 0 {
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
	shape := g.tet[g.tetRot]
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {
			if shape[i][j] > 0 {
				g.playfield[int(g.tetPos.y)+i][int(g.tetPos.x)+j] = shape[i][j]
			}
		}
	}
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
			g.clearLines()
			g.tet = newRandomTetrimino()
			g.tetPos = position{x: 0, y: 0}
			g.tetRot = 0
		}
		g.latestFall = time.Now()
	}
}

func (g *game) draw(r *sdl.Renderer, p position, res *resources) {
	// Draw playfield
	// Start from the first visible row
	for i := invisiblePlayfieldHeight; i < playfieldHeight; i++ {
		for j := 0; j < playfieldWidth; j++ {
			var t int8
			shape := g.tet[g.tetRot]
			if g.playfield[i][j] > 0 {
				// Draw  block from the matrix
				t = g.playfield[i][j]
			} else if j >= int(g.tetPos.x) &&
				j < int(g.tetPos.x)+len(shape[0]) &&
				i >= int(g.tetPos.y) &&
				i < int(g.tetPos.y)+len(shape) &&
				shape[i-int(g.tetPos.y)][j-int(g.tetPos.x)] > 0 {
				// or from the mapped player's tetrimino
				t = shape[i-int(g.tetPos.y)][j-int(g.tetPos.x)]
			}

			r.Copy(res.tex[t],
				&sdl.Rect{X: 0, Y: 0, W: res.tW, H: res.tH},
				&sdl.Rect{
					X: p.x + int32(j)*res.tW,
					Y: p.y + int32(i-invisiblePlayfieldHeight)*res.tH,
					W: res.tW, H: res.tH,
				})
		}
	}
}
