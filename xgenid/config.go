package xgenid

type IdGeneratorConfig struct {
	SchemaName   string `yaml:"schemaName" env:"SCHEMA_NAME" env-default:""`
	ExpireInSecs string `yaml:"expireInSecs" env:"EXPIRE_IN_SECONDS" env-default:"600"`
}
