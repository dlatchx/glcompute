package main

import (
	"github.com/dlatchx/glcompute/glc"

	"fmt"
	"time"
)

const BUFLEN = 65536 * 256 * 8

func init() {
	glc.Init()
}

func main() {
	p := glc.LoadProgram("demo.comp.glsl")

	// generate input data
	bufSlice := make([]float32, BUFLEN)
	for i := 0; i < BUFLEN; i++ {
		x := float32(BUFLEN - i - 1)
		bufSlice[i] = x * x
	}

	// send it to the GPU's dedicated memory
	tic_all := time.Now()
	buf := glc.LoadBufferStorage(&bufSlice, glc.BUF_STREAM_COPY)

	glc.Sync()
	tic_compute := time.Now()

	// run the shader
	p.Call(BUFLEN/256, 1, 1, buf)

	glc.Sync()
	toc_compute := time.Now()

	// retreive result
	bufSlice2 := make([]float32, BUFLEN)
	buf.Download(&bufSlice2)
	toc_all := time.Now()

	// show start and end of computed buffer
	for i := 0; i < 16; i++ {
		fmt.Println(bufSlice2[i])
	}
	fmt.Println("...")
	for i := BUFLEN - 16; i < BUFLEN; i++ {
		fmt.Println(bufSlice2[i])
	}

	fmt.Println("compute time:", float64(toc_compute.Sub(tic_compute).Nanoseconds())*1e-6, "ms")
	fmt.Println("all time:", float64(toc_all.Sub(tic_all).Nanoseconds())*1e-6, "ms")
}
