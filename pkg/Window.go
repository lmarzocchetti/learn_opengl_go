package pkg

import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	Window *glfw.Window
}

func MakeWindow(width int,
	height int,
	title string,
	processInput glfw.KeyCallback,
	framebufferSizeCallback glfw.FramebufferSizeCallback) Window {
	window := Window{}
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)

	windowGlfw, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		glfw.Terminate()
		panic(err)
	}
	window.Window = windowGlfw

	windowGlfw.MakeContextCurrent()
	if gl.Init() != nil {
		glfw.Terminate()
		panic(err)
	}

	glfw.SwapInterval(1)
	windowGlfw.SetKeyCallback(processInput)
	windowGlfw.SetFramebufferSizeCallback(framebufferSizeCallback)

	return window
}

func (window Window) Destroy() {
	window.Window.Destroy()
	glfw.Terminate()
}
