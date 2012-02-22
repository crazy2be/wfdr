#include <stdio.h>
#include <stdlib.h>

int main(int argc, char *argv) {
    if (argc < 2) {
        while (1) {
            char x = getchar();
            if (x == 47) { //if /
                char y = getchar();
                if (y == 47) { // If another / then remove till newline
                    while (1) {
                        char z = getchar();
                        if (z == 10 || z == 13) { // Found end of line; break
                            break;
                        }
                    }
                } else if (y == 42) {
                    //printf ("Into Long comment loop");
                    while (1) {
                        char yy = getchar();
                        if (yy == 42) {
                            char zz = getchar();
                            if (zz == 47) {
                                break;
                            }
                        }
                    }
                } else {
                    printf ("/%c", y); // Replace y, since it is not a comment
                }
            } else if (x == NULL) {
                break;
            } else {
                printf ("%c", x); //Print the char, not a /
            }
        }
    } else {
        printf ("CSS De-Commenter\n");
        printf ("================\n");
        printf ("By: Systemtwo\n\n");
        printf ("Usage: Input thru stdin, it will give uncommented code \n(C-style or CSS-style [// or /*foobar*/]) \nthru stdout.\n\n");
        printf ("File Based Input not implemented in this version\n");
    }
}


