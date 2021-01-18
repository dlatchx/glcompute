package img

import (
	"image"
	//"image/color"
	"io"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func DecodeImg(r io.Reader) (*CpuImg, error) {
	i, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	bounds := i.Bounds()
	width := bounds.Max.X - bounds.Min.X + 1
	height := bounds.Max.Y - bounds.Min.Y + 1

	img := NewCpuImg(width, height)

	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
			ix := (y-bounds.Min.Y)*width + (x - bounds.Min.X)

			r, g, b, a := i.At(x, y).RGBA()
			img.Raster[ix][0] = float32(r) / 65535
			img.Raster[ix][1] = float32(g) / 65535
			img.Raster[ix][2] = float32(b) / 65535
			img.Raster[ix][3] = float32(a) / 65535
		}
	}

	return img, nil
}

func LoadImg(filename string) *CpuImg {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	im, err := DecodeImg(f)
	if err != nil {
		panic(err)
	}

	return im
}
