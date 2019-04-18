package main

import (
	"math/rand"
	"time"
)

var (
	tetriminos = [...]tetrimino{
		// 0, I, cyan
		tetrimino{
			{
				{0, 0, 0, 0},
				{1, 1, 1, 1},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			{
				{0, 0, 1, 0},
				{0, 0, 1, 0},
				{0, 0, 1, 0},
				{0, 0, 1, 0},
			},
			{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{1, 1, 1, 1},
				{0, 0, 0, 0},
			},
			{
				{0, 1, 0, 0},
				{0, 1, 0, 0},
				{0, 1, 0, 0},
				{0, 1, 0, 0},
			},
		},
		// 1, J, blue
		tetrimino{
			{
				{2, 0, 0},
				{2, 2, 2},
				{0, 0, 0},
			},
			{
				{0, 2, 2},
				{0, 2, 0},
				{0, 2, 0},
			},
			{
				{0, 0, 0},
				{2, 2, 2},
				{0, 0, 2},
			},
			{
				{0, 2, 0},
				{0, 2, 0},
				{2, 2, 0},
			},
		},
		// 2, L, orange
		tetrimino{
			{
				{0, 0, 3},
				{3, 3, 3},
				{0, 0, 0},
			},
			{
				{0, 3, 0},
				{0, 3, 0},
				{0, 3, 3},
			},
			{
				{0, 0, 0},
				{3, 3, 3},
				{3, 0, 0},
			},
			{
				{3, 3, 0},
				{0, 3, 0},
				{0, 3, 0},
			},
		},
		// 3, O, square, yellow
		tetrimino{
			{
				{0, 4, 4, 0},
				{0, 4, 4, 0},
			},
			{
				{0, 4, 4, 0},
				{0, 4, 4, 0},
			},
			{
				{0, 4, 4, 0},
				{0, 4, 4, 0},
			},
			{
				{0, 4, 4, 0},
				{0, 4, 4, 0},
			},
		},
		// 4, S, green
		tetrimino{
			{
				{0, 5, 5},
				{5, 5, 0},
			},
			{
				{0, 5, 0},
				{0, 5, 5},
				{0, 0, 5},
			},
			{
				{0, 0, 0},
				{0, 5, 5},
				{5, 5, 0},
			},
			{
				{5, 0, 0},
				{5, 5, 0},
				{0, 5, 0},
			},
		},
		// 5, T, magenta
		tetrimino{
			{
				{0, 6, 0},
				{6, 6, 6},
				{0, 0, 0},
			},
			{
				{0, 6, 0},
				{0, 6, 6},
				{0, 6, 0},
			},
			{
				{0, 0, 0},
				{6, 6, 6},
				{0, 6, 0},
			},
			{
				{0, 6, 0},
				{6, 6, 0},
				{0, 6, 0},
			},
		},
		// 6, Z, red
		tetrimino{
			{
				{7, 7, 0},
				{0, 7, 7},
				{0, 0, 0},
			},
			{
				{0, 0, 7},
				{0, 7, 7},
				{0, 7, 0},
			},
			{
				{0, 0, 0},
				{7, 7, 0},
				{0, 7, 7},
			},
			{
				{0, 7, 0},
				{7, 7, 0},
				{7, 0, 0},
			},
		},
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type tShape [][]int8

type tetrimino [4]tShape

// type tetrimino struct {
// 	tex      *sdl.Texture
// 	shapes   *[4]shape
// 	rotation int8
// 	pos      position
// }

// func (r *tetrimino) rotate() {
// 	r.rotation++
// 	if r.rotation == 4 {
// 		r.rotation = 0
// 	}
// }

// func (r tetrimino) shape() shape {
// 	return r.shapes[r.rotation]
// }

// func newTet(t int) *tetrimino {
// 	return &tetrimino{
// 		tex:      getTex(tetTextures[t]),
// 		shapes:   &tetShapes[t],
// 		rotation: 0,
// 	}
// }

func newRandomTetrimino() *tetrimino {
	return &tetriminos[rand.Intn(len(tetriminos))]
}
