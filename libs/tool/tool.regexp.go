package tool

import (
	"regexp"
	// "github.com/dlclark/regexp2"
)

// 检查密码强度
var checkPasswordRegexp = handleCheckPasswordRegexp()

func handleCheckPasswordRegexp() []*regexp.Regexp {
	var data = []string{
		`[0-9]+`,
		`[a-z]+`,
		`[A-Z]+`,
		`[~!@#$%^&amp;*?_-]+`,
	}
	// match, _ := regexp.MatchString(`^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,10}$`, password)

	result := []*regexp.Regexp{}

	for _, v := range data {
		re, err := regexp.Compile(v)

		if err != nil {
			// fmt.Println(err)
			continue
		}

		result = append(result, re)
	}

	return result
}

func CheckPassword(password string, min, max int) int {
	level := 0

	if len(password) < min {
		return -1
	}

	if len(password) > max {
		return 5
	}

	for _, v := range checkPasswordRegexp {
		match := v.MatchString(password)

		if match == true {
			level++
		}
	}

	return level
}

// 清除非数字符
var notANumberRegex = regexp.MustCompile(`\D`)

func ClearNotANumber(str string) string {
	return notANumberRegex.ReplaceAllString(str, "")
}
