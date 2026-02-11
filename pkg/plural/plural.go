package plural

func Make(n int, one string, few string, many string) string {
	n = n % 100
	if n >= 11 && n <= 19 {
		return many
	}

	switch n % 10 {
	case 1:
		return one
	case 2, 3, 4:
		return few
	default:
		return many
	}
}
