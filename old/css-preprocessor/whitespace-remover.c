#include <stdio.h>
#include <stdlib.h>

/*
 * CSS Whitespace remover
 * Created by Systemtwo
 * March 2011
 *
 *
 */

int main(int argc, char * argv[]) {
	if (argc < 2) {
		int count;
		while (1) {
			char x = getchar();
			if (x == 13 || x == 10) { //Check for newlines
                count++;
			} else if (x == 32) {
				char y = getchar();
				if (y == 32) { //Check for a second space (will remove two or more spaces in a row)
                    count++;
										count++;
										while (1) { //Continue geting spaces, if not, then print and break
                        char z = getchar();
												if (z == 32) {
													count++;
												} else {
													printf("%c", z);
													break;
												}
										}
				} else {
					printf(" %c", y);
				}
			} else if (x == 9) { //Check for tabs
                count++;
			} else if (x == EOF) {
				break;
			} else {
				printf("%c", x);
			}
		}
	} else {
		puts("CSS Whitespace remover");
		puts("======================");
		puts("By Systemtwo");
		puts("This program will remove all tabs, newlines, linefeeds as");
		puts("well as two or more spaces in a row");
		puts("");
		puts("Usage:");
		puts("csswtsprem");
		puts("");
		puts("Feed in CSS thru stdin, will output thru stdout");
	}
}