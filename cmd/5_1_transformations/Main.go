package main

import (
	_ "image/jpeg"
	_ "image/png"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	. "learn_opengl_go/pkg"
)

const binName string = "5_1_transformations"
const prefix = "./cmd/" + binName + "/"

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func processInput(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape {
		w.SetShouldClose(true)
	}
}

func framebufferSizeCallback(w *glfw.Window, width int, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func main() {
	window := MakeWindow(800, 600, binName, processInput, framebufferSizeCallback)

	shaderProgram := MakeShader(prefix+"shaders/shader.vert", prefix+"shaders/shader.frag")
	defer shaderProgram.Delete()

	vertices := [...]float32{
		// positions      // colors        // texture coord
		0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, // top left
		0.5, 0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 1.0, // top right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // bottom right
	}

	indices := [...]uint32{
		2, 0, 1, 2, 1, 3,
	}

	var VBO, VAO, EBO uint32
	gl.GenVertexArrays(1, &VAO)
	defer gl.DeleteVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	defer gl.DeleteBuffers(1, &VBO)
	gl.GenBuffers(1, &EBO)
	defer gl.DeleteBuffers(1, &EBO)
	gl.BindVertexArray(VAO)

	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*int(unsafe.Sizeof(vertices[0])), unsafe.Pointer(&vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(unsafe.Sizeof(indices[0])), unsafe.Pointer(&indices), gl.STATIC_DRAW)

	// Position Attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*int32(unsafe.Sizeof(vertices[0])), nil)
	gl.EnableVertexAttribArray(0)

	// Color Attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*int32(unsafe.Sizeof(vertices[1])), unsafe.Pointer(3*unsafe.Sizeof(vertices[0])))
	gl.EnableVertexAttribArray(1)

	// Texture Attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*int32(unsafe.Sizeof(vertices[1])), unsafe.Pointer(6*unsafe.Sizeof(vertices[0])))
	gl.EnableVertexAttribArray(2)

	var containerTexture, emojiTexture uint32
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRRORED_REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	gl.GenTextures(1, &containerTexture)
	gl.BindTexture(gl.TEXTURE_2D, containerTexture)
	{
		containerImage, width, height, err := OpenImageRGBA(prefix + "resources/container.jpg")
		if err != nil {
			panic(err)
		}

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGB,
			int32(width),
			int32(height),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			unsafe.Pointer(&containerImage.Pix[0]),
		)
		gl.GenerateMipmap(gl.TEXTURE_2D)
	}

	gl.GenTextures(1, &emojiTexture)
	gl.BindTexture(gl.TEXTURE_2D, emojiTexture)
	{
		emojiImage, width, height, err := OpenImageRGBA(prefix + "resources/awesomeface.png")
		if err != nil {
			panic(err)
		}

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(width),
			int32(height),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			unsafe.Pointer(&emojiImage.Pix[0]),
		)
		gl.GenerateMipmap(gl.TEXTURE_2D)
	}

	shaderProgram.Use()
	shaderProgram.SetInt("container_texture", 0)
	shaderProgram.SetInt("emoji_texture", 1)

	transformLoc := shaderProgram.GetUniform("transform")

	var timeValue float64
	for !window.Window.ShouldClose() {
		timeValue = glfw.GetTime()

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, containerTexture)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, emojiTexture)

		translation := mgl32.Translate3D(0.3, -0.3, 0.0)
		rotation := mgl32.HomogRotate3D(float32(timeValue), mgl32.Vec3{0.0, 0.0, 1.0})
		transformations := translation.Mul4(rotation)

		shaderProgram.Use()
		gl.UniformMatrix4fv(transformLoc, 1, false, &transformations[0])

		gl.BindVertexArray(VAO)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		// Do OpenGL stuff.
		window.Window.SwapBuffers()
		glfw.PollEvents()
	}
}
