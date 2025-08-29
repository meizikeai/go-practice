package tool

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unsafe"

	mathrand "math/rand"
)

type Units struct{}

func NewUnits() *Units {
	return &Units{}
}

// 字符串转换
// https://github.com/spf13/cast

func (u *Units) MarshalJson(date any) string {
	res, err := json.Marshal(date)

	if err != nil {
		return ""
	}

	return string(res)
}

func (u *Units) UnmarshalJson(date string) map[string]any {
	res := make(map[string]any, 0)

	_ = json.Unmarshal([]byte(date), &res)

	return res
}

func (u *Units) CheckPassword(password string, min, max int) int {
	level := 0

	if len(password) < min {
		return -1
	}

	if len(password) > max {
		return 5
	}

	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&amp;*?_-]+`}

	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, password)

		if match == true {
			level++
		}
	}

	return level
}

var ranSource = mathrand.NewSource(time.Now().UnixNano())

const (
	letterIdBits = 6
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func getLetter(types int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

	switch types {
	case 1:
		str = "1234567890"
	case 2:
		str = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	case 3:
		str = "abcdefghijklmnopqrstuvwxyz"
	case 4:
		str = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	return str
}

func (u *Units) CreateRandom(types, length int) string {
	b := make([]byte, length)
	letters := getLetter(types)

	for i, cache, remain := length-1, ranSource.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = ranSource.Int63(), letterIdMax
		}

		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}

		cache >>= letterIdBits
		remain--
	}

	result := *(*string)(unsafe.Pointer(&b))
	// fmt.Println(result)

	return result
}

func (u *Units) GenerateRandomNumber(start, end, count int) ([]int, error) {
	var result []int

	for range count {
		rangeBig := big.NewInt(int64(end - start + 1))
		n, err := rand.Int(rand.Reader, rangeBig)

		if err != nil {
			return nil, err
		}

		num := int(n.Int64()) + start

		result = append(result, num)
	}

	return result, nil
}

func (u *Units) RemoveDuplicateElement(strs []string) []string {
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

func (u *Units) IsSlice(v any) bool {
	kind := reflect.ValueOf(v).Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		return true
	}

	return false
}

func (u *Units) ArrayIntToString(array []int64, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(array), " ", delim, -1), "[]")
}
