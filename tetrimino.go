package main

import (
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	tetI = iota
	tetJ
	tetL
	tetO
	tetS
	tetT
	tetZ
)

var (
	tetTextures = [...]string{
		"cyan.png",
		"blue.png",
		"orange.png",
		"yellow.png",
		"green.png",
		"magenta.png",
		"red.png",
	}

	tetShapes = [...][4]shape{
		// 0, I, cyan
		[4]shape{
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
		[4]shape{
			{
				{1, 0, 0},
				{1, 1, 1},
				{0, 0, 0},
			},
			{
				{0, 1, 1},
				{0, 1, 0},
				{0, 1, 0},
			},
			{
				{0, 0, 0},
				{1, 1, 1},
				{0, 0, 1},
			},
			{
				{0, 1, 0},
				{0, 1, 0},
				{1, 1, 0},
			},
		},
		// 2, L, orange
		[4]shape{
			{
				{0, 0, 1},
				{1, 1, 1},
				{0, 0, 0},
			},
			{
				{0, 1, 0},
				{0, 1, 0},
				{0, 1, 1},
			},
			{
				{0, 0, 0},
				{1, 1, 1},
				{1, 0, 0},
			},
			{
				{1, 1, 0},
				{0, 1, 0},
				{0, 1, 0},
			},
		},
		// 3, O, square, yellow
		[4]shape{
			{
				{0, 1, 1, 0},
				{0, 1, 1, 0},
			},
			{
				{0, 1, 1, 0},
				{0, 1, 1, 0},
			},
			{
				{0, 1, 1, 0},
				{0, 1, 1, 0},
			},
			{
				{0, 1, 1, 0},
				{0, 1, 1, 0},
			},
		},
		// 4, S, green
		[4]shape{
			{
				{0, 1, 1},
				{1, 1, 0},
			},
			{
				{0, 1, 0},
				{0, 1, 1},
				{0, 0, 1},
			},
			{
				{0, 0, 0},
				{0, 1, 1},
				{1, 1, 0},
			},
			{
				{1, 0, 0},
				{1, 1, 0},
				{0, 1, 0},
			},
		},
		// 5, T, magenta
		[4]shape{
			{
				{0, 1, 0},
				{1, 1, 1},
				{0, 0, 0},
			},
			{
				{0, 1, 0},
				{0, 1, 1},
				{0, 1, 0},
			},
			{
				{0, 0, 0},
				{1, 1, 1},
				{0, 1, 0},
			},
			{
				{0, 1, 0},
				{1, 1, 0},
				{0, 1, 0},
			},
		},
		// 6, Z, red
		[4]shape{
			{
				{1, 1, 0},
				{0, 1, 1},
				{0, 0, 0},
			},
			{
				{0, 0, 1},
				{0, 1, 1},
				{0, 1, 0},
			},
			{
				{0, 0, 0},
				{1, 1, 0},
				{0, 1, 1},
			},
			{
				{0, 1, 0},
				{1, 1, 0},
				{1, 0, 0},
			},
		},
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type shape [][]int8

type tetrimino struct {
	tex      *sdl.Texture
	shapes   *[4]shape
	rotation int8
	pos      position
}

func (r *tetrimino) rotate() {
	r.rotation++
	if r.rotation == 4 {
		r.rotation = 0
	}
}

func (r tetrimino) currentShape() shape {
	return r.shapes[r.rotation]
}

func newTetrimino(t int) *tetrimino {
	return &tetrimino{
		tex:      getTex(tetTextures[t]),
		shapes:   &tetShapes[t],
		rotation: 0,
	}
}

func newRandomTet() *tetrimino {
	return newTetrimino(rand.Intn(len(tetShapes)))
}
