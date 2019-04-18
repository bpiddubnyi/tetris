package main

import "github.com/veandco/go-sdl2/sdl"

type kbd struct {
	keys []uint8
}

func (k *kbd) poll() {
	sdl.PumpEvents()
	keys := sdl.GetKeyboardState()

	if len(k.keys) == 0 {
		k.keys = make([]uint8, len(keys))
	}

	for i, key := range keys {
		if k.keys[i] == 0 {
			k.keys[i] = key
		} else {
			if keys[i] == 0 {
				k.keys[i] = 0
			} else {
				k.keys[i] = 2
			}
		}
	}
}

func (k kbd) pressed(key uint8) bool {
	return k.keys[key] > 0
}

func (k *kbd) justPressed(key int8) bool {
	return k.keys[key] == 1
}
