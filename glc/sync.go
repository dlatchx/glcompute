package glc

import (
	"time"

	"github.com/go-gl/gl/v4.3-core/gl"
)

func Barrier() {
	gl.MemoryBarrier(gl.ALL_BARRIER_BITS)
}

func Flush() {
	gl.Flush()
}

func Sync() {
	Barrier()
	Flush()
	time.Sleep(0)
	gl.Finish()
}
