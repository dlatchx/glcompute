package img

import (
	"github.com/adraenwan/glcompute/glc"
)

// image metadata
type Img struct {
	width  int
	height int
	nbChan int
}

func (i Img) Width() int {
	return i.width
}

func (i Img) Height() int {
	return i.height
}

func (i Img) NbChan() int {
	return i.nbChan
}

// wrapper for cpu access of image
type CpuImg struct {
	Img

	Chans [][]float32

	raster []float32
}

func NewCpuImg(width, height, nbChan int) *CpuImg {
	rasterLen := width * height * nbChan

	raster := make([]float32, rasterLen)

	chans := make([][]float32, nbChan)
	for i := 0; i < nbChan; i++ {
		chans[i] = raster[rasterLen*i : rasterLen*(i+1)]
	}

	return &CpuImg{
		Img: Img{
			width:  width,
			height: height,
			nbChan: nbChan,
		},
		Chans:  chans,
		raster: raster,
	}
}

// uploads CPU image to dst.
// if dst is nil, a new GPU buffer
// is allocated
func (ci *CpuImg) Upload(dst *GpuImg) *GpuImg {
	if dst == nil {
		dst = NewGpuImg(ci.Width(), ci.Height(), ci.NbChan())
	} else if ci.Width() != dst.Width() || ci.Height() != dst.Height() || ci.NbChan() != dst.NbChan() {
		panic("shapes do not match")
	}

	dst.buf.Upload(&ci.raster)

	return dst
}

// wrapper for GPU access of image
type GpuImg struct {
	Img

	buf *glc.Buffer
}

func NewGpuImg(width, height, nbChan int) *GpuImg {
	return &GpuImg{
		Img: Img{
			width:  width,
			height: height,
			nbChan: nbChan,
		},
		buf: glc.NewBufferStorage(width*height*nbChan, 4, glc.BUF_DYNAMIC_COPY),
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
	rasterLen := gi.Width() * gi.Height() * gi.NbChan()
	var raster []float32
	gi.buf.Map(&raster, glc.MAP_READ|glc.MAP_WRITE)
	defer gi.buf.Unmap()

	chans := make([][]float32, gi.NbChan())
	for i := 0; i < gi.NbChan(); i++ {
		chans[i] = raster[rasterLen*i : rasterLen*(i+1)]
	}

	fn(CpuImg{
		Img: Img{
			gi.width,
			gi.height,
			gi.nbChan,
		},
		Chans: chans,
	})
}

// downloads GPU image to CPU.
// if dst is nil, a new CPU
// image is allocated
func (gi *GpuImg) Download(dst *CpuImg) *CpuImg {
	if dst == nil {
		dst = NewCpuImg(gi.Width(), gi.Height(), gi.NbChan())
	} else if gi.Width() != dst.Width() || gi.Height() != dst.Height() || gi.NbChan() != dst.NbChan() {
		panic("shapes do not match")
	}

	gi.Access(func(cSrc CpuImg) {
		for i := 0; i < gi.NbChan(); i++ {
			copy(dst.Chans[i], cSrc.Chans[i])
		}
	})

	return dst
}
