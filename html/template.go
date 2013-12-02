// *
// * ---------------
// * template.go ::: template store
// * ---------------
// * Copyright (c) 2013-Present, Scott Cagno & Greg Pechiro. All rights reserved.
// * BSD-style license can be found at https://sites.google.com/site/bsdc3license
// *

package html

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"sort"
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
func InitTemplateStore(dir, base string) *TemplateStore {
	return &TemplateStore{
		dir:    dir,
		base:   base,
		cached: make(map[string]*template.Template),
		funcs: template.FuncMap{
			"title": strings.Title,
			"safe":  safe,
			"eq":    eq,
			"add":   add,
			"sub":   sub,
			"decr":  decr,
			"incr":  incr,
			"split": strings.Split,
			"sort":  sorter,
			"enc":   Encode,
			"dec":   Decode,
		},
	}
}

func Encode(s string, a ...interface{}) string {
	return base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(fmt.Sprintf(s, a...))))
}

func Decode(s string) string {
	val, _ := base64.StdEncoding.DecodeString(s)
	val2, _ := url.QueryUnescape(string(val))
	return val2
}

// sorts a string slice
func sorter(s []string) []string {
	sort.Strings(s)
	return s
}

// html safe escaper
func safe(html string) template.HTML {
	return template.HTML(html)
}

// check for equality
func eq(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
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
