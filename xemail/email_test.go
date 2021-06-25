package xemail

import (
	"github.com/duchiporexia/goutils/xconfig"
	"github.com/jordan-wright/email"
	"testing"
)

func TestSendEmail(t *testing.T) {
	var cfg EmailConfig
	xconfig.LoadConfig("", &cfg)
	Init(&cfg)
	e := email.NewEmail()
	e.From = cfg.From
	e.To = []string{"David Xia <duchiporexia@gmail.com>"}
	e.Bcc = []string{}
	e.Cc = []string{"duchipore@qq.com"}
	e.Subject = "Signup [TEST]"
	//e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h>This is Email Demo</h>")
	SendEmail(e)
}
