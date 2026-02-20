package smtp

import (
	"bytes"
	"io"
	"net/mail"
	"strings"

	"github.com/emersion/go-smtp"
	"sink.io/m/src/store"
	"sink.io/m/src/ws"
)

type Backend struct {
	Store *store.Store
	Hub   *ws.Hub
}

type Session struct {
	backend *Backend
	from    string
	to      string
}

func (b *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{backend: b}, nil
}

func (s *Session) AuthPlain(_, _ string) error { return nil }

func (s *Session) Mail(from string, _ *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, _ *smtp.RcptOptions) error {
	s.to = strings.ToLower(to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	msg, err := mail.ReadMessage(bytes.NewReader(raw))
	if err != nil {
		return err
	}

	subject := msg.Header.Get("Subject")
	bodyBytes, _ := io.ReadAll(msg.Body)

	email := s.backend.Store.Add(store.Email{
		To:      s.to,
		From:    s.from,
		Subject: subject,
		Body:    string(bodyBytes),
	})

	s.backend.Hub.Broadcast(s.to, email)
	return nil
}

func (s *Session) Reset()        {}
func (s *Session) Logout() error { return nil }
