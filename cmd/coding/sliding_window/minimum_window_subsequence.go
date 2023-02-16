package main

import "math"

func main() {
}

func proc(str1 string, str2 string) string {
	sizeStr1, sizeStr2 := len(str1), len(str2)
	resLen := math.MaxInt64
	res := ""
	idx1, idx2 := 0, 0

	for idx1 < sizeStr1 {
		if str1[idx1] == str2[idx2] {
			idx2++
			if idx2 == sizeStr2 {
				start := idx1
				end := idx1 + 1
				idx2 -= 1
				for idx2 >= 0 {
					if str1[start] == str2[idx2] {
						idx2 -= 1
					}
					start -= 1
				}
				start += 1
				if end-start < resLen {
					resLen = end - start
					res = str1[start:end]
				}
				idx1 = start
				idx2 = 0
			}
		}
		idx1++
	}

	return res
}
