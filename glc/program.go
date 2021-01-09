package glc

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type Program struct {
	id uint32
}

func (p *Program) Ready() bool {
	return p.id != 0
}

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

func NewProgram() *Program {
	programId := gl.CreateProgram()

	program := &Program{id: programId}

	runtime.SetFinalizer(program, func(sp *Program) {
		sp.Delete()
	})

	return program
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

func (p *Program) LoadBuf(bytecode []byte) error {
	gl.ProgramBinary(p.id, gl.SPIR_V_BINARY_ARB, gl.Ptr(bytecode), int32(len(bytecode)))
	return p.check()
}

func (p *Program) LoadFile(path string) error {
	bytecode, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return p.LoadBuf(bytecode)
}

func (p *Program) LoadSrc(path string) error {
	s := NewShader()
	defer s.Delete()

	err := s.CompileFile(path)
	if err != nil {
		return err
	}

	return p.LinkShader(s)
}

func (p Program) Dispatch(x, y, z uint32) {
	p.use()
	gl.DispatchCompute(x, y, z)
}
