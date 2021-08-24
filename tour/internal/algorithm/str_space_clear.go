package algorithm

func ClearSpace(str string) (string, error) {

	var result []byte
	for i := 0; i < len(str); i++ {
		if str[i] != byte(' ') {
			result = append(result, str[i])
		}
	}
	return string(result), nil
}
