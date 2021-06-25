package xemail

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"time"
)

type EmailConfig struct {
	Host     string `yaml:"host" env:"HOST" env-default:""`
	Port     string `yaml:"port" env:"PORT" env-default:"587"`
	User     string `yaml:"user" env:"USER" env-default:""`
	Password string `yaml:"password" env:"PASS" env-default:""`
	From     string `yaml:"from" env:"FROM" env-default:""`
}

var pool *email.Pool

var EmailFrom string

func Init(cfg *EmailConfig) {
	EmailFrom = cfg.From
	plainAuth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	p, err := email.NewPool(cfg.Host+":"+cfg.Port, 6, plainAuth)
	if err != nil {
		panic(err)
	}
	pool = p
}

func SendEmail(email *email.Email) error {
	return pool.Send(email, 5*time.Second)
}
