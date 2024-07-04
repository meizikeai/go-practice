package tool

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	mathrand "math/rand"
)

func Contain(arr []string, element string) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}

func MarshalJson(date any) []byte {
	res, err := json.Marshal(date)

	if err != nil {
		fmt.Println(err)
	}

	return res
}

func UnmarshalJson(date string) map[string]any {
	var res map[string]any

	_ = json.Unmarshal([]byte(date), &res)

	return res
}

func GetRandmod(length int) int64 {
	result := int64(0)
	res, err := rand.Int(rand.Reader, big.NewInt(int64(length)))

	if err != nil {
		return result
	}

	return res.Int64()
}

func IntToString(value int64) string {
	v := strconv.FormatInt(value, 10)

	return v
}

func StringToInt(value string) int64 {
	res, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		res = 0
	}

	return res
}

var compileRegex = regexp.MustCompile(`\D`)

func ClearNotaNumber(str string) string {
	return compileRegex.ReplaceAllString(str, "")
}

func HandleEscape(source string) string {
	var j int = 0

	if len(source) == 0 {
		return ""
	}

	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)

	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte

		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}

		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}

	return string(desc[0:j])
}

func GenerateRandomNumber(start int, end int, count int) string {
	if end < start || (end-start) < count {
		return ""
	}

	result := make([]string, 0)
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	for len(result) < count {
		r := r.Intn((end - start)) + start
		num := IntToString(int64(r))

		exist := false

		for _, v := range result {
			if v == num {
				exist = true
				break
			}
		}

		if exist == false {
			result = append(result, num)
		}
	}

	return strings.Join(result, "")
}

func RemoveDuplicateElement(strs []string) []string {
	result := make([]string, 0, len(strs))
	temp := map[string]struct{}{}

	for _, item := range strs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func IsSlice(v any) bool {
	kind := reflect.ValueOf(v).Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		return true
	}

	return false
}

func StringToArray(data string) []string {
	result := []string{}

	if len(data) > 0 {
		result = strings.Split(data, ",")
	}

	return result
}

func GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
