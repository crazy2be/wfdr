all: commentstrip varfinder whitespacerm

commentstrip: comment-stripper.c
	gcc -o commentstrip comment-stripper.c

varfinder: variable-finder.c
	gcc -o varfinder variable-finder.c

whitespacerm: whitespace-remover.c
	gcc -o whitespacerm whitespace-remover.c

install: all
	mkdir -p ../../bin
	cp commentstrip ../../bin/
	cp varfinder ../../bin/
	cp whitespacerm ../../bin

clean:
	rm commentstrip
	rm varfinder
	rm whitespacerm