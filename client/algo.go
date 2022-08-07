package main

import "github.com/smirzoavliyoev/word_of_wisdom_test/pkg/algo"

func binPow(a, b int64) int64 {
	if b <= 1 {
		return a
	}

	var z = binPow(a, b/2)

	if a%2 == 1 {
		return z * z * a
	}

	return z * z
}

// tried to use lexigrafical ordering of strings in specific length but its not good idea
// may be will try to use 3sum or APSP or OV^k
func GetNumberOfString(s string) (number int64) {

	var l int64 = 1
	var r int64 = binPow(26, int64(len(s)))

	//issue here...DDOS is possible if client hacks the algorithm...26 notation problem
	for l < r {
		var m int64 = (l + r) >> 1
		if algo.GetNumberOfString(m) < s {
			l = m + 1
		} else {
			r = m
		}
	}

	if algo.GetNumberOfString(l) < s {
		l++
	}

	return l
}
