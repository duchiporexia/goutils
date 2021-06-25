package batch_handler

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

//go:generate qtc -dir=.

type handlerData struct {
	Package    string
	Name       string
	KeyType    *goType
	ParamsType *goType
	ValueType  *goType

	IsBatched   bool
	HandlerType HandlerType
}

func checkParameters(name string, keyType string, paramsType string, valueType string, handlerType HandlerType) error {
	if name == "" {
		return fmt.Errorf("empty name")
	}

	if !handlerType.isValid() {
		return fmt.Errorf("invalid handlerType")
	}
	switch handlerType {
	case HandlerTypeCreate:
		if keyType != "" {
			return fmt.Errorf("keyType should be empty")
		}
		if paramsType == "" {
			return fmt.Errorf("empty paramsType")
		}
	case HandlerTypeRead:
		if keyType == "" {
			return fmt.Errorf("empty keyType")
		}
		if valueType == "" {
			return fmt.Errorf("empty valueType")
		}
	case HandlerTypeUpdate:
		if keyType == "" {
			return fmt.Errorf("empty keyType")
		}
	case HandlerTypeDelete:
		if keyType == "" {
			return fmt.Errorf("empty keyType")
		}
	}
	return nil
}

type HandlerType int

const (
	HandlerTypeCreate HandlerType = 1 // No Key | No Cache | Has Params
	HandlerTypeRead   HandlerType = 2 // Has Key | CacheBatchGet, CacheBatchSet | Has Value
	HandlerTypeUpdate HandlerType = 3 // Has Key | CacheDel
	HandlerTypeDelete HandlerType = 4 // Has Key | No Cache
)

func (o HandlerType) isValid() bool {
	return o >= 1 && o <= 4
}

func (o HandlerType) String() string {
	switch o {
	case HandlerTypeCreate:
		return "Create"
	case HandlerTypeRead:
		return "Read"
	case HandlerTypeUpdate:
		return "Update"
	case HandlerTypeDelete:
		return "Delete"
	default:
		return "unknown"
	}
}

func (o HandlerType) hasKey() bool {
	return o != HandlerTypeCreate
}

func (o HandlerType) hasCacheGet() bool {
	return o == HandlerTypeRead
}

func (o HandlerType) hasCacheSet() bool {
	return o == HandlerTypeRead
}

func (o HandlerType) hasCacheDel() bool {
	return o == HandlerTypeUpdate
}

type goType struct {
	Modifiers  string
	ImportPath string
	ImportName string
	Name       string
}

func (t *goType) String() string {
	if t.ImportName != "" {
		return t.Modifiers + t.ImportName + "." + t.Name
	}

	return t.Modifiers + t.Name
}

func (t *goType) IsPtr() bool {
	return strings.HasPrefix(t.Modifiers, "*")
}

func (t *goType) IsSlice() bool {
	return strings.HasPrefix(t.Modifiers, "[]")
}

var partsRe = regexp.MustCompile(`^([\[\]\*]*)(.*?)(\.\w*)?$`)

func parseType(str string) (*goType, error) {
	if str == "" {
		return nil, nil
	}
	parts := partsRe.FindStringSubmatch(str)
	if len(parts) != 4 {
		return nil, fmt.Errorf("type must be in the form []*github.com/import/path.Name")
	}

	t := &goType{
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

func GetGoTypeName(gotype *goType) string {
	if gotype == nil {
		return ""
	}
	return gotype.String()
}

func GetZeroValue(gotype *goType) string {
	if gotype == nil {
		return ""
	}
	if gotype.IsPtr() {
		return "nil"
	}
	basicType := gotype.Name
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
		return gotype.String() + "{}"
	}
}

func Generate(name string, keyType string, paramsType string, valueType string, isBatched bool, handlerType HandlerType) error {
	if err := checkParameters(name, keyType, paramsType, valueType, handlerType); err != nil {
		return err
	}
	wd, _ := os.Getwd()
	data, err := getHandlerData(name, keyType, paramsType, valueType, wd, isBatched, handlerType)
	if err != nil {
		return err
	}

	filename := strings.ToLower(data.Name) + "_gen.go"

	if err := writeTemplate(filepath.Join(wd, filename), data); err != nil {
		return err
	}

	return nil
}

func getHandlerData(name string, keyType string, paramsType string, valueType string, wd string, isBatched bool, handlerType HandlerType) (*handlerData, error) {
	var data handlerData

	genPkg := getPackage(wd)
	if genPkg == nil {
		return nil, fmt.Errorf("unable to find package info for " + wd)
	}

	var err error
	data.Name = name
	data.Package = genPkg.Name

	//////////////////////////////////////////////////////
	data.KeyType, err = parseType(keyType)
	if err != nil {
		return nil, fmt.Errorf("KeyType: %s", err.Error())
	}
	if data.KeyType != nil && genPkg.PkgPath == data.KeyType.ImportPath {
		data.KeyType.ImportName = ""
		data.KeyType.ImportPath = ""
	}
	//////////////////////////////////////////////////////
	data.ParamsType, err = parseType(paramsType)
	if err != nil {
		return nil, fmt.Errorf("ParamsType: %s", err.Error())
	}
	if data.ParamsType != nil && genPkg.PkgPath == data.ParamsType.ImportPath {
		data.ParamsType.ImportName = ""
		data.ParamsType.ImportPath = ""
	}
	//////////////////////////////////////////////////////
	data.ValueType, err = parseType(valueType)
	if err != nil {
		return nil, fmt.Errorf("ValueType: %s", err.Error())
	}
	if data.ValueType != nil && genPkg.PkgPath == data.ValueType.ImportPath {
		data.ValueType.ImportName = ""
		data.ValueType.ImportPath = ""
	}

	data.IsBatched = isBatched
	data.HandlerType = handlerType

	return &data, nil
}

func getPackage(dir string) *packages.Package {
	p, _ := packages.Load(&packages.Config{
		Dir: dir,
	}, ".")

	if len(p) != 1 {
		return nil
	}

	return p[0]
}

func writeTemplate(filepath string, data *handlerData) error {
	var buf bytes.Buffer

	WriteGenerateBatchHandler(&buf, data)

	src, err := imports.Process(filepath, buf.Bytes(), nil)
	if err != nil {
		return errors.Wrap(err, "unable to gofmt")
	}

	//src := buf.Bytes()

	if err := ioutil.WriteFile(filepath, src, 0644); err != nil {
		return errors.Wrap(err, "writing output")
	}

	return nil
}

func lcFirst(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
