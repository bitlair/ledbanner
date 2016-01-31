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
	img.thumbnail((led.size[0], float('inf')))
	return img

led    = ledbanner.LEDBanner()
images = []

for f in sys.argv[1:]:
	try:
		images.append(read_image(f))
	except:
		pass
if len(images) == 0:
	print('Unable to read any of the specified images')
	exit(1)

while 1:
	for image in images:
		for i in range(-led.size.y, image.size[1]):
			frame = led.make_frame()
			for x in range(0, led.size.x):
				for y in range(0, led.size.y):
					scroll = y + i
					if scroll >= 0 and scroll < image.size[1]:
						pix = image.getpixel((x, scroll))
						if type(pix) is int: # Monochrome
							pix = [pix, pix, pix]
						frame.set(x, y, (pix[0] / 255, pix[1] / 255, pix[2] / 255))

			led.set_frame(frame)
			time.sleep(1 / led.fps)
