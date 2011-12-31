#include <stdio.h>
#include <assert.h>

#include "startswith.c"

int tagContent(char* buf, int n) {
	int pos = 0;
	char c;
	for (pos = 0; pos < n; pos++) {
		c = getchar();
		if (c == '>') {
			buf[pos] = '\0';
			return pos;
		}
		buf[pos] = c;
	}
}

char isWhitespace(char c) {
	if (c == ' ') return 1;
	if (c == '\t') return 1;
	if (c == '\n') return 1;
	return 0;
}

int main() {
	//fprintf(stderr, "Removing whitespace!\n");
	char rmwhitespace = 1;
	char c;
	char pc = ' ';
	for (c = getchar(); c != EOF; c = getchar()) {
		// Special cases for <pre> tags
		if (c == '<') {
			int TAG_MAX_LENGTH = 2048;
			char tagName[TAG_MAX_LENGTH];
			int tagLen = tagContent(&tagName, TAG_MAX_LENGTH);
			assert(tagLen < TAG_MAX_LENGTH);
			if (startsWith(tagName, TAG_MAX_LENGTH, "pre", 3)) {
				rmwhitespace = 0;
			}
			if (startsWith(tagName, TAG_MAX_LENGTH, "/pre", 4)) {
				rmwhitespace = 1;
			}
			putchar(c);
			fputs(tagName, stdout);
			putchar('>');
			continue;
		}
		if (rmwhitespace) {
			if (isWhitespace(pc) && isWhitespace(c)) {
			} else if (pc == '>' && isWhitespace(c)) {
			} else {
				putchar(c);
			}
		} else {
			putchar(c);
		}
		pc = c;
	}
	return 0;
}