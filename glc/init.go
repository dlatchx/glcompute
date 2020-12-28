package glc

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/remogatto/egl"

	"runtime"
	"fmt"
)

var (
	disp egl.Display
)

func eglError() error {
	errno := egl.GetError()
	if errno == egl.SUCCESS {
		return nil
	}

	return fmt.Errorf("%s", egl.NewError(errno).Error())
}

func initEGL() {
	disp = egl.GetDisplay(egl.DEFAULT_DISPLAY)

	var major int32
	var minor int32
	if !egl.Initialize(disp, &major, &minor) {
		panic("glcompute: could not initialize EGL: " + eglError().Error())
	}

	configAttribs := []int32{
		egl.SURFACE_TYPE, egl.PBUFFER_BIT,
		egl.BLUE_SIZE, 8,
		egl.GREEN_SIZE, 8,
		egl.RED_SIZE, 8,
		egl.DEPTH_SIZE, 8,
		egl.RENDERABLE_TYPE, egl.OPENGL_BIT,
		egl.NONE,
	}
	var config egl.Config
	var numConfigs int32
	if !egl.ChooseConfig(disp, configAttribs, &config, 1, &numConfigs) {
		panic("glcompute: could not set EGL config: " + eglError().Error())
	}

	if !egl.BindAPI(egl.OPENGL_API) {
		panic("glcompute: could not bind API: " + eglError().Error())
	}

	ctx := egl.CreateContext(disp, config, egl.NO_CONTEXT, nil)

	if !egl.MakeCurrent(disp, egl.NO_SURFACE, egl.NO_SURFACE, ctx) {
		panic("glcompute: could not make OpenGL context the current context: " + eglError().Error())
	}
}

func Init() error {
	runtime.LockOSThread()

	initEGL()

	err := gl.Init()
	return err
}
