// --------
// email.go ::: smtp wrapper
// --------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package mail

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
)

const AUTH_KEY = "038e376187f47b718b9fac83dab476d9ecfb7f3f4955f96135f571c7b9324ba2c3395b"

// email structure
type Email struct {
	Host_, From_, Reply_, Subject_, To_, Body_ string
}

// get new email instance, supply from and string
func NewEmail(from, host string) *Email {
	return &Email{
		Host_:  host,
		From_:  from,
		Reply_: from,
	}
}

// set subject
func (self *Email) Subject(subject string) *Email {
	self.Subject_ = subject
	return self
}

// set to
func (self *Email) To(to string) *Email {
	self.To_ = to
	return self
}

// set body
func (self *Email) Body(body string) *Email {
	self.Body_ = body
	return self
}

// send mail
func (self *Email) SendMail() {
	c, err := smtp.Dial(self.Host_)
	if err != nil {
		c.Reset()
		log.Fatal(err)
	}
	c.Mail(self.From_)
	c.Rcpt(self.To_)
	wc, err := c.Data()
	if err != nil {
		c.Reset()
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString("To: " + self.To_ + "\nSubject: " + self.Subject_ + "\n" + self.Body_)
	if _, err = buf.WriteTo(wc); err != nil {
		c.Reset()
		log.Fatal(err)
	}
}

// encode body
func EncodeBody(body string) string {
	return url.QueryEscape(body)
}

// decode body
func DecodeBody(body string) string {
	dat, _ := url.QueryUnescape(body)
	return dat
}

// http api hook
func HttpApiHook(w http.ResponseWriter, r *http.Request) {
	auth := r.FormValue("auth")
	if len(auth) <= 0 || auth != AUTH_KEY {
		http.Redirect(w, r, "/error/405", 303)
		return
	}
	host := r.FormValue("host")
	from := r.FormValue("from")
	repl := r.FormValue("reply")
	sub := r.FormValue("subject")
	to := r.FormValue("to")
	body := r.FormValue("body")
	if len(host) <= 0 || len(from) <= 0 || len(repl) <= 0 || len(sub) <= 0 || len(to) <= 0 || len(body) <= 0 {
		http.Redirect(w, r, "/error/405", 303)
		return
	}
	email := &Email{
		Host_:    host,
		From_:    from,
		Reply_:   repl,
		Subject_: sub,
		To_:      to,
		Body_:    DecodeBody(body),
	}
	email.SendMail()
	fmt.Fprintln(w, "got it, thanks.")
}
