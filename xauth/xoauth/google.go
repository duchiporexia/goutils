package xoauth

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/duchiporexia/goutils/xerr"
	"github.com/duchiporexia/goutils/xjson"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	oauthGoogleCfg *oauth2.Config
	paths          = [][]string{
		{"id"},
		{"email"},
		{"verified_email"},
		{"name"},
		{"given_name"}, {"family_name"},
	}
)

func Init(cfg *XAuthGoogleConfig) {
	oauthGoogleCfg = &oauth2.Config{
		RedirectURL:  cfg.CallbackUrl,
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

var (
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

func RedirectUrlGoogle(ctx *fiber.Ctx) error {
	url := oauthGoogleCfg.AuthCodeURL(oauthStateString)
	return ctx.SendString(url)
}

func GetUserInfoGoogle(code string, state string) (string, error) {
	if state != oauthStateString {
		return "", fmt.Errorf("invalid oauth state")
	}
	token, err := oauthGoogleCfg.Exchange(oauth2.NoContext, code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %s", err)
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %s", err)
	}
	defer response.Body.Close()
	userInfo, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading response body: %s", err)
	}
	return string(userInfo), nil
}

func ParseUserInfoGoogle(content string) (*OauthInfo, error) {
	var info OauthInfo
	var err error
	xjson.ParseJson(paths, []byte(content), func(idx int, value []byte, vt jsonparser.ValueType, er error) {
		if er != nil {
			err = er
		}
		switch idx {
		case 0:
			info.Uuid = string(value)
		case 1:
			info.Email = string(value)
		case 2:
			v, _ := jsonparser.ParseBoolean(value)
			info.VerifiedEmail = v
		case 3:
			info.Name = string(value)
		case 4:
			info.FirstName = string(value)
		case 5:
			info.LastName = string(value)
		}
	})
	if err != nil || info.Uuid == "" || info.Email == "" {
		return nil, xerr.ErrOauthInfo
	}
	info.Email = strings.ToLower(info.Email)
	return &info, nil
}
