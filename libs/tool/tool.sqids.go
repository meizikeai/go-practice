package tool

import (
	"github.com/sqids/sqids-go"
)

type Sqids struct{}

func NewSqids() *Sqids {
	return &Sqids{}
}

var ids, _ = sqids.New(sqids.Options{
	Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123567890", // 注意，与默认的值不同，故意少 4 这个数字
	MinLength: 6,
	Blocklist: []string{},
})

func (s *Sqids) SqidsEncode(uid int64) string {
	r, _ := ids.Encode([]uint64{uint64(uid)})
	return r
}

func (s *Sqids) SqidsDecode(id string) []uint64 {
	res := ids.Decode(id)
	return res
}
