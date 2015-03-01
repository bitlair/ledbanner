/*
 * Copyright (c) 2014 PolyFloyd
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"runtime"
	"unsafe"
	gl   "github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

const (
	INFO          = "BitBanner Simulator v0.1"
	MATRIX_X      = 150
	MATRIX_Y      = 16
	MAGNIFICATION = 12
)

func main() {
	l   := flag.String("l",         ":54746",      "The TCP host and port for incoming connections")
	x   := flag.Int("x",            MATRIX_X,      "The width of the matrix")
	y   := flag.Int("y",            MATRIX_Y,      "The height of the matrix")
	mag := flag.Int("magification", MAGNIFICATION, "The level of detail")
	flag.Parse()

	banner := Banner{
		lenX: *x,
		lenY: *y,
		magnification: *mag,
	}

	go banner.RunServer(*l)
	banner.RunDisplay()
}


type Banner struct {
	lenX              int
	lenY              int
	magnification     int
	buffer           []float32
	bufferStream     chan []float32
}

func (banner *Banner) RunDisplay() error {
	banner.buffer = make([]float32, banner.NumPixels() * 3)
	banner.bufferStream = make(chan []float32, 1)

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	runtime.GOMAXPROCS(runtime.NumCPU())

	glfw.SetErrorCallback(func(code glfw.ErrorCode, desc string) {
		fmt.Printf("GLFW Error: %v\n", desc)
	})
	if !glfw.Init() {
		return fmt.Errorf("Can't init GLFW!")
	}
	win, err := glfw.CreateWindow(banner.lenX * banner.magnification, banner.lenY * banner.magnification, INFO, nil, nil)
	if err != nil {
		return err
	}
	defer win.Destroy()
	win.MakeContextCurrent()
	glfw.SwapInterval(1)
	if gl.Init() != gl.FALSE {
		return fmt.Errorf("Could not initialize OpenGL")
	}

	for !win.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		w, h := win.GetFramebufferSize()
		gl.Viewport(0, 0, w, h)
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(
			0,
			float64(banner.lenX),
			float64(banner.lenY),
			0,
			-1,
			1,
		)
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()

		select {
		case banner.buffer = <- banner.bufferStream:
		default:
		}

		for x := 0; x < banner.lenX; x++ {
			for y := 0; y < banner.lenY; y++ {
				gl.Begin(gl.QUADS)
				i := (x * banner.lenY + y) * 3
				gl.Color3f(banner.buffer[i], banner.buffer[i+1], banner.buffer[i+2])
				gl.Vertex2i(x,   y)
				gl.Vertex2i(x+1, y)
				gl.Vertex2i(x+1, y+1)
				gl.Vertex2i(x,   y+1)
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


func (banner *Banner) RunServer(listen string) {
	listener, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go communicate(conn, banner)
	}
}

func communicate(conn net.Conn, banner *Banner) {
	buf := make([]byte, banner.NumPixels() * 3)
	var backBuffer []float32
	main: for {
		_, err := conn.Read(buf[:3])
		if err != nil {
			break
		}
		switch string(buf[:3]) {
		case "ver":
			conn.Write([]byte(INFO))

		case "inf":
			x := *(*[4]byte)(unsafe.Pointer(&banner.lenX))
			conn.Write(x[:])
			y := *(*[4]byte)(unsafe.Pointer(&banner.lenY))
			conn.Write(y[:])
			one := 1
			z := *(*[4]byte)(unsafe.Pointer(&one))
			conn.Write(z[:])
			conn.Write([]byte{ 3 })
			conn.Write([]byte{ 60 })

		case "put":
			backBuffer = make([]float32, banner.NumPixels() * 3)
			for completed := 0; completed < banner.NumPixels()*3; {
				read, err := conn.Read(buf[:banner.NumPixels()*3 - completed])
				if err != nil {
					break main
				}
				for i, b := range buf[:read] {
					backBuffer[completed+i] = float32(b) / 256
				}
				completed += read
			}

		case "swp":
			banner.bufferStream <- backBuffer
		}
	}
}
