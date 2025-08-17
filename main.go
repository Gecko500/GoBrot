package main

import (
	"fmt"
	"math"
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/Gecko500/GoBrot/openGL"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	vid_mode := glfw.GetPrimaryMonitor().GetVideoMode()

	window, prog, err := openGL.InitGLFW(vid_mode.Width, vid_mode.Height)
	if err != nil {
		panic(err)
	}
	defer gl.DeleteProgram(prog)
	defer window.Destroy()

	gl.UseProgram(prog)

	// Set the clear color to a light blue
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// Create a fullscreen quad VAO
	fullscreen_vao := openGL.MakeVao([]float32{
		-1.0, -1.0, 0.0,
		1.0, -1.0, 0.0,
		-1.0, 1.0, 0.0,
		1.0, 1.0, 0.0,
	})

	max_iter := 1000                               // Maximum iterations for the Mandelbrot set
	zoom := float64(2.5)                           // Initial zoom level
	offsetX, offsetY := float64(0.0), float64(0.0) // Initial offsets
	cursorX, cursorY := float64(0), float64(0)
	mousePressed := false

	window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {

		offsetX += float64(cursorX) * 0.16 * zoom // Adjust offset based on cursor position
		offsetY += float64(cursorY) * 0.16 * zoom // Adjust offset based on cursor position

		zoom -= float64(yoff) * 0.2 * zoom // Adjust zoom sensitivity

		// Increase iterations as you zoom in for more detail
		max_iter = int(1000.0 + 10.0*math.Pow(1.0/zoom, 0.3))
	})

	xpos_old, ypos_old := 0.0, 0.0

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		width, height := w.GetSize()
		cursorX = xpos/float64(width)*2.0 - 1.0           // Normalize to [-1, 1]
		cursorY = -((ypos / float64(height) * 2.0) - 1.0) // Invert Y axis and normalize

		if mousePressed {
			if int(xpos) != int(xpos_old) || int(ypos) != int(ypos_old) {
				offsetX -= float64(int(xpos)-int(xpos_old)) * 0.001 * zoom
				offsetY += float64(int(ypos)-int(ypos_old)) * 0.001 * zoom
			}
		}
		xpos_old, ypos_old = xpos, ypos // Always update to current float values
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if button == glfw.MouseButtonLeft {
			mousePressed = action == glfw.Press
		}
	})

	width, height := window.GetSize()
	setShaderVarVec2(prog, "u_resolution", float32(width), float32(height))

	draw(window, prog, fullscreen_vao, zoom, offsetX, offsetY, max_iter)

	zoom_old, offsetX_old, offsetY_old := zoom, offsetX, offsetY
	for !window.ShouldClose() {
		// draw if values have changed
		if zoom != zoom_old || offsetX != offsetX_old || offsetY != offsetY_old {
			draw(window, prog, fullscreen_vao, zoom, offsetX, offsetY, max_iter)
			zoom_old, offsetX_old, offsetY_old = zoom, offsetX, offsetY
		}
		glfw.PollEvents()
	}

	fmt.Println("Exiting...")
}

func draw(window *glfw.Window, program uint32, vao uint32, zoom, offsetX, offsetY float64, max_iter int) {
	start := glfw.GetTime()

	setShaderVarFloat(program, "u_zoom", zoom)
	setShaderVarVec2d(program, "u_offset", offsetX, offsetY)
	setShaderVarInt(program, "max_iter", max_iter)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	window.SwapBuffers()

	// Only print render info if frame took longer than 0.1s (for debug)
	renderTime := glfw.GetTime() - start
	if renderTime > 0.1 {
		fmt.Printf("Z: %f X: %f Y: %f N: %d in %fs\n", zoom, offsetX, offsetY, max_iter, renderTime)
	}
}

func setShaderVarFloat(program uint32, name string, value float64) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform1d(loc, value)
}

func setShaderVarVec2d(program uint32, name string, x, y float64) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform2d(loc, x, y)
}

func setShaderVarVec2(program uint32, name string, x, y float32) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform2f(loc, x, y)
}

func setShaderVarInt(program uint32, name string, value int) {
	loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform1i(loc, int32(value))
}
