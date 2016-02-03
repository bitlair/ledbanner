#
# Copyright (c) 2015, Bitlair
#

from collections import namedtuple
import os
import socket
import sys

Vector = namedtuple('Vector', 'x y')

def determine_connection():
	addr = os.getenv("LEDBANNER_ADDR")
	port = os.getenv("LEDBANNER_PORT")

	for (i, arg) in enumerate(sys.argv[1:]):
		if arg == "-a":
			addr = sys.argv[i + 2]
		elif arg == "-p":
			port = int(sys.argv[i + 2])

	if not addr:
		addr = "127.0.0.1"
	if not port:
		port = 8230

	return (addr, port)


class LEDBanner(object):

	def __init__(self, server=determine_connection()):
		self.socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
		self.server = server

		self.size            = Vector(150, 16)
		self.bytes_per_pixel = 3
		self.fps             = 60

	def send(self, packet):
		self.socket.sendto(packet, self.server)

	def set_frame(self, data, swap=True):
		packet = bytearray(1 + 4 + 2 + len(data))

		packet[0]             = 0x01                                      # type
		packet[1:4]           = (0).to_bytes(4,       byteorder="little") # start index
		packet[5:6]           = len(data).to_bytes(2, byteorder="little") # length
		packet[7:7+len(data)] = data

		self.send(packet)
		if swap:
			self.swap()

	def make_frame(self):
		return Frame(self.size, self.bytes_per_pixel)

	def swap(self):
		self.send(bytes([0x02]))


class Frame(bytearray):

	def __init__(self, size, bytes_per_pixel):
		super(Frame, self).__init__(size.x * size.y * bytes_per_pixel)
		self.size            = size
		self.bytes_per_pixel = bytes_per_pixel

	def index(self, x, y):
		x, y = int(x), int(y)
		if 0 <= x < self.size.x and 0 <= y < self.size.y:
			return (y * self.size.x + x) * self.bytes_per_pixel
		return -1

	def get(self, x, y):
		i = self.index(x, y)
		if i == -1:
			raise IndexError("(%s, %s) is outside screenspace" % (x, y))
		return (self[i] / 255, self[i + 1] / 255, self[i + 2] / 255)

	def set(self, x, y, pixel, clip=True):
		i = self.index(x, y)
		if i != -1:
			for j in range(0, self.bytes_per_pixel):
				self[i + j] = int(pixel[j] * 255)
		elif not clip:
			raise IndexError("(%s, %s) is outside screenspace" % (x, y))
