package api

import (
	"crypto/md5"
	"fmt"
)

type Auth struct {
	CP     string `json:"cp"`
	Seckey string `json:"seckey"`
}

// Sign ctime 请求时间(时间戳，到秒)
func (a Auth) Sign(ctime int64) string {
	s := fmt.Sprintf("%s%d%s", a.CP, ctime, a.Seckey)
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
