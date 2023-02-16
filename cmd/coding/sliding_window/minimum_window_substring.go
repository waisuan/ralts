package main

import (
	"fmt"
	"math"
	"strings"
)

func minWindow(s, t string) string {
	// Empty string scenario
	if t == "" {
		return ""
	}

	// Creating the two hash maps
	rCount := make(map[byte]int)
	window := make(map[byte]int)

	// Populating the rCount
	for c, _ := range t {
		if _, ok := rCount[t[c]]; ok {
			rCount[t[c]] = 1 + rCount[t[c]]
		} else {
			rCount[t[c]] = 1
		}
	}

	// Setting up the conditional variables
	current, required := 0, len(rCount)

	// Setting up a variable containing the result's starting and ending point with default values and a length variable
	res, resLen := []int{-1, -1}, math.MaxInt64

	// Setting up the sliding window pointers
	l := 0
	for r, _ := range s {
		c := s[r]

		// Populating the window hashmap
		if _, ok := window[c]; ok {
			window[c] = 1 + window[c]
		} else {
			window[c] = 1
		}

		// Updating the current variable
		if _, ok := rCount[c]; ok && window[c] == rCount[c] {
			current += 1
		}

		// Sliding Window in working
		for current == required {
			// Update our result
			if (r - l + 1) < resLen {
				res = []int{l, r}
				resLen = (r - l + 1)
			}

			// Pop from the left of our window
			window[s[l]] -= 1
			if _, ok := rCount[s[l]]; ok && window[s[l]] < rCount[s[l]] {
				current -= 1
			}
			l += 1
		}
	}
	l, r := res[0], res[1]

	if resLen != math.MaxInt64 {
		return s[l : r+1]
	} else {
		return ""
	}
}

// Driver code
func main() {
	s := []string{"PATTERN", "LIFE", "ABRACADABRA", "STRIKER", "DFFDFDFVD"}
	t := []string{"TN", "I", "ABC", "RK", "VDD"}

	for i, _ := range s {
		fmt.Printf("%d.\ts: \"%s\"\n\tt: \"%s\"\n\tThe minimum substring containing \"%s\" is \"%s\"\n", i+1, s[i], t[i], t[i], minWindow(s[i], t[i]))
		fmt.Printf("%s\n", strings.Repeat("-", 100))
	}
}
