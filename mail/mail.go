// -------
// mail.go ::: simple smtp wrapper
// -------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package mail

import (
	"bytes"
	"log"
	"net/smtp"
)

type Email struct {
	MailFrom    string
	MailReplyTo string
	MailRcptTo  []string
	MailSubject string
	MailBody    string
}

func NewEmail(from string, reply ...string) *Email {
	var replyTo string
	if len(reply) > 0 {
		replyTo = reply[0]
	} else {
		replyTo = from
	}
	return &Email{
		MailFrom:    from,
		MailReplyTo: replyTo,
	}
}

func (self *Email) To(to ...string) *Email {
	self.MailRcptTo = to
	return self
}

func (self *Email) Subject(subject string) *Email {
	self.MailSubject = subject
	return self
}

func (self *Email) Body(body string) *Email {
	self.MailBody = body
	return self
}

func (self *Email) SendFrom(host string) {
	c, err := smtp.Dial(host + ":25")
	defer c.Quit()
	if err != nil {
		c.Reset()
		log.Fatal(err)
	}
	c.Mail(self.MailFrom)
	for _, to := range self.MailRcptTo {
		c.Rcpt(to)
	}
	wc, err := c.Data()
	if err != nil {
		c.Reset()
		log.Fatal(err)
	}
	defer wc.Close()
	buf := bytes.NewBufferString("Subject: " + self.MailSubject + "\n" + self.MailBody)
	if _, err = buf.WriteTo(wc); err != nil {
		c.Reset()
		log.Fatal(err)
	}
}
