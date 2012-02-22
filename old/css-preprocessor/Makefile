all: css-commentstrip css-varfinder css-whitespacerm

css-commentstrip: comment-stripper.c
	gcc -o css-commentstrip comment-stripper.c

css-varfinder: variable-finder.c
	gcc -o css-varfinder variable-finder.c

css-whitespacerm: whitespace-remover.c
	gcc -o css-whitespacerm whitespace-remover.c

install: all
	mkdir -p ../../bin
	cp css-commentstrip ../../bin/
	cp css-varfinder ../../bin/
	cp css-whitespacerm ../../bin

clean:
	rm css-commentstrip
	rm css-varfinder
	rm css-whitespacerm