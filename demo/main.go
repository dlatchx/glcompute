package main

import (
	"github.com/adraenwan/glcompute/glc"

	"fmt"
	"time"
)

func main() {
	err := glc.Init()
	if err != nil {
		panic(err)
	}

	p := glc.NewProgram()
	err = p.LoadSrc("demo.comp.glsl")
	if err != nil {
		panic(err)
	}

	buf := glc.NewBufferStorage(65536, 4, glc.BUF_STREAM_COPY)

	bufSlice := make([]int32, 65536)
	for i := 0; i < len(bufSlice); i++ {
		bufSlice[i] = int32(i)
	}
	buf.Upload(&bufSlice)

	glc.Sync()
	tic := time.Now()

	p.Call(65536/4, 1, 1, buf)

	glc.Sync()
	toc := time.Now()

	bufSlice2 := make([]int32, 65536)
	buf.Download(&bufSlice2)
	for i := 0; i < len(bufSlice2); i++ {
		fmt.Println(bufSlice2[i])
	}

	fmt.Println("dispatch time:", toc.Sub(tic).Milliseconds(), "ms")
}
