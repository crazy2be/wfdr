all: tmpl-commentstrip tmpl-whitespacerm

tmpl-commentstrip: comment-stripper.c
	gcc -o tmpl-commentstrip comment-stripper.c

tmpl-whitespacerm: whitespace-remover.c
	gcc -o tmpl-whitespacerm whitespace-remover.c

install: all
	mkdir -p ../../bin
	cp tmpl-commentstrip ../../bin/
	cp tmpl-whitespacerm ../../bin

clean:
	rm tmpl-commentstrip
	rm tmpl-whitespacerm