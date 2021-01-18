package main

import (
	"github.com/adraenwan/glcompute/glc"
	"github.com/adraenwan/glcompute/glc/img"
)

func init() {
	glc.Init()
}

func main() {
	p := glc.LoadProgram("sobel.glsl")
	defer p.Delete()

	im := img.LoadImg("toto.png")
	gi1 := im.Upload(nil)
	gi2 := img.NewGpuImg(im.Width(), im.Height())

	glc.Sync()

	p.Call(uint32(im.Width()), uint32(im.Height()), 1, gi1, gi2)

	glc.Sync()

	(gi2).Download(nil).Show("truc")
}
