char startsWith(char* str, int n, char* sub, int m) {
	int i, j;
	if (m > n) return 0;
	for (i = 0; i < n; i++) {
		if (str[i] != sub[i]) return 0;
	}
	return 1;
}