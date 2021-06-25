package main

import (
	"flag"
	bh "github.com/duchiporexia/goutils/batch_handler/internal"
	"log"
	"os"
)

type config struct {
	name        string
	keyType     string
	paramsType  string
	valueType   string
	isBatched   bool
	handlerType string
}

var cnf config

func parseFlags() {
	flag.StringVar(&cnf.name, "name", "", "[Required] The name of batch handler")
	flag.StringVar(&cnf.keyType, "keyType", "", "[Optional] The key type")
	flag.StringVar(&cnf.paramsType, "paramsType", "", "[Optional] The params type")
	flag.StringVar(&cnf.valueType, "valueType", "", "[Optional] The valueType type")
	flag.BoolVar(&cnf.isBatched, "isBatched", true, "[Optional] is batched or not")
	flag.StringVar(&cnf.handlerType, "handlerType", "", "[Required] Handler Type: one of [create,update,read,delete]")
	flag.Parse()

	if cnf.name == "" || cnf.handlerType == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	//go:generate go run github.com/duchiporexia/goutils/batch_handler -name TReadNpHandler -keyType string -valueType *servicehub/common/batch_handler.Dog -handlerType read
	parseFlags()
	var handlerType bh.HandlerType
	switch cnf.handlerType {
	case "create":
		handlerType = bh.HandlerTypeCreate
	case "read":
		handlerType = bh.HandlerTypeRead
	case "update":
		handlerType = bh.HandlerTypeUpdate
	case "delete":
		handlerType = bh.HandlerTypeDelete
	default:
		log.Fatalln("unknown handler type")
	}

	if err := bh.Generate(cnf.name, cnf.keyType, cnf.paramsType, cnf.valueType, cnf.isBatched, handlerType); err != nil {
		log.Fatalf("Error : %v", err.Error())
	}
}
