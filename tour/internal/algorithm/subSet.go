package algorithm

func SubSets1(str string) ([]string, int) {

	if len(str) < 1 {
		return []string{""}, 1
	}

	// 取出字符串最后一个字符
	lastChar := str[len(str)-1]

	// 取除了最后一个字符的子字符串
	subStr := str[:len(str)-1]

	// fmt.Printf("lastChar %v, subStr %v \n", string(lastChar), string(subStr))

	// 递归调用求子字符串的函数
	strList, n := SubSets1(subStr)
	size := len(strList)
	for i := 0; i < size; i++ {
		tmpStr := strList[i] + string(lastChar)
		strList = append(strList, tmpStr)
		n++
	}
	return strList, n
}

var subSets []string
var n int

func SubSets2(str string) ([]string, int) {
	var subStr string
	backTrack(str, 0, subStr)
	return subSets, n
}

func backTrack(str string, start int, subStr string) {

	// 前序遍历的位置
	subSets = append(subSets, subStr)
	n++

	for i := start; i < len(str); i++ {
		subStr = subStr + string(str[i])
		backTrack(str, i+1, subStr)
		subStr = subStr[:len(subStr)-1]
	}
}
