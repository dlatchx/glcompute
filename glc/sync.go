package glc

import (
	"time"

    gl "github.com/go-gl/gl/v3.1/gles2"
)

// memory barrier
func Barrier() {
	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)
}

// waits for instruction
// pipeline to be empty
func Flush() {
	gl.Flush()
}

// wait until al GPU tasks are done
// and returns when memory is synchronized
func Sync() {
	Barrier()
	Flush()
	time.Sleep(0)
	gl.Finish()
}
