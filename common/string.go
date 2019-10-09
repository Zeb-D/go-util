package common

//StringReverse support uft-8 string
func StringReverse(src string) string {
	bytes := []rune(src)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	return string(bytes)
}

func SubString(src string, begin, end int) string {
	//you will be args check
	rs := []rune(src)
	return string(rs[begin:end])
}
