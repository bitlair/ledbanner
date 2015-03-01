/*
 * Copyright (c) 2014 PolyFloyd
 */

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"runtime"
	gl   "github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

const (
	INFO               = "BitBanner Simulator 2016 Xtreme Edition(tm) v0.1"
	NET_TYPE_DATA byte = 0x01
	NET_TYPE_SWAP byte = 0x02
)

func main() {
	host := flag.String("host",      "0.0.0.0",     "The UDP bind address")
	port := flag.Int("port",         8230,          "The UDP port")
	x    := flag.Int("x",            150,           "The width of the matrix")
	y    := flag.Int("y",            16,            "The height of the matrix")
	mag  := flag.Int("magification", 12,            "Amount of pixels per dot")
	flag.Parse()

	banner := Banner{
		lenX: *x,
		lenY: *y,
		magnification: *mag,
	}
	go banner.RunServer(&net.UDPAddr{
		Port: *port,
		IP:   net.ParseIP(*host),
	})
	banner.RunDisplay()
}


type Banner struct {
	lenX          int
	lenY          int
	magnification int
	buffer        []float32
	bufferStream  chan []float32
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


func (banner *Banner) RunServer(addr *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		buf        := make([]byte,    banner.NumPixels() * 3 + 1+4+2)
		backBuffer := make([]float32, banner.NumPixels() * 3)

		for {
			read, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println(err)
				break
			}

			switch buf[0] {
			case NET_TYPE_SWAP:
				banner.bufferStream <- backBuffer

			case NET_TYPE_DATA:
				if read < 1 + 4 + 2 {
					fmt.Printf("%v error: missing meta information\n", addr)
					continue
				}

				order := binary.LittleEndian
				start  := int(order.Uint32(buf[1:5]))
				length := int(order.Uint16(buf[5:7]))

				if start > len(backBuffer) {
					fmt.Printf("%v error: start index out of range: %v\n", addr, start)
					continue
				}
				if start + length > len(backBuffer) || length == 0 {
					fmt.Printf("%v error: length out of range: %v\n", addr, length)
					continue
				}

				data := buf[7:7+length-1]
				for i, b := range data {
					backBuffer[start+i] = float32(b) / 256
				}
			}
		}
	}
}
