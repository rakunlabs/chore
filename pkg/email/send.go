package email

import (
	"io"
	"time"

	gomail "gopkg.in/mail.v2"
)

var EmailTimeout = 1 * time.Minute

type Client struct {
	d *gomail.Dialer
}

type Attach struct {
	FileName string
	Content  io.Reader
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

// Send with headers.
// Headers should not be empty string array!
func (c *Client) Send(msg []byte, headers map[string][]string, attachments []Attach) error {
	m := gomail.NewMessage()

	m.SetHeaders(headers)

	for _, attach := range attachments {
		m.AttachReader(attach.FileName, attach.Content)
	}

	m.SetBody("text/html", string(msg))

	return c.d.DialAndSend(m)
}
