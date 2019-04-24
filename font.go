package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type font struct {
	font  *ttf.Font
	color sdl.Color
	cache map[string]*sdl.Texture
}

func newFont(name string, size int, color sdl.Color) (*font, error) {
	f, err := ttf.OpenFont(assetPath(name), size)
	if err != nil {
		return nil, err
	}
	return &font{font: f, color: color, cache: map[string]*sdl.Texture{}}, nil
}

func (f *font) free() {
	f.font.Close()
	for k, v := range f.cache {
		v.Destroy()
		delete(f.cache, k)
	}
}

func (f *font) sprint(r *sdl.Renderer, s string) (*sdl.Texture, int32, int32) {
	if tex, ok := f.cache[s]; ok {
		_, _, w, h, err := tex.Query()
		if err != nil {
			panic(err)
		}
		return tex, w, h
	}

	sur, err := f.font.RenderUTF8Solid(s, f.color)
	if err != nil {
		panic(err)
	}
	defer sur.Free()

	tex, err := r.CreateTextureFromSurface(sur)
	if err != nil {
		panic(err)
	}
	f.cache[s] = tex

	return tex, sur.W, sur.H
}

func (f *font) sprintf(r *sdl.Renderer, format string, a interface{}) (*sdl.Texture, int32, int32) {
	return f.sprint(r, fmt.Sprintf(format, a))
}
