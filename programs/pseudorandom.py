#!/usr/bin/env python3

import math
import random
import sys
import time

DISP_WIDTH  = 150
DISP_HEIGHT = 16


steps         = [math.sin(i / 40 * math.pi / 2) for i in range(0, 40)]
frame_current = bytearray(DISP_WIDTH * DISP_HEIGHT * 3)
frame_target  = bytearray(DISP_WIDTH * DISP_HEIGHT * 3)

while 1:
	frame_source = frame_target
	frame_target = bytearray(DISP_WIDTH * DISP_HEIGHT * 3)
	for i in range(0, len(frame_target)):
		frame_target[i] = random.randint(0, 255)
	for m in steps:
		for j in range(0, len(frame_target)):
			frame_current[j] = int(frame_source[j] * (1 - m) + frame_target[j] * m)
		sys.stdout.buffer.write(frame_current)
		time.sleep(1 / 60)
