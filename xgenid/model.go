package xgenid

import (
	"time"
)

type idGenerator struct {
	Id         int
	Name       string
	Version    int
	ExpireTime time.Time
}

var idGeneratorTableName = "db_id_generator"

func (*idGenerator) TableName() string {
	if schemaName == "" {
		return idGeneratorTableName
	}
	return schemaName + "." + idGeneratorTableName
}
