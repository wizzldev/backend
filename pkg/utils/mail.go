package utils

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/wizzldev/chat/pkg/configs"
	"gopkg.in/gomail.v2"
)

type Mail struct {
	fromAddress string
	toAddresses []string
	subject     string
	body        string
	isHTML      bool
}

func NewMail(from string, to ...string) *Mail {
	return &Mail{
		fromAddress: from,
		toAddresses: to,
	}
}

func (m *Mail) Subject(s string) *Mail {
	m.subject = s
	return m
}

func (m *Mail) Body(b string, html bool) *Mail {
	m.body = b
	m.isHTML = html
	return m
}

func (m *Mail) TemplateBody(t string, props map[string]string, otherwise string) *Mail {
	m.isHTML = true
	pwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Failed to get current working directory:", err)
		m.body = otherwise
		return m
	}

	path := pwd + "/templates/" + t + ".html"
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Println("Failed to open template file:", err)
		m.body = otherwise
		return m
	}

	dataStr := string(data)

	for k, v := range props {
		dataStr = strings.Replace(dataStr, "@"+k, v, 1)
	}

	m.body = dataStr
	return m
}

func (m *Mail) Send() error {
	gm := gomail.NewMessage()
	gm.SetHeader("From", m.fromAddress)
	gm.SetHeader("To", m.toAddresses...)
	gm.SetHeader("Subject", m.subject)

	contentType := "text/html"
	if !m.isHTML {
		contentType = "text/plain"
	}

	gm.SetBody(contentType, m.body)

	d := gomail.NewDialer(configs.Env.Email.SMTPHost, configs.Env.Email.SMTPPort, configs.Env.Email.Username, configs.Env.Email.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(gm); err != nil {
		return err
	}

	return nil
}
