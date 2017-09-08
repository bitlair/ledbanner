package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func main() {
	x := flag.Int("x", 150, "The width of the matrix")
	y := flag.Int("y", 16, "The height of the matrix")
	mag := flag.Int("magification", 12, "Amount of pixels per cell")
	pixelSize := flag.Float64("pixelsize", 0.3, "The space per cell that will be lit")
	flag.Parse()

	banner := Banner{
		lenX:          *x,
		lenY:          *y,
		magnification: *mag,
		pixelSize:     *pixelSize,
	}
	go banner.RunInput(os.Stdin)
	banner.RunDisplay()
}

type Banner struct {
	lenX          int
	lenY          int
	magnification int
	pixelSize     float64
	buffer        []uint8
	bufferStream  chan []uint8
}

func (banner *Banner) RunDisplay() error {
	banner.buffer = make([]uint8, banner.NumPixels()*3)
	banner.bufferStream = make(chan []uint8, 1)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := glfw.Init(); err != nil {
		return fmt.Errorf("Can't init GLFW: %v", err)
	}
	win, err := glfw.CreateWindow(banner.lenX*banner.magnification, banner.lenY*banner.magnification, "LED-Banner", nil, nil)
	if err != nil {
		return err
	}
	defer win.Destroy()
	win.MakeContextCurrent()
	glfw.SwapInterval(1)
	if err := gl.Init(); err != nil {
		return err
	}

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		w, h := win.GetFramebufferSize()
		cellSizeX, cellSizeY := float64(w)/float64(banner.lenX), float64(h)/float64(banner.lenY)
		gl.Viewport(0, 0, int32(w), int32(h))
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(-cellSizeX/2, float64(w)-cellSizeX/2, float64(h)-cellSizeY/2, -cellSizeY/2, -1, 1)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()

		select {
		case banner.buffer = <-banner.bufferStream:
		default:
		}

		pixelSize := banner.pixelSize * cellSizeX
		for x := 0; x < banner.lenX; x++ {
			for y := 0; y < banner.lenY; y++ {
				i := (y*banner.lenX + x) * 3
				gl.Color3ub(banner.buffer[i], banner.buffer[i+1], banner.buffer[i+2])

				rx := float64(x) * cellSizeX
				ry := float64(y) * cellSizeY
				gl.Begin(gl.QUADS)
				gl.Vertex2d(rx-pixelSize/2, ry-pixelSize/2)
				gl.Vertex2d(rx+pixelSize/2, ry-pixelSize/2)
				gl.Vertex2d(rx+pixelSize/2, ry+pixelSize/2)
				gl.Vertex2d(rx-pixelSize/2, ry+pixelSize/2)
				gl.End()
			}
		}

		win.SwapBuffers()
		glfw.PollEvents()
	}

	return nil
}

func (banner *Banner) NumPixels() int {
	return banner.lenX * banner.lenY
}

func (banner *Banner) RunInput(input io.Reader) {
	for {
		buf := make([]uint8, banner.NumPixels()*3)
		io.ReadFull(input, buf)
		banner.bufferStream <- buf
	}
}
