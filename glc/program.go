package glc

import (
	"fmt"
	"runtime"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"
)

// Hight-level wrapper for shader programs
type Program struct {
	id uint32
}

// returns true if the buiffer is ready to be used
func (p *Program) Ready() bool {
	return p.id != 0
}

// free GPU memory
func (p *Program) Delete() {
	if p.Ready() {
		gl.DeleteProgram(p.id)
		p.id = 0
	}
}

func (p *Program) getUniformLocation(name string) int32 {
	str := gl.Str(name + "\000")

	return gl.GetUniformLocation(p.id, str)
}

func (p Program) use() {
	gl.UseProgram(p.id)
}

// returns a new empty program
// it must be linked with a shader before
// it can be used
func NewProgram() *Program {
	programId := gl.CreateProgram()

	program := &Program{id: programId}

	runtime.SetFinalizer(program, func(sp *Program) {
		sp.Delete()
	})

	return program
}

// init a program then link it
// with a shader copiled from the source
// file passed as argument
func LoadProgram(path string) *Program {
	prgm := NewProgram()

	err := prgm.LoadSrc(path)
	if err != nil {
		panic(err)
	}

	return prgm
}

func (p *Program) check() error {
	var status int32
	gl.GetProgramiv(p.id, gl.LINK_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(p.id, gl.INFO_LOG_LENGTH, &logLength)

		errLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(p.id, logLength, nil, gl.Str(errLog))

		return fmt.Errorf("failed to link program: %v", errLog)
	}

	return nil
}

func (p *Program) LinkShader(shader *Shader) error {
	gl.AttachShader(p.id, shader.id)
	gl.LinkProgram(p.id)
	return p.check()
}

// compiles a shader from a file
// and links it with program
func (p *Program) LoadSrc(path string) error {
	s := NewShader()
	defer s.Delete()

	err := s.CompileFile(path)
	if err != nil {
		return err
	}

	return p.LinkShader(s)
}

// calls glDispatchCompute()
func (p Program) Dispatch(x, y, z uint32) {
	p.use()
	gl.DispatchCompute(x, y, z)
}

// bind the list of buffers passed in ascending order
// then calls glDispatchCompute
// ie. p.Call(a, b, c, buf1, nil, buf2) is equivalent to
// buf1.Bind(0)
// buf2.Bind(2)
// p.Dispacth(a, b, c)
func (p Program) Call(x, y, z uint32, args ...*Buffer) {
	for i, b := range args {
		if b != nil {
			b.Bind(uint32(i))
		}
	}

	p.Dispatch(x, y, z)
}
