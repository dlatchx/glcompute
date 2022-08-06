package img

import (
	"github.com/dlatchx/glcompute/glc"
)

// image metadata
type Img struct {
	width  int
	height int
}

func (i Img) Width() int {
	return i.width
}

func (i Img) Height() int {
	return i.height
}

// wrapper for cpu access of image
type CpuImg struct {
	Img

	Raster [][4]float32
}

func NewCpuImg(width, height int) *CpuImg {
	raster := make([][4]float32, width*height)

	return &CpuImg{
		Img: Img{
			width:  width,
			height: height,
		},
		Raster: raster,
	}
}

// uploads CPU image to dst.
// if dst is nil, a new GPU buffer
// is allocated
func (ci *CpuImg) Upload(dst *GpuImg) *GpuImg {
	if dst == nil {
		dst = NewGpuImg(ci.Width(), ci.Height())
	} else if ci.Width() != dst.Width() || ci.Height() != dst.Height() {
		panic("shapes do not match")
	}

	dst.buf.Upload(&ci.Raster)

	return dst
}

// wrapper for GPU access of image
type GpuImg struct {
	Img

	buf *glc.Buffer
}

func NewGpuImg(width, height int) *GpuImg {
	return &GpuImg{
		Img: Img{
			width:  width,
			height: height,
		},
		buf: glc.NewBufferStorage(width*height, 4*4, glc.BUF_DYNAMIC_COPY),
	}
}

// free GPU memory
func (gi *GpuImg) Delete() {
	gi.buf.Delete()
}

// maps image to CPU virtual memory,
// pass the slice to fn,
// then unmap
func (gi *GpuImg) Access(fn func(CpuImg)) {
	var raster [][4]float32
	gi.buf.Map(&raster, glc.MAP_READ|glc.MAP_WRITE)
	defer gi.buf.Unmap()

	fn(CpuImg{
		Img: Img{
			gi.width,
			gi.height,
		},
		Raster: raster,
	})
}

// downloads GPU image to CPU.
// if dst is nil, a new CPU
// image is allocated
func (gi *GpuImg) Download(dst *CpuImg) *CpuImg {
	if dst == nil {
		dst = NewCpuImg(gi.Width(), gi.Height())
	} else if gi.Width() != dst.Width() || gi.Height() != dst.Height() {
		panic("shapes do not match")
	}

	gi.Access(func(cSrc CpuImg) {
		copy(dst.Raster, cSrc.Raster)
	})

	return dst
}

func (gi *GpuImg) Bind(slot uint32) {
	gi.buf.Bind(slot)
}
