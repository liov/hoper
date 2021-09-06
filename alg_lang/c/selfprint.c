#include<stdio.h>
main() {
	char *s = "#include<stdio.h>%cmain(){%c%cchar *s=%c%s%c;%c%cprintf(s,10,10,9,34,s,34,10,9,10);%c}";
	printf(s, 10, 10, 9, 34, s, 34, 10, 9, 10);

	system("pause");
}
