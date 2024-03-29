package main

import (
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	playfieldHeight          = 26
	playfieldWidth           = 10
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

type t struct {
	position
	shape *tetrimino
	r     int8
}

func newT() t {
	return t{shape: getRandomTetrimino()}
}

type score struct {
	lines  int
	points int
	level  int
}

var (
	scores = [...]int{
		40, 100, 300, 1200,
	}
)

func (s *score) addLines(n int) {
	s.lines += n
	s.points += scores[n-1] * (s.level + 1)
}

func (s *score) draw(r *sdl.Renderer, p position, res *resources) {
	lines, lW, lH := res.font.sprintf(r, "lines: %d", s.lines)
	if err := r.Copy(lines,
		&sdl.Rect{X: 0, Y: 0, W: lW, H: lH},
		&sdl.Rect{X: p.x, Y: p.y, W: lW, H: lH}); err != nil {
		panic(err)
	}

	score, sW, sH := res.font.sprintf(r, "score: %d", s.points)
	if err := r.Copy(score,
		&sdl.Rect{X: 0, Y: 0, W: sW, H: sH},
		&sdl.Rect{X: p.x, Y: p.y + lH + 2, W: sW, H: sH}); err != nil {
		panic(err)
	}
}

type game struct {
	playfield  [playfieldHeight][playfieldWidth]int8
	t          t
	latestMove time.Time
	latestFall time.Time
	score      score
	speed      float32
	pause      bool
	loose      bool
}

func newGame() game {
	return game{
		speed: 2,
		t:     newT(),
	}
}

func (g *game) reset() {
	*g = newGame()
}

func (g *game) tetMoveIsPossible(pos position, rot int8) bool {
	shape := g.t.shape[rot]
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
		if g.tetMoveIsPossible(position{x: g.t.x - 1, y: g.t.y}, g.t.r) {
			g.t.x--
		}
	case dirRight:
		if g.tetMoveIsPossible(position{x: g.t.x + 1, y: g.t.y}, g.t.r) {
			g.t.x++
		}
	case dirDown:
		if g.tetMoveIsPossible(position{x: g.t.x, y: g.t.y + 1}, g.t.r) {
			g.t.y++
		} else {
			return false
		}
	}
	return true
}

func (g *game) rotateTet() bool {
	nR := g.t.r + 1
	if nR == 4 {
		nR = 0
	}
	if g.tetMoveIsPossible(g.t.position, nR) {
		g.t.r = nR
		return true
	}
	return false
}

func (g *game) collapseAndScore(start, n int) {
	copy(g.playfield[n:start+1], g.playfield[0:start-n+1])
	for i := 0; i < n; i++ {
		for j := 0; j < len(g.playfield[i]); j++ {
			g.playfield[i][j] = 0
		}
	}

	g.score.addLines(n)
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
					g.collapseAndScore(startLine, lines)
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
		g.collapseAndScore(startLine, lines)
	}
}

func (g *game) mergeTet() {
	shape := g.t.shape[g.t.r]
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(shape[i]); j++ {
			if shape[i][j] > 0 {
				g.playfield[int(g.t.y)+i][int(g.t.x)+j] = shape[i][j]
			}
		}
	}
}

func (g *game) update(kbd *kbd) {
	kbd.poll()

	if kbd.justPressed(sdl.SCANCODE_ESCAPE) {
		if !g.loose {
			g.pause = !g.pause
		} else {
			g.reset()
		}
	}

	if g.pause || g.loose {
		return
	}

	if kbd.justPressed(sdl.SCANCODE_UP) {
		g.rotateTet()
		g.latestMove = time.Now()
	}

	if kbd.justPressed(sdl.SCANCODE_SPACE) {
		for g.moveTet(dirDown) {
		}
		g.latestMove = time.Now()
	}

	if time.Since(g.latestMove) >= tick {
		if kbd.pressed(sdl.SCANCODE_LEFT) {
			g.moveTet(dirLeft)
			g.latestMove = time.Now()
		}
		if kbd.pressed(sdl.SCANCODE_RIGHT) {
			g.moveTet(dirRight)
			g.latestMove = time.Now()
		}
		if kbd.pressed(sdl.SCANCODE_DOWN) {
			g.moveTet(dirDown)
			g.latestMove = time.Now()
		}
	}

	if time.Since(g.latestFall) >= time.Duration(float32(tick)*g.speed) {
		if !g.moveTet(dirDown) {
			g.mergeTet()
			g.clearLines()
			g.t = newT()
			if !g.tetMoveIsPossible(g.t.position, 0) {
				g.loose = true
				// TODO: Draw this
				log.Println("you loose")
			}
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
			shape := g.t.shape[g.t.r]
			if g.playfield[i][j] > 0 {
				// Draw  block from the matrix
				t = g.playfield[i][j]
			} else if j >= int(g.t.x) &&
				j < int(g.t.x)+len(shape[0]) &&
				i >= int(g.t.y) &&
				i < int(g.t.y)+len(shape) &&
				shape[i-int(g.t.y)][j-int(g.t.x)] > 0 {
				// or from the mapped player's tetrimino
				t = shape[i-int(g.t.y)][j-int(g.t.x)]
			}

			if err := r.Copy(res.tex[t],
				&sdl.Rect{X: 0, Y: 0, W: res.tW, H: res.tH},
				&sdl.Rect{
					X: p.x + int32(j)*res.tW,
					Y: p.y + int32(i-invisiblePlayfieldHeight)*res.tH,
					W: res.tW, H: res.tH,
				}); err != nil {
				panic(err)
			}
		}
	}
}
