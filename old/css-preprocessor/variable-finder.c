#include <stdio.h>
#include <stdlib.h>
#include <string.h>

char * readglobal(char * fstr, char * gfile);


int readcssstd(char * gfile) { // for read from stdin
//gfile is global var file
int i;
char var[100];
while (1) {
	int x = getchar();
	if (x == -1) { //check for end of file (null)
            break;
	} else if (x == 36) { //Check for $ (variable)
            //printf ("FOUND VAR!");
            var[0] = 36;
						for (i = 1; i <= 100; i++) {
							var[i] = (char)getchar(); // Load var name into var array
							if (var[i] == 59) { //stop on ; and print variable then print ;
                    var[i] = NULL;
										char *ffstr = readglobal(var, gfile);  //pass pointer (array)
										//printf ("%d%c", strlen(ffstr), ffstr[strlen(ffstr)-1]);
										if (ffstr[strlen(ffstr)-1] == 13 || ffstr[strlen(ffstr)-1] == 10) {
											ffstr[strlen(ffstr)-1] = NULL; // Chomp off the last \n if there is any (Set char to null)
										}
										printf ("%s", ffstr); // Print the result of readglobal
										printf(";");
										break;
							}
							else if (var[i] == 32) { // stop on space and not print
                    break;
							}
						}
						
	} else {
		printf ("%c", (char)x);
	}
}
//printf ("\n\nVar is:%s", var);
return 0;
}

int readcss(char * cfile, char * gfile) {
	//cfile is css to variable-ize
	//gfile is global var file
	FILE *fp;
	fp = fopen (cfile, "r");
	int i;
	char var[100];
	while (1) {
		int x = fgetc(fp);
		if (x == -1) { //check for end of file (null)
	break;
		} else if (x == 36) { //Check for $ (variable)
	//printf ("FOUND VAR!");
	var[0] = 36;
	for (i = 1; i <= 100; i++) {
		var[i] = (char)fgetc(fp); // Load var name into var array
		if (var[i] == 59) { //stop on ; and print variable then print ;
                    var[i] = NULL;
										char *ffstr = readglobal(var, gfile);  //pass pointer (array)
	//printf ("%d%c", strlen(ffstr), ffstr[strlen(ffstr)-1]);
	if (ffstr[strlen(ffstr)-1] == 13 || ffstr[strlen(ffstr)-1] == 10) {
		ffstr[strlen(ffstr)-1] = NULL; // Chomp off the last \n if there is any (Set char to null)
	}
	printf ("%s", ffstr); // Print the result of readglobal
	printf(";");
	break;
		}
		else if (var[i] == 32) { // stop on space and not print
                    break;
		}
	}
	
		} else {
			printf ("%c", (char)x);
		}
	}
	//printf ("\n\nVar is:%s", var);
	return 0;
}

char * readglobal(char * fstr, char * gfile) {
	//gfile is global.txt or equivalant
	char find[] = " $abc";
	
	FILE *fp;
	fp = fopen (gfile, "r");
	char str[1024];
	int i;
	while (1) {
		if (fgets (str, 1024, fp) == NULL) { //Check for end of file
            break;
		}
		//printf ("%s", str);
		char *str2 = strtok(str, "=");
		if (strcmp(fstr, str2) == 0) { // If find matches the beginning of the line, give the end of the line
            //printf ("SAME");
            char *str3 = strtok(NULL, "="); //Remember, a pointer is an array
						return str3;
		}
		//printf ("%s\n", strtok(NULL, "="));
	}
}

int main(int argc, char * argv[])
{
	if (argc < 2) {
		printf ("NOT ENOUGH VARIABLES\n");
		printf ("cssvar [Location of CSS file (opetional)] <Location of global.txt>");
		return 1;
	} else if (argc == 2) {
		readcssstd(argv[1]);
	} else {
		readcss(argv[1], argv[2]);
	}
}