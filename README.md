BitBanner (name subject to change) Simulator
============================================

## Building

Make sure you have [Go](http://golang.org/dl), GLFW3 and GLEW installed.
On Debian (and maybe Ubuntu) systems, you can install all dependencies with this command:
```sh
	sudo apt-get install git golang libgflw3-dev libglew-dev
```

Finally, install and build bitbanner:
```sh
	git clone REPO bitbanner
	cd bitbanner
	GOPATH=$PWD/gopath go get; go build
```

## Command Arguments

	-host="0.0.0.0":  The UDP bind address
	-port=8230:       The UDP port
	-magification=12: Amount of pixels per dot
	-x=150:           The width of the matrix
	-y=16:            The height of the matrix

## Network Protocol

Type field enumeration:
* 0x01: Data
* 0x02: Swap

### Packet Format

#### Data
A data packet contains a portion of a frame.

| Start - End Index  | Data           | Description                     |
| ------------------ | -------------- | ------------------------------- |
| 0 - 0              | uint8          | Type                            |
| 1 - 4              | uint32         | Start index int the framebuffer |
| 5 - 6              | uint16         | Length of the data              |
| 7 - 7 + Length - 1 | uint8 * Length | Data                            |

Data is encoded as 24 bit RGB

#### Swap
Upon receiving a swap packet, the displaybuffers are swapped.

| Start - End Index | Data       | Description |
| ----------------- | ---------- | ----------- |
| 0  - 3            | uint8=0x02 | Type        |
