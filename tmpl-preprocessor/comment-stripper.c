#include <stdio.h>
#include <assert.h>

#include "startswith.c"

int main() {
	int i;
	char c;
	char inComment = 0;
	char* beginChars = "<!--";
	int beginPos = 0;
	char* endChars = "-->";
	int endPos = 0;
	while (c != EOF) {
		c = getchar();
		if (inComment) {
			if (c == endChars[endPos]) {
				endPos++;
				if (endPos >= strlen(endChars)) {
					inComment = 0;
					endPos = 0;
				}
				continue;
			} else {
				endPos = 0;
			}
		} else {
			if (c == beginChars[beginPos]) {
				beginPos++;
				if (beginPos >= strlen(beginChars)) {
					inComment = 1;
					beginPos = 0;
				}
				continue;
			} else {
				for (i = 0; i < beginPos; i++) {
					putchar(beginChars[i]);
				}
				beginPos = 0;
			}
		}
		if (c != EOF) {
			printf("%c", c);
		}
	}
}