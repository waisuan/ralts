package main

import (
	"sort"
)

func main() {
}

func findSumOfThree(nums []int, target int) bool {
	sort.Ints(nums)

	for i := 0; i < len(nums); i++ {
		a := nums[i]

		j := i + 1
		if j == len(nums)-1 {
			return false
		}

		k := len(nums) - 1
		for j < k {
			if a+nums[j]+nums[k] < target {
				j++
			} else if a+nums[j]+nums[k] > target {
				k--
			} else {
				return true
			}
		}
	}

	return false
}
