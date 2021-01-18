package img

import (
	"gocv.io/x/gocv"
)

func (ci CpuImg) ToMat() gocv.Mat {
	var mat gocv.Mat

	mat = gocv.NewMatWithSizes([]int{ci.Height(), ci.Width()}, gocv.MatTypeCV32FC4)
	for x := 0; x < ci.Width(); x++ {
		for y := 0; y < ci.Height(); y++ {
			ptr, _ := mat.DataPtrFloat32()
			ptr[4*y*ci.Width()+4*x+0] = ci.Raster[y*ci.Width()+x][2]
			ptr[4*y*ci.Width()+4*x+1] = ci.Raster[y*ci.Width()+x][1]
			ptr[4*y*ci.Width()+4*x+2] = ci.Raster[y*ci.Width()+x][0]
			ptr[4*y*ci.Width()+4*x+3] = ci.Raster[y*ci.Width()+x][3]
		}
	}

	return mat
}

func (ci CpuImg) Show(winTitle string) {
	win := gocv.NewWindow("glc/img")
	defer win.Close()

	mat := ci.ToMat()
	defer mat.Close()

	win.IMShow(mat)
	gocv.WaitKey(0)
}
