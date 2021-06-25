package xjson

import (
	"github.com/buger/jsonparser"
)

func ParseJson(paths [][]string, data []byte, cb func(int, []byte, jsonparser.ValueType, error)) {
	jsonparser.EachKey(data, cb, paths...)
}
