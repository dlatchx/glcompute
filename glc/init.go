package glc

import (
	"fmt"
	"runtime"

	"github.com/remogatto/egl"
    gl "github.com/go-gl/gl/v3.1/gles2"
)

var (
	disp egl.Display
)

// fetch error message form C
// errno system
func eglError() error {
	errno := egl.GetError()
	if errno == egl.SUCCESS {
		return nil
	}

	return fmt.Errorf("%s", egl.NewError(errno).Error())
}

// creates a surfaceless OpenGL context
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

// Init OpenGL context
// must be called before calling any other function from glc
// this function calls glInit()
func Init() {
	runtime.LockOSThread()

	initEGL()

	err := gl.Init()
	if err != nil {
		panic("glompute: could not initialize OpenGL function bindings")
	}
}
