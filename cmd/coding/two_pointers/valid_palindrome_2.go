package main

func main() {
	var s string

	i := 0
	j := len(s) - 1
	for i < j {
		if s[i] == s[j] {
			i++
			j--
			continue
		}
		if isPalindrome(s[0:i] + s[i+1:]) {

		} else if isPalindrome(s[0:j] + s[j+1:]) {

		} else {

		}
	}
}

func isPalindrome(s string) bool {
	i := 0
	j := len(s) - 1
	for i < j {
		if s[i] != s[j] {
			return false
		}
		i++
		j--
	}

	return true
}
