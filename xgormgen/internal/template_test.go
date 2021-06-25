package internal

import (
	"fmt"
	"testing"
)

//go:generate qtc -dir=.

func TestGenerateFunc(t *testing.T) {
	data := structsConfig{
		PkgName: "gormgen",
		Structs: []structConfig{{
			StructName:       "User",
			QueryBuilderName: "UserQueryBuilder",
			Fields: []fieldConfig{{
				FieldName:  "Id",
				ColumnName: "id",
				FieldType:  "int",
			}, {
				FieldName:  "Name",
				ColumnName: "name",
				FieldType:  "string",
			},
			},
		}},
	}
	fmt.Printf("%s", GenerateFunc(&data))
}
