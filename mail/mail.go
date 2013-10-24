// *
// * -------
// * mail.go ::: email client
// * -------
// * Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// * License can be found at sites.google.com/site/bsdc3license
// *

package mail

import (
	"net/smtp"
	"log"
)

// email account
type account struct {
	email, pass, host 	string
	auth 				smtp.Auth
}

func InitAccount(email, pass string, h ...string) *account {
	self := &account{email:email,pass:pass}
	if len(h) == 1 
{		self.host = h[0]
	}
}

func (self *account) Send(to)