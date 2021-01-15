package glc

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
)

// wrapper for shaders
// not necessary since Program
// handles the creation/deletion
// of shaders
type Shader struct {
	id uint32
}

func (s *Shader) Ready() bool {
	return s.id != 0
}

func (s *Shader) Delete() {
	if s.Ready() {
		gl.DeleteShader(s.id)
		s.id = 0
	}
}

func NewShader() *Shader {
	shaderId := gl.CreateShader(gl.COMPUTE_SHADER)

	shader := &Shader{id: shaderId}

	runtime.SetFinalizer(shader, func(s *Shader) {
		s.Delete()
	})

	return shader
}

func (s *Shader) check() error {
	var status int32
	gl.GetShaderiv(s.id, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(s.id, gl.INFO_LOG_LENGTH, &logLength)

		errLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(s.id, logLength, nil, gl.Str(errLog))

		return fmt.Errorf("%s", errLog)
	}

	return nil
}

func (s *Shader) CompileStr(source string) error {
	csource, csourceFree := gl.Strs(source)
	defer csourceFree()
	ln := int32(len(source))
	gl.ShaderSource(s.id, 1, csource, &ln)
	gl.CompileShader(s.id)

	err := s.check()
	if err != nil {
		return fmt.Errorf("failed to compile shader %v: %v", source, err.Error())
	}

	return nil
}

func (s *Shader) CompileFile(path string) error {
	shaderSource, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return s.CompileStr(string(shaderSource))
}
