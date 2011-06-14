#!/usr/bin/env python
from BaseHTTPServer import BaseHTTPRequestHandler,HTTPServer
import sys
try:
	sys.path.insert(1, "lib")
	import pystache
except ImportError as e:
	print "pystache not found in lib/"

global_login_page = {'Name': 'auth/global-login', 'Title': 'Global Authentication', 'Object': {}}

#Subclass means modify exitsting class thingy
class MyServer(BaseHTTPRequestHandler):
	def do_GET(self):
		self.send_response(200, 'OK')
		self.send_header('Content-type', 'text/html')
		self.end_headers()
		self.wfile.write(pystache.execute(global_login_page))
	@staticmethod
	#"Model" server and start server
	def serve_forever(port):
		HTTPServer(('', port), MyServer).serve_forever()

if __name__ == "__main__":
	# Make object called httpd from MyServer
	httpd = MyServer
	httpd.serve_forever(8081)