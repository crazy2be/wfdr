# ifile must be file object, not path
# sep is symbol to seperate
def lstify(ifile, sep):
	i = []
	while True:
		line = ifile.readline()
		print "Line:", line
		# If the length of the line is 0, the end of the file has been reached. See http://docs.python.org/tutorial/inputoutput.html. An empty line will have a length of one, literally '\n'.
		if len(line) == 0: #This is BAD! It will not read the rest of the file
			print "Warning: *Empty line!"
			break
		else:
			if line[0] == '#':
				# Comment
				continue
		
			t = line.split('=', 2) #should be sep!!!!!!!!!!!!!!!!!!!
		
			# No equals sign in line, either blank or invalid. Skip in any case.
			if len(t) < 2:
				continue
		
			i.append (t[0])
			i.append (t[1].rstrip('\n'))
		
	print i
	return i


# def tuplefy (lst): #Found problem. lst (the argument) does not count the blank spots, making it too low, excluding ones at end.
#	print "Tuplefy says:", len(lst)
#	for i in range (0, len(lst)-1, 2):
#		print (lst[i], lst[i+1])
#	return [(lst[i],lst[i+1]) for i in range(0,len(lst)-1,2)] #len(lst)-2 CHANGED to len(lst)-1. Seems to fix the problem, because the last number would be cut off

def tuplefy (lst):
	return [(lst[i],lst[i+1]) for i in range(0,len(lst)-1,2)]

	
def triplefy (lst): #For short summary in clist.txt
	return [(lst[i],lst[i+1], lst[i+2]) for i in range(0,len(lst)-3,3)]
	

# Escape HTML
#
# Needed for clubss
# Why?
# Mustache escapes html automatically...
# For raw html, use {{{Name}}} rather than {{Name}}
html_escape_table = {
    "&": "&amp;",
    '"': "&quot;",
    "'": "&apos;",
    ">": "&gt;",
    "<": "&lt;",
    }

def html_escape(text):
    """Produce entities within text."""
    return "".join(html_escape_table.get(c,c) for c in text)

