package main

import (
	"fmt"
)

func main() {
	input := []int{10, 6, 9, -3, 23, -1, 34, 56, 67, -1, -4, -8, -2, 9, 10, 34, 67}
	size := 2

	res := make([]int, 0)
	window := make([]int, 0)

	for i := 0; i < size; i++ {
		if len(window) != 0 && input[i] >= input[window[len(window)-1]] {
			window = window[:len(window)-1]
		}
		window = append(window, i)
	}

	fmt.Println(window)
	res = append(res, input[window[0]])

	for i := size; i < len(input); i++ {
		fmt.Println(fmt.Sprintf("%d: %v", i, window))
		if len(window) != 0 && input[i] >= input[window[len(window)-1]] {
			window = window[:len(window)-1]
		}
		fmt.Println(fmt.Sprintf("%d: %v", i, window))
		if len(window) != 0 && window[0] <= (i-size) {
			window = window[1:]
		}
		fmt.Println(fmt.Sprintf("%d: %v", i, window))
		window = append(window, i)
		res = append(res, input[window[0]])
	}

	// [10, 9, 9, 23, 23, 34, 56, 67, 67, -1, -4, -2, 9, 10, 34, 67]
	// [10 9 9 23 23 34 56 67 67 -1 -4 -2 9 10 34 67]
	fmt.Println(res)
}
