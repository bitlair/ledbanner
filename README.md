BitBanner (name subject to change) Simulator
============================================

## Building

Make sure you have [Go](http://golang.org/dl), GLFW3 and GLEW installed.
On Debian (and maybe Ubuntu) systems, you can install all dependencies with this command:
```sh
sudo apt-get install git golang libgflw3-dev libglew-dev
```
Build:
```sh
git clone https://github.com/bitlair/ledbanner-sim.git bitbanner
cd bitbanner
./just install && ./just build
```

## Usage
Pipe an RGB24 image with a size of 150x16 pixels to stdin.

## Command Arguments
    -magification=12: Amount of pixels per dot
    -x=150:           The width of the matrix
    -y=16:            The height of the matrix
