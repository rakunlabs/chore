package email

import (
	gomail "gopkg.in/mail.v2"
)

type Client struct {
	d *gomail.Dialer
}

func NewClient(host string, port int, mail, password string) Client {
	d := gomail.NewDialer(host, port, mail, password)

	return Client{d}
}

func (c *Client) Send(msg []byte, headers map[string][]string) error {
	m := gomail.NewMessage()

	m.SetHeaders(headers)

	m.SetBody("text/html", string(msg))

	return c.d.DialAndSend(m)
}
