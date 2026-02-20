package smtp

import (
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

func StartSMTP(b *Backend, addr string) {
	srv := smtp.NewServer(b)
	srv.Addr = addr
	srv.Domain = "sink.io.local"
	srv.ReadTimeout = 10 * time.Second
	srv.AllowInsecureAuth = true

	log.Printf("SMTP listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
