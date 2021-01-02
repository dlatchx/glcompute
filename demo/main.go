package main

import (
	"github.com/adraenwan/glcompute/glc"

	//"fmt"
)

func main() {
	err := glc.Init()
	if err != nil {
		panic(err)
	}

	p := glc.NewProgram()
	err = p.LoadSrc("demo.glsl.comp")
	if err != nil {
		panic(err)
	}

	buf := glc.NewBufferStorage(65535, 4, glc.BUF_STREAM_COPY)
	var bufSlice []int32

	buf.Map(&bufSlice, glc.MAP_WRITE|glc.MAP_INVALIDATE_BUFFER)
	for i := 0; i < len(bufSlice); i++ {
		bufSlice[i] = int32(i)
	}
	if !buf.Unmap() {
		panic("unmap error")
	}

	buf.Bind(4)

	glc.Sync()

	const split = 256
	for i := 0; i < split; i++ {
		p.Dispatch(65535 / split, 1, 1)
		glc.Sync()
		//fmt.Println(i)
	}

	glc.Sync()

	buf.Map(&bufSlice, glc.MAP_READ)
	//for i := 0; i < len(bufSlice); i++ {
	//	fmt.Println(bufSlice[i])
	//}
	if !buf.Unmap() {
		panic("unmap error")
	}
}
