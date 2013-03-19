// --------------
// multiplexer.go ::: http multiplexer
// --------------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package web

import (
	"net/http"
	"net/url"
	"strings"
)

// http multiplexer
type Multiplexer struct {
	handlers map[string][]*Handler
}

// return new multiplexer instance
func NewMultiplexer() *Multiplexer {
	return &Multiplexer{
		handlers: make(map[string][]*Handler),
	}
}

// match request against registered handlers, server http
func (self *Multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range self.handlers[r.Method] {
		if params, ok := h.parse(r.URL.Path); ok {
			if len(params) > 0 {
				r.URL.RawQuery = url.Values(params).Encode() + "&" + r.URL.RawQuery
			}
			h.ServeHTTP(w, r)
			return
		}
	}
	allowed := make([]string, 0, len(self.handlers))
	for method, handlers := range self.handlers {
		if method == r.Method {
			continue
		}
		for _, h := range handlers {
			if _, ok := h.parse(r.URL.Path); ok {
				allowed = append(allowed, method)
			}
		}
	}
	if len(allowed) == 0 {
		http.Redirect(w, r, "/error/404", 303)
		return
	}
	w.Header().Add("Allow", strings.Join(allowed, ", "))
	http.Redirect(w, r, "/error/405", 303)
}

// register an http hander for a particular method and path
func (self *Multiplexer) Handle(method, path string, h http.Handler) {
	self.handlers[method] = append(self.handlers[method], &Handler{path, h})
	n := len(path)
	if n > 0 && path[n-1] == '/' {
		self.Handle(method, path[:n-1], http.RedirectHandler(path, 301))
	}
}

// wrapper for Handle to allow use of handler functions
func (self *Multiplexer) HandleFunc(method, path string, h http.HandlerFunc) {
	self.Handle(method, path, h)
}

// register an http handler func for the get method
func (self *Multiplexer) Get(path string, h http.HandlerFunc) {
	self.Handle("GET", path, h)
}

// register an http handler func for the post method
func (self *Multiplexer) Post(path string, h http.HandlerFunc) {
	self.Handle("POST", path, h)
}

// register an http handler func for the put method
func (self *Multiplexer) Put(path string, h http.HandlerFunc) {
	self.Handle("PUT", path, h)
}

// register an http handler func for the delete method
func (self *Multiplexer) Delete(path string, h http.HandlerFunc) {
	self.Handle("DELETE", path, h)
}

// register forwarder handler
func (self *Multiplexer) Forward(path, newpath string) {
	self.Handle("GET", path, http.RedirectHandler(newpath, 301))
}

func (self *Multiplexer) Static(path, folder string) {
	n := len(path)
	if n > 0 && path[n-1] != '/' {
		path = path + "/"
	}
	h := http.StripPrefix(path, http.FileServer(http.Dir(folder)))
	self.Handle("GET", path, h)
}

// handler
type Handler struct {
	path string
	http.Handler
}

// parse registered pattern
func (self *Handler) parse(path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(self.path):
			if self.path != "/" && len(self.path) > 0 && self.path[len(self.path)-1] == '/' {
				return p, true
			}
			return nil, false
		case self.path[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(self.path, isBoth, j+1)
			val, _, i = match(path, byteParse(nextc), i)
			p.Add(":"+name, val)
		case path[i] == self.path[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(self.path) {
		return nil, false
	}
	return p, true
}

// match path with registered handler
func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

// determine type of byte
func byteParse(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

// test for alpha byte
func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// test for numerical byte
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// test for alpha or numerical byte
func isBoth(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
