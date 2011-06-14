#! /usr/bin/env python
# Read Only, no write yet!
# Write may not be a good idea... yet
#
#
# Terminology and notes for this module
# 
# If you haven't noticed, some of the stuff
# on this document reflect past online memes
#
# The plate variable holds cookies. (What else
# do you do with plates?)
#
# Cookies will be stored as cookievalue


# This returns cookies from the header string
def findcookie(string):
	print "Session says: I have a cookie!"
	print "Its very long!"
	print "Ready?"
	print "Here we go!"
	print string
	if string.has_key('Cookie'):
		print "Has a cookie!"
	else:
		print "I can't haz cookie?"
	
	plate = string['Cookie']
		
def readcookie(cookie):
	pass


if __name__ == "session":
	print "Session Initalized!"
