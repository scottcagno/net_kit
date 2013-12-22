// -----------
// template.go ::: html template store
// -----------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package tmpl

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"
)

// map type
type M map[string]interface{}

// template
type TemplateStore struct {
	dir    string
	base   string
	cached map[string]*template.Template
	funcs  template.FuncMap
	mu     sync.Mutex
}

// return a new template store instace
func NewTemplateStore(dir, base string) *TemplateStore {
	return &TemplateStore{
		dir:    dir,
		base:   base,
		cached: make(map[string]*template.Template),
		funcs: template.FuncMap{
			"title": strings.Title,
			"safe":  safe,
			"add":   add,
			"sub":   sub,
			"decr":  decr,
			"incr":  incr,
			"split": strings.Split,
		},
	}
}

// html safe escaper
func safe(html string) template.HTML {
	return template.HTML(html)
}

// decrement
func decr(a int) int { return a - 1 }

// increment
func incr(a int) int { return a + 1 }

// add
func add(a, b int) int { return a + b }

// subtract
func sub(a, b int) int { return a - b }

// load template files associated with base into cache
func (self *TemplateStore) Load(name ...string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	for i := 0; i < len(name); i++ {
		t := template.New(self.base).Funcs(self.funcs)
		t = template.Must(t.ParseFiles(self.dir+"/"+self.base, self.dir+"/"+name[i]))
		self.cached[name[i]] = t
	}
}

// register raw html strings into cache (must supply keys)
func (self *TemplateStore) Register(m map[string]string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	for k, v := range m {
		self.cached[k] = template.Must(template.New("tmpl").Funcs(self.funcs).Parse(v))
	}
}

// render a template by name
func (self *TemplateStore) Render(w http.ResponseWriter, name string, m interface{}) {
	self.cached[name].Execute(w, m)
}

// render raw data
func (self *TemplateStore) Raw(w http.ResponseWriter, format string, a ...interface{}) {
	fmt.Fprintf(w, format, a...)
}

// set the header content type
func (self *TemplateStore) ContentType(w http.ResponseWriter, typ string) {
	w.Header().Set("Content Type", typ)
}

/*
// simple form validater
func (self *TemplateStore) Valid(w http.ResponseWriter, v interface{}) (map[string]string, bool) {
	m, ok := map[string]string{"errors": "error"}, true
	for k, s := range v.(M) {
		s = strings.TrimSpace(s.(string))
		m[k] = s.(string)
		switch {
		case s.(string) == "":
			m["errors"] = m["errors"] + ", " + k + " required"
			ok = false
		case strings.Contains(s.(string), ","):
			m[k] = strings.Replace(s.(string), ",", "", -1)
		case k == "email":
			if strings.Count(s.(string), "@") != 1 {
				m["errors"] = m["errors"] + " invalid email"
				ok = false
				break
			}
		case k == "pass":
			if len(s.(string)) < 6 {
				m["errors"] = m["errors"] + " min length 6"
				ok = false
				break

			}
		case k == "confirm":
			if s.(string) != v.(M)["pass"].(string) {
				m["errors"] = m["errors"] + " pass does not match"
				ok = false
				break
			}
		}
	}

		if _, ok := v.(M)["confirm"]; ok {
			if v.(M)["confirm"].(string) != v.(M)["pass"].(string) {
				m["errors"] = m["errors"]+", pass does not match"
				ok = false
			}
		}

	return m, ok
}
*/
