package tool

import (
	"github.com/sqids/sqids-go"
)

var s, _ = sqids.New(sqids.Options{
	Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123567890", // 注意，与默认的值不同，故意少 4 这个数字
	MinLength: 6,
	Blocklist: []string{},
})

func HashidsEncode(uid int64) string {
	r, _ := s.Encode([]uint64{uint64(uid)})
	return r
}

func HashidsDecode(id string) []uint64 {
	res := s.Decode(id)
	return res
}
