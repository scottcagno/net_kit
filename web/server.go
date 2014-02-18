// ---------
// server.go ::: http server
// ---------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package web

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type WebServer struct {
	http.Server
}

func NewWebServer() *WebServer {
	server := &WebServer{}
	server.ReadTimeout = 120 * time.Second
	server.WriteTimeout = 30 * time.Second
	server.MaxHeaderBytes = 1 << 23
	return server
}

// listen and serve tls
func (self *WebServer) ServeTLS(host, crt, key string, handler http.Handler) {
	self.Addr = host
	self.Handler = handler
	err := self.ListenAndServeTLS(crt, key)
	if err != nil {
		log.Fatal(fmt.Sprintf("%v : %s\n", time.Now(), err))
	}
}

// listen and serve
func (self *WebServer) Serve(host string, handler http.Handler) {
	self.Addr = host
	self.Handler = handler
	err := self.ListenAndServe()
	if err != nil {
		log.Fatal(fmt.Sprintf("%v : %s\n", time.Now(), err))
	}
}
