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

	gi1.Bind(0)
	gi2.Bind(1)
	p.Dispatch(uint32(im.Width()), uint32(im.Height()), 1)

	glc.Sync()

	(gi2).Download(nil).Show("truc")
}
