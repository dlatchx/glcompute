package main

import (
	"github.com/adraenwan/glcompute/glc"

	"fmt"
	"time"
)

const BUFLEN = 65536 * 256 * 16

func main() {
	err := glc.Init()
	if err != nil {
		panic(err)
	}

	p := glc.LoadProgram("demo.comp.glsl")

	bufSlice := make([]int32, BUFLEN)
	for i := 0; i < len(bufSlice); i++ {
		bufSlice[i] = int32(i)
	}
	buf := glc.LoadBufferStorage(&bufSlice, glc.BUF_STREAM_COPY)

	glc.Sync()
	tic := time.Now()

	p.Call(BUFLEN/256, 1, 1, buf)

	glc.Sync()
	toc := time.Now()

	bufSlice2 := make([]int32, BUFLEN)
	buf.Download(&bufSlice2)
	for i := 0; i < len(bufSlice2); i++ {
		//fmt.Println(bufSlice2[i])
	}

	fmt.Println("dispatch time:", float64(toc.Sub(tic).Nanoseconds())*1e-6, "ms")
}
