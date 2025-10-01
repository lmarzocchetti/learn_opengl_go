package pkg

// #include <stdlib.h>
import "C"

import (
	"log"
	"os"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Shader uint32

func MakeShader(vertexPath string, fragmentPath string) Shader {
	vertexCode, err := os.ReadFile(vertexPath)
	if err != nil {
		panic(err)
	}
	vertexCodeCStr := C.CString(string(vertexCode))
	defer C.free(unsafe.Pointer(vertexCodeCStr))

	fragmentCode, err := os.ReadFile(fragmentPath)
	if err != nil {
		panic(err)
	}
	fragmentCodeCStr := C.CString(string(fragmentCode))
	defer C.free(unsafe.Pointer(fragmentCodeCStr))

	success := int32(0)
	infoLog := [512]byte{}

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	defer gl.DeleteShader(vertexShader)
	gl.ShaderSource(vertexShader, 1, (**uint8)(unsafe.Pointer(&vertexCodeCStr)), nil)
	//gl.ShaderSource(vertexShader, 1, (**uint8)(unsafe.Pointer(&vertexCode[0])), nil)
	gl.CompileShader(vertexShader)
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		gl.GetShaderInfoLog(vertexShader, 512, nil, &infoLog[0])
		log.Fatalf("failed to compile vertex: %v", string(infoLog[:]))
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(fragmentShader)
	gl.ShaderSource(fragmentShader, 1, (**uint8)(unsafe.Pointer(&fragmentCodeCStr)), nil)
	//gl.ShaderSource(fragmentShader, 1, (**uint8)(unsafe.Pointer(&fragmentCode[0])), nil)
	gl.CompileShader(fragmentShader)
	gl.GetShaderiv(fragmentShader, gl.COMPILE_STATUS, &success)
	if success == gl.FALSE {
		gl.GetShaderInfoLog(fragmentShader, 512, nil, &infoLog[0])
		log.Fatalf("failed to compile fragment: %v", string(infoLog[:]))
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		gl.GetProgramInfoLog(shaderProgram, 512, nil, &infoLog[0])
		log.Fatalf("failed to link shader: %v", string(infoLog[:]))
	}

	return Shader(shaderProgram)
}

func (s Shader) Delete() {
	gl.DeleteProgram(uint32(s))
}

func (s Shader) Use() {
	gl.UseProgram(uint32(s))
}

func (s Shader) GetUniform(name string) int32 {
	nameCStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameCStr))
	return gl.GetUniformLocation(uint32(s), (*uint8)(unsafe.Pointer(nameCStr)))
}

func (s Shader) SetBool(name string, value bool) {
	var valueInt int32
	if value {
		valueInt = 1
	} else {
		valueInt = 0
	}
	gl.Uniform1i(s.GetUniform(name), valueInt)
}

func (s Shader) SetInt(name string, value int32) {
	gl.Uniform1i(s.GetUniform(name), value)
}

func (s Shader) SetFloat(name string, value float32) {
	gl.Uniform1f(s.GetUniform(name), value)
}
