package openGL

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	BROT_SHADER_FILE = "brotShader.glsl"
)

// MakeVao initializes and returns a vertex array from the points provided.
func MakeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func InitGLFW(width, height int) (window *glfw.Window, prog uint32, err error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	window, err = glfw.CreateWindow(width, height, "GoBrot", nil, nil)
	if err != nil {
		return nil, 0, err
	}
	window.MakeContextCurrent()

	prog, err = initGL(window)
	if err != nil {
		return nil, 0, err
	}

	return window, prog, nil
}

func initGL(window *glfw.Window) (uint32, error) {
	if err := gl.Init(); err != nil {
		return 0, err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	glsl := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
	fmt.Println("OpenGL version", version, glsl)

	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	// compile shaders
	fragmentShader, err := compileShader(BROT_SHADER_FILE, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	// add shaders
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog, err
}

func compileShader(shaderFilePath string, shaderType uint32) (uint32, error) {

	//read the shader source code from file
	source, err := os.ReadFile(shaderFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read shader file %s: %v", shaderFilePath, err)
	}

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(string(source) + "\x00")
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
