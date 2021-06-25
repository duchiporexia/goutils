package main

import (
	"flag"
	gen "github.com/duchiporexia/goutils/xgormgen/internal"
	"github.com/duchiporexia/goutils/xstrings"
	"log"
	"os"
)

type config struct {
	output   string
	entities []string
}

var cnf config

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func parseFlags() {
	var output string
	var entities arrayFlags
	flag.Var(&entities, "entity", "[Required] entity list")
	flag.StringVar(&output, "output", "", "[Optional] The name of the output file")
	flag.Parse()

	if len(entities) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if output == "" {
		output = xstrings.ToSnakeCase(entities[0]) + "_gorm_gen.go"
	}

	cnf = config{
		output:   output,
		entities: entities,
	}
}

func main() {
	//-structs User -output user_gen.go
	parseFlags()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("entities:%s\n", cnf.entities)
	log.Printf("output:%s\n", cnf.output)

	parser := gen.NewParser()
	parser.ParseDir(wd)

	generator := gen.NewGenerator(cnf.output)
	if err := generator.Init(parser, cnf.entities); err != nil {
		log.Fatalf("Error Initializing Generator: %v", err.Error())
	}
	if err := generator.Generate(); err != nil {
		log.Fatalf("Error Generating file: %v", err.Error())
	}
	log.Printf("done\n")
}
