// ----------
// session.go ::: session manager
// ----------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package sess

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	MIN   = 60
	HOUR  = MIN * 60
	DAY   = HOUR * 24
	WEEK  = DAY * 7
	MONTH = DAY * 30
	YEAR  = WEEK * 52
)

type Store struct {
	cookieId string
	rate     int64
	sessions map[string]*Session
	mu       sync.Mutex
}

func NewSessionStore(cookieId string, rate int64) *Store {
	store := &Store{
		cookieId: cookieId,
		rate:     rate,
		sessions: make(map[string]*Session, 0),
	}
	store.GC()
	return store
}

func (self *Store) FreshCookie(sid string) http.Cookie {
	return http.Cookie{
		Name:     self.cookieId,
		Value:    url.QueryEscape(sid),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(self.rate),
	}
}

func (self *Store) FreshSession(sid string) *Session {
	return &Session{
		sid:   sid,
		store: self,
		ts:    time.Now(),
		vals:  make(map[string]string, 0),
	}
}

func (self *Store) NewSession(w http.ResponseWriter, r *http.Request) *Session {
	self.mu.Lock()
	defer self.mu.Unlock()
	sid := Random(32)
	session := self.FreshSession(sid)
	self.sessions[sid] = session
	cookie := self.FreshCookie(sid)
	http.SetCookie(w, &cookie)
	return session
}

func (self *Store) GetSession(w http.ResponseWriter, r *http.Request) *Session {
	self.mu.Lock()
	defer self.mu.Unlock()
	var session *Session
	cookie, err := r.Cookie(self.cookieId)
	if err != nil || cookie.Value == "" {
		sid := Random(32)
		session = self.FreshSession(sid)
		self.sessions[sid] = session
		cookie := self.FreshCookie(sid)
		http.SetCookie(w, &cookie)
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session = self.sessions[sid]
	}
	return session
}

func (self *Store) DelSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(self.cookieId)
	if err != nil && err == http.ErrNoCookie || cookie.Value == "" {
		return
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	sid, _ := url.QueryUnescape(cookie.Value)
	delete(self.sessions, sid)
	*cookie = self.FreshCookie(sid)
	cookie.MaxAge = -1
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
}

func (self *Store) ExtSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(self.cookieId)
	if err != nil && err == http.ErrNoCookie || cookie.Value == "" {
		return
	}
	self.mu.Lock()
	defer self.mu.Unlock()
	sid, _ := url.QueryUnescape(cookie.Value)
	if session, ok := self.sessions[sid]; ok {
		*cookie = self.FreshCookie(sid)
		currentTime := time.Now()
		session.ts = currentTime
		cookie.Expires = currentTime.Add(time.Duration(self.rate) * time.Second)
		http.SetCookie(w, cookie)
	}
}

func (self *Store) Update(sid string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	if session, ok := self.sessions[sid]; ok {
		session.ts = time.Now()
	}
}

func (self *Store) GC() {
	self.mu.Lock()
	defer self.mu.Unlock()
	for sid, session := range self.sessions {
		if (session.ts.Unix() + self.rate) < time.Now().Unix() {
			delete(self.sessions, sid)
		} else {
			break
		}
	}
	time.AfterFunc(time.Duration(self.rate)*time.Second, func() {
		self.GC()
	})
}

func (self *Store) ViewSessions() {
	for k, v := range self.sessions {
		fmt.Printf("key: %v\nval: %v\n\n", k, v)
	}
}

type Session struct {
	sid   string
	store *Store
	ts    time.Time
	vals  map[string]string
}

func (self *Session) SetFlash(style, key, val string) {
	self.vals["flash-"+key] = style + ":" + val
	self.store.Update(self.sid)
}

func (self *Session) GetFlash(key string) map[string]string {
	vals := make([]string, 2)
	flash := make(map[string]string, 0)
	if val, ok := self.vals["flash-"+key]; ok {
		vals = strings.Split(val, ":")
		flash["style"] = vals[0]
		flash["flash"] = vals[1]
		delete(self.vals, "flash-"+key)
		self.store.Update(self.sid)
	} else {
		flash = nil
	}
	return flash
}

func (self *Session) Set(key, val string) {
	self.vals[key] = val
	self.store.Update(self.sid)
}

func (self *Session) Get(key string) string {
	value := ""
	if val, ok := self.vals[key]; ok {
		value = val
		self.store.Update(self.sid)
	}
	return value
}

func (self *Session) Del(key string) {
	delete(self.vals, key)
	self.store.Update(self.sid)
}

func (self *Session) Id() string {
	id := self.sid
	self.store.Update(self.sid)
	return id
}

func Random(n int) string {
	e := make([]byte, n)
	rand.Read(e)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(e)))
	base64.URLEncoding.Encode(b, e)
	return string(b)
}
