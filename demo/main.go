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

	buf := glc.NewBufferStorage(65535, 4, glc.BUF_STREAM_COPY)

	bufSlice := make([]int32, 65535)
	for i := 0; i < len(bufSlice); i++ {
		bufSlice[i] = int32(i)
	}
	buf.Upload(bufSlice)

	buf.Bind(1)

	glc.Sync()

	tic := time.Now()

	const split = 256
	for i := 0; i < split; i++ {
		p.Dispatch(65535/split, 1, 1)
		glc.Sync()
		//fmt.Println(i)
	}

	glc.Sync()

	toc := time.Now()

	bufSlice2 := make([]int32, 65535)
	buf.Download(bufSlice2)
	for i := 0; i < len(bufSlice2); i++ {
		fmt.Println(bufSlice[i])
	}

	fmt.Println("dispatch time:", toc.Sub(tic).Milliseconds(), "ms")
}
