package main

func main() {
	num := []int{0, 1, 0, 10, 0, 1, 0, 10, 0, -1, -2, -1}
	solution(num)
}
func solution(numLog []int) string {
	answer := ""

	for i, _ := range numLog {
		if i == len(numLog)-1 {
			break
		}

		val := numLog[i+1] - numLog[i]

		if val == 1 {
			answer += "w"
		} else if val == -1 {
			answer += "s"
		} else if val == 10 {
			answer += "d"
		} else if val == -10 {
			answer += "a"
		}
	}

	return answer
}
