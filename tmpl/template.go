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
			"safe": safe,
		},
	}
}

// html safe escaper
func safe(html string) template.HTML {
	return template.HTML(html)
}

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
