package internal

import (
	"github.com/pkg/errors"
	"go/format"
	"golang.org/x/tools/imports"
	"gorm.io/gorm/schema"

	"bytes"
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
)

type structHelpers struct {
	Titelize func(string) string
}

type fieldConfig struct {
	FieldName  string
	ColumnName string
	FieldType  string
}

type structConfig struct {
	StructName       string
	QueryBuilderName string
	Fields           []fieldConfig
}

type structsConfig struct {
	PkgName string
	Helpers structHelpers
	Structs []structConfig
}

// The Generator is the one responsible for generating the code, adding the imports, formating, and writing it to the file.
type Generator struct {
	buf        *bytes.Buffer
	outputFile string
	config     structsConfig
}

// NewGenerator function creates an instance of the generator given the name of the output file as an argument.
func NewGenerator(outputFile string) *Generator {
	return &Generator{
		buf:        new(bytes.Buffer),
		outputFile: outputFile,
	}
}

func Title(str string) string {
	return strings.Title(str)
}

// Init function should be called before any other function is called. It takes a parser that has already parsed the directory
// that contains the types we want to generate code for. It also takes the name of the structs that we want to generate code for.
func (g *Generator) Init(parser *Parser, structs []string) error {
	if err := g.validateStructs(parser, structs); err != nil {
		return err
	}
	g.config = structsConfig{
		PkgName: parser.pkgName,
		Helpers: structHelpers{
			Titelize: strings.Title,
		},
	}
	for _, st := range structs {
		g.config.Structs = append(g.config.Structs, *g.buildConfig(parser, st))
	}

	return nil
}

func (g *Generator) validateStructs(parser *Parser, structs []string) error {
	for _, st := range structs {
		if parser.GetTypeByName(st) == nil {
			return fmt.Errorf("Type %v is not found", st)
		}
	}
	return nil
}

func (g *Generator) buildConfig(parser *Parser, structName string) *structConfig {
	cnf := &structConfig{
		StructName:       structName,
		QueryBuilderName: fmt.Sprintf("%sQueryBuilder", structName),
	}
	structType := parser.GetTypeByName(structName)
	cnf.Fields = g.buildFieldConfig(parser, structType)
	return cnf
}

func (g *Generator) parseGormStructTag(tagLine string) map[string]string {
	ret := make(map[string]string)
	tag := reflect.StructTag(strings.Trim(tagLine, "`"))
	if section, ok := tag.Lookup("gorm"); ok {
		if section == "-" {
			ret["-"] = "-"
			return ret
		}
		parts := strings.Split(section, ";")
		for _, part := range parts {
			kv := strings.Split(part, ":")
			ret[kv[0]] = strings.Join(kv[1:], ":")
		}
	}
	return ret
}

func (g *Generator) buildFieldConfig(parser *Parser, structType *types.Struct) []fieldConfig {
	fields := []fieldConfig{}
	ns := schema.NamingStrategy{}
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		if !field.Exported() {
			continue
		}
		tag := g.parseGormStructTag(structType.Tag(i))
		if _, ok := tag["-"]; ok {
			continue
		}
		if field.Anonymous() {
			fields = append(fields, g.buildFieldConfig(parser, field.Type().Underlying().(*types.Struct))...)
			continue
		}
		// A hack to ignore non-stdlib named types for now
		fieldType := field.Type().String()
		if strings.Contains(fieldType, "/") {
			switch fieldType {
			case "github.com/duchiporexia/goutils/xmsg.LongID":
				fieldType = "int64"
			case "github.com/duchiporexia/goutils/xmsg.DateTime":
				fieldType = "time.Time"
			case "*github.com/duchiporexia/goutils/xmsg.DateTime":
				fieldType = "*time.Time"
			default:
				continue
			}
		}
		columnName := ns.ColumnName("", field.Name())
		if cname, ok := tag["column"]; ok {
			columnName = cname
		}
		fields = append(fields, fieldConfig{
			FieldName:  field.Name(),
			ColumnName: columnName,
			FieldType:  fieldType,
		})
	}
	return fields
}

// Generate executes the template and store it in an internal buffer.
func (g *Generator) Generate() error {
	var buf bytes.Buffer

	WriteGenerateFunc(&buf, &g.config)

	src, err := imports.Process(g.outputFile, buf.Bytes(), nil)
	if err != nil {
		return errors.Wrap(err, "unable to gofmt")
	}
	formatedOutput, err := format.Source(g.buf.Bytes())
	if err != nil {
		return err
	}
	g.buf = bytes.NewBuffer(formatedOutput)

	//src := buf.Bytes()

	if err := ioutil.WriteFile(g.outputFile, src, 0644); err != nil {
		return errors.Wrap(err, "writing output")
	}

	return nil
}

// Imports function adds the missing imports in the generated code.
func (g *Generator) Imports() error {
	wd, err := os.Getwd()
	formatedOutput, err := imports.Process(wd, g.buf.Bytes(), nil)
	if err != nil {
		return err
	}
	g.buf = bytes.NewBuffer(formatedOutput)
	return nil
}

// Flush function writes the output to the output file.
func (g *Generator) Flush() error {
	return ioutil.WriteFile(g.outputFile, g.buf.Bytes(), 0644)
}
