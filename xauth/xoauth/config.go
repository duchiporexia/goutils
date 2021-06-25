package xoauth

type XAuthGoogleConfig struct {
	ClientId     string `yaml:"clientId" env:"CLIENT_ID" env-default:""`
	ClientSecret string `yaml:"clientSecret" env:"CLIENT_SECRET" env-default:""`
	CallbackUrl  string `yaml:"callbackUrl" env:"CB_URL" env-default:""`
}
