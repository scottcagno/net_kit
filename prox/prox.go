// -------
// prox.go ::: reverse proxy server
// -------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"sync"
)

const HOST = ":80"

func main() {
	revprox := NewReverseProxy()
	revprox.Register("localhost/one", "localhost:8081")
	revprox.Register("localhost/two", "localhost:8081")
	revprox.Register("localhost/three", "localhost:8081")
	fd, _ := strconv.Atoi(os.Getenv("RUNSIT_PORTFD_http"))
	log.Fatal(http.Serve(listen(fd, HOST), revprox))
}

func listen(fd int, addr string) net.Listener {
	var l net.Listener
	var err error
	if fd >= 3 {
		l, err = net.FileListener(os.NewFile(uintptr(fd), "http"))
	} else {
		l, err = net.Listen("tcp", addr)
	}
	if err != nil {
		log.Fatal(err)
	}
	return l
}

type ReverseProxy struct {
	services map[string]http.Handler
	mu       sync.RWMutex
}

func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{
		services: make(map[string]http.Handler),
	}
}

func (self *ReverseProxy) Register(host, proxyto string) *ReverseProxy {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.services[host] = self.proxyHandler(proxyto)
	return self
}

func (self *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if proxHdl := self.handleProxReq(r); proxHdl != nil {
		proxHdl.ServeHTTP(w, r)
		return
	}
	fmt.Fprintf(w, "%d %s", 404, http.StatusText(404))
}

func (self *ReverseProxy) handleProxReq(r *http.Request) http.Handler {
	self.mu.Lock()
	defer self.mu.Unlock()
	host := r.Host
	if i := strings.Index(host, ":"); i >= 0 {
		host = host[:i]
	}
	if proxHdl, ok := self.services[host]; ok {
		return proxHdl
	}
	return nil
}

func (self *ReverseProxy) proxyHandler(s string) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = s
		},
	}
}
