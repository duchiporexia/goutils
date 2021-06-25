package xjson

import (
	"fmt"
	"github.com/buger/jsonparser"
	"testing"
)

func TestDemo(t *testing.T) {
	data := []byte(`{
  "id": "534534",
  "email": "duchiporexia@gmail.com",
  "verified_email": true,
  "name": "David Xia",
  "given_name": "David",
  "family_name": "Xia",
  "picture": "https://lh3.googleusercontent.com/a-/t4d=s96-c",
  "locale": "en-GB"
}`)
	paths := [][]string{
		{"id"},
		{"email"},
		{"verified_email"},
		{"name"},
		{"given_name"},
		{"family_name"},
	}

	var obj struct {
		Uuid          string
		Email         string
		VerifiedEmail bool
		Name          string
		FirstName     string
		LastName      string
	}

	ParseJson(paths, data, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		switch idx {
		case 0:
			obj.Uuid = string(value)
		case 1:
			obj.Email = string(value)
		case 2:
			v, _ := jsonparser.ParseBoolean(value)
			obj.VerifiedEmail = v
		case 3:
			obj.Name = string(value)
		case 4:
			obj.FirstName = string(value)
		case 5:
			obj.LastName = string(value)
		}
	})

	fmt.Printf("%v\n", obj)
}
