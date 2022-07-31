package algo

func CheckNumber(number int64, req string) bool {
	return GetNumberOfString(number) == req
}

func GetNumberOfString(m int64) string {

	ans := []byte{}

	for m > 0 {
		ans = append(ans, byte(m%26+'a'))
		m /= 26
	}

	for i := 0; i < len(ans)/2; i++ {
		ans[i], ans[len(ans)-i-1] = ans[len(ans)-i-1], ans[i]
	}

	return string(ans)
}
