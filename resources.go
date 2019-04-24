package main

import (
	"fmt"
	"path"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	blockNum = 8

	fntName = "Xolonium-Regular.ttf"
	fntSize = 24
)

var (
	tetTextureNames = [blockNum]string{
		"white.png",
		"cyan.png",
		"blue.png",
		"orange.png",
		"yellow.png",
		"green.png",
		"magenta.png",
		"red.png",
	}
)

type resources struct {
	tex    [blockNum]*sdl.Texture
	font   *font
	tW, tH int32
}

func assetPath(name string) string {
	return path.Join(assetDir, name)
}

func textureFromFile(rndr *sdl.Renderer, name string) (*sdl.Texture, error) {
	p := assetPath(name)
	t, err := img.LoadTexture(rndr, p)
	if err != nil {
		return nil, fmt.Errorf("failed to load texture from %s: %v", p, err)
	}
	return t, nil
}

func loadGameResources(r *sdl.Renderer) (*resources, error) {
	var (
		res resources
		err error
	)

	for i, n := range tetTextureNames {
		res.tex[i], err = textureFromFile(r, n)
		if err != nil {
			return nil, err
		}
	}
	_, _, tW, tH, err := res.tex[0].Query()
	if err != nil {
		return nil, fmt.Errorf("failed to query texture: %s", err)
	}

	res.tW = tW
	res.tH = tH

	res.font, err = newFont(fntName, fntSize, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		return nil, fmt.Errorf("failed to open font: %s", err)
	}

	return &res, nil
}

func (r *resources) free() {
	for _, t := range r.tex {
		_ = t.Destroy()
	}
	r.font.free()
}
