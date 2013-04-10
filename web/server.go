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
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 22
	server.TLSConfig = nil
	return server
}

func (self *WebServer) Serve(host string, handler http.Handler) {
	self.Addr = host
	self.Handler = handler
	err := self.ListenAndServe()
	if err != nil {
		log.Fatal(fmt.Sprintf("%v : %s\n", time.Now(), err))
	}
}
