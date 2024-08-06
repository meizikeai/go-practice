package tool

import (
	"regexp"
	// "github.com/dlclark/regexp2"
)

type Regexp struct{}

func NewRegexp() *Regexp {
	return &Regexp{}
}

var notANumberRegexp = regexp.MustCompile(`\D`)
var noSpaceRegexp = regexp.MustCompile(`\s+`)
var noLineFeedRegexp = regexp.MustCompile(`\n|\r|\t`)

func (r *Regexp) CleanNotANumber(str string) string {
	return notANumberRegexp.ReplaceAllString(str, "")
}

func (r *Regexp) CleanSpace(str string) string {
	return noSpaceRegexp.ReplaceAllString(str, "")
}

func (r *Regexp) CleanLineFeed(str string) string {
	return noLineFeedRegexp.ReplaceAllString(str, " ")
}

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
			continue
		}

		result = append(result, re)
	}

	return result
}

func (r *Regexp) CheckPassword(password string, min, max int) int {
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
