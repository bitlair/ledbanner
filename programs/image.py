#! /usr/bin/env python3

#
# Copyright (c) 2015 Bitlair
#

# Requires the Pillow library

import ledbanner
import sys
import time
from PIL import Image

if len(sys.argv) == 1:
	print('Usage: %s <image 1> [<image 2> <image 3> ...]' % sys.argv[0])
	exit(1)

def read_image(filename):
	img = Image.open(filename)
	img.thumbnail((led.size.x, float('inf')))
	return img

led    = ledbanner.LEDBanner()
images = []

for f in sys.argv[1:]:
	try:
		image = read_image(f)
		image_size = ledbanner.Vector(image.size[0], image.size[1])
		raw = bytearray(image_size.x * image_size.y * led.bytes_per_pixel)
		for x in range(0, image_size.x):
			for y in range(0, image_size.y):
				pix = image.getpixel((x, y))
				if type(pix) is int: # Monochrome
					pix = (pix, pix, pix)
				i = (y * led.size.x + x) * led.bytes_per_pixel
				raw[i:i+3] = pix
		images.append((raw, image_size))
	except:
		pass
if len(images) == 0:
	print('Unable to read any of the specified images')
	exit(1)

while 1:
	for image, image_size in images:
		for scroll in range(led.size.y, -image_size.y, -1):
			frame = led.make_frame()

			frame_start_y = scroll
			if frame_start_y < 0:
				frame_start_y = 0
			frame_stop_y  = scroll + image_size.y
			if frame_stop_y > led.size.y:
				frame_stop_y = led.size.y
			image_start_y = -scroll
			if image_start_y < 0:
				image_start_y = 0
			image_stop_y = image_start_y + frame_stop_y - frame_start_y
			if image_stop_y > image_size.y:
				image_stop_y = image_size.y

			bytes_per_y = led.size.x * led.bytes_per_pixel
			frame[frame_start_y * bytes_per_y:frame_stop_y * bytes_per_y] = \
				image[image_start_y*bytes_per_y:image_stop_y*bytes_per_y]

			led.set_frame(frame)
			time.sleep(1 / led.fps)
