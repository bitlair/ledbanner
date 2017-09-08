#!/usr/bin/env python3

from collections import namedtuple
import sys
import time
from PIL import Image # Requires the Pillow library

DISP_WIDTH  = 150
DISP_HEIGHT = 16


Vector = namedtuple('Vector', 'x y')

if len(sys.argv) == 1:
	print('Usage: %s <image 1> [<image 2> <image 3> ...]' % sys.argv[0], file=sys.stderr)
	exit(1)

def read_image(filename):
	img = Image.open(filename)
	img.thumbnail((DISP_WIDTH, float('inf')))
	return img

images = []
for f in sys.argv[1:]:
	try:
		image = read_image(f)
		image_size = Vector(image.size[0], image.size[1])
		raw = bytearray(image_size.x * image_size.y * 3)
		for x in range(0, image_size.x):
			for y in range(0, image_size.y):
				pix = image.getpixel((x, y))
				if type(pix) is int: # Monochrome
					pix = (pix, pix, pix)
				elif len(pix) == 4: # Alhpa
					pix = pix[:3]
				i = (y * DISP_WIDTH + x) * 3
				raw[i:i+3] = pix
		images.append((raw, image_size))
	except Exception as ex:
		print('Unable to read %s: %s' % (f, ex), file=sys.stderr)
if len(images) == 0:
	print('Unable to read any of the specified images', file=sys.stderr)
	exit(1)

while 1:
	for image, image_size in images:
		for scroll in range(DISP_HEIGHT, -image_size.y, -1):
			frame = bytearray(DISP_WIDTH * DISP_HEIGHT * 3)

			frame_start_y = scroll
			if frame_start_y < 0:
				frame_start_y = 0
			frame_stop_y  = scroll + image_size.y
			if frame_stop_y > DISP_HEIGHT:
				frame_stop_y = DISP_HEIGHT
			image_start_y = -scroll
			if image_start_y < 0:
				image_start_y = 0
			image_stop_y = image_start_y + frame_stop_y - frame_start_y
			if image_stop_y > image_size.y:
				image_stop_y = image_size.y

			bytes_per_y = DISP_WIDTH * 3
			frame[frame_start_y * bytes_per_y:frame_stop_y * bytes_per_y] = \
				image[image_start_y*bytes_per_y:image_stop_y*bytes_per_y]

			sys.stdout.buffer.write(frame)
			time.sleep(1 / 60 * 2)
