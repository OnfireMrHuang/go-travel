package algorithm

import (
	"strings"
)

// 通过字符c分割字符串str
func SplitStr(c byte, str string) ([]string, error) {

	// 判断如果没有包含字符，则直接返回空切片
	if strings.Contains(str, string(c)) == false {
		return []string{}, nil
	}
	var list []string
	var tmp []byte
	for i := 0; i < len(str); i++ {
		// 判断字符是不是等于分割字符
		if str[i] == c {
			if string(tmp) != "" {
				list = append(list, string(tmp))
			}
			tmp = []byte{}
		} else {
			tmp = append(tmp, str[i])
		}
	}
	if string(tmp) != "" {
		list = append(list, string(tmp))
	}
	return list, nil
}
