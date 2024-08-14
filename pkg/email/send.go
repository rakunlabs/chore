package email

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"

	gomail "gopkg.in/mail.v2"
)

var EmailTimeout = 1 * time.Minute

type Client struct {
	d *gomail.Dialer
}

type Attach struct {
	FileName string `json:"filename"`
	Content  string `json:"content"`
	Filetype string `json:"filetype"`
}

func NewClient(host string, port int, noAuth bool, mail, password string) Client {
	var mailDialer *gomail.Dialer

	if noAuth {
		mailDialer = &gomail.Dialer{Host: host, Port: port}
	} else {
		mailDialer = gomail.NewDialer(host, port, mail, password)
	}

	mailDialer.Timeout = EmailTimeout

	return Client{mailDialer}
}

// fileFromTextConvert converting

func (c *Client) fileFromTextConvert(attach Attach) Attach {
	var newAttach Attach
	var newContent string
	switch attach.Filetype {
	case "csv":
		newAttach.FileName = attach.FileName + ".csv"
		rows := strings.Split(attach.Content, " ")
		for _, row := range rows {
			newContent = newContent + row + "\n"
		}
		newAttach.Content = newContent
		newAttach.Filetype = "csv"
	default:
		newAttach = attach
	}
	return newAttach
}

// Send with headers.
// Headers should not be empty string array!
func (c *Client) Send(msg []byte, headers map[string][]string) error {
	m := gomail.NewMessage()

	m.SetHeaders(headers)

	body := struct {
		Attachments []Attach `json:"attachments"`
		Body        []byte   `json:"body"`
	}{}

	err := json.Unmarshal(msg, &body)
	if err != nil {
		log.Info().Msgf("error: %v", fmt.Errorf("failed to unmarshal attachments and body: %w", err))
	}

	for _, el := range body.Attachments {
		as := c.fileFromTextConvert(el)
		m.AttachReader(as.FileName, strings.NewReader(as.Content))
	}
	if len(body.Attachments) == 0 {
		m.SetBody("text/html", string(msg))
		return c.d.DialAndSend(m)
	}

	m.SetBody("text/html", string(body.Body))

	return c.d.DialAndSend(m)
}
