package main

func main() {
	solution(0, "wsdawsdassw")
}

func solution(n int, control string) int {
	for _, c := range control {
		switch c {
		case 'w':
			n++
		case 's':
			n--
		case 'd':
			n += 10
		case 'a':
			n -= 10
		}
	}

	return n
}
