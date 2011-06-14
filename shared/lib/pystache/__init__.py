from pystache.template import Template
from pystache.view import View
import string

def execute(context):
	if not context:
		print "Context must be defined!"
		return
	try: tmp = context['Name']
	except KeyError as e:
		print "Context name required to determine what template to use..."
		return
	try: tmp = context['ModuleName']
	except KeyError as e:
		# Tries to guess the ModuleName based on the Name of the template.
		names = string.split(context['Name'], "/", 2)
		# Convert the first letter to uppercase.
		modName = string.upper(names[0][:1]) + names[0][1:]

		context['ModuleName'] = {}
		context['ModuleName'][modName] = True
	try: tmp = context['Title']
	except KeyError as e:
		print "Warning: no title will be set for page"
		
	# Actual rendering part.
	# TODO: Add mobile ability.
	view = View(context=context)
	view.template_file = 'tmpl/desktop/'+context['Name']
	return view.render()

def render(template, context=None, **kwargs):
	context = context and context.copy() or {}
	context.update(kwargs)
	return Template(template, context).render()
	
print("Correct pystache imported")