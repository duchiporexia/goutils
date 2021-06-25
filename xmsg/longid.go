package xmsg

import (
	"strconv"
	"strings"
)

//go:generate msgp --tests=false

//msgp:shim LongID as:string using:(LongID).String/ParseString
type LongID int64

func (s LongID) ToInt64() int64 {
	return int64(s)
}

func (s LongID) String() string {
	return strconv.FormatInt(int64(s), 10)
}

func ParseString(s string) LongID {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return LongID(id)
}

func (t LongID) MarshalJSON() ([]byte, error) {
	str := strconv.FormatInt(int64(t), 10)
	return []byte(str), nil
}

func (t *LongID) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return err
	}
	*t = LongID(i)
	return nil
}
