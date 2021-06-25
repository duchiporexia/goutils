package utilscommon

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"regexp"
	"strings"
)

type GoType struct {
	Modifiers  string
	ImportPath string
	ImportName string
	Name       string
}

func (t *GoType) String() string {
	if t.ImportName != "" {
		return t.Modifiers + t.ImportName + "." + t.Name
	}

	return t.Modifiers + t.Name
}

func (t *GoType) IsPtr() bool {
	return strings.HasPrefix(t.Modifiers, "*")
}

func (t *GoType) IsSlice() bool {
	return strings.HasPrefix(t.Modifiers, "[]")
}

var partsRe = regexp.MustCompile(`^([\[\]\*]*)(.*?)(\.\w*)?$`)

func ParseType(str string) (*GoType, error) {
	if str == "" {
		return nil, nil
	}
	parts := partsRe.FindStringSubmatch(str)
	if len(parts) != 4 {
		return nil, fmt.Errorf("type must be in the form []*github.com/import/path.Name")
	}

	t := &GoType{
		Modifiers:  parts[1],
		ImportPath: parts[2],
		Name:       strings.TrimPrefix(parts[3], "."),
	}

	if t.Name == "" {
		t.Name = t.ImportPath
		t.ImportPath = ""
	}

	if t.ImportPath != "" {
		p, err := packages.Load(&packages.Config{Mode: packages.NeedName}, t.ImportPath)
		if err != nil {
			return nil, err
		}
		if len(p) != 1 {
			return nil, fmt.Errorf("not found")
		}

		t.ImportName = p[0].Name
	}

	return t, nil
}

func (s *GoType) ZeroValue() string {
	if s.IsPtr() {
		return "nil"
	}
	basicType := s.Name
	switch basicType {
	case "bool":
		return "false"
	case "string":
		return "\"\""
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"byte",
		"rune",
		"float32", "float64",
		"complex64", "complex128":
		return basicType + "(0)"
	default:
		return s.String() + "{}"
	}
}
