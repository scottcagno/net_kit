// --------
// forms.go ::: html form generator and validator
// --------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package frms

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

/*
type MultiForm struct {
	Action, Button, ButtonName string
	Forms                      map[string]*template.Template
}
*/

type Form struct {
	Inputs []Input
	Errors map[string]string
}

type Input struct {
	Name, Type, Value, Holder, Error, Class string
	Required                                bool
	Max, Min                                int
}

func (self *Form) Render(t *template.Template) string {
	var html bytes.Buffer
	t.Execute(&html, self)
	return html.String()
}

func (self *Form) SetError(inputName, errStr string) {
	for i := 0; i < len(self.Inputs); i++ {
		if self.Inputs[i].Name == inputName {
			self.Inputs[i].Error = errStr
			self.Errors[self.Inputs[i].Name] = errStr
		}
	}
}

func (self *Form) IsValid(r *http.Request) bool {
	isValid := true
	for i := 0; i < len(self.Inputs); i++ {
		formVal := r.FormValue(self.Inputs[i].Name)
		if self.Inputs[i].Required {
			//formVal := r.FormValue(self.Inputs[i].Name)
			if len(formVal) < 1 {
				err := fmt.Sprintf("*%s is a required field", strings.Title(self.Inputs[i].Name))
				self.Inputs[i].Error = err
				self.Errors[self.Inputs[i].Name] = err
				isValid = false
				//break
			} else {
				self.Inputs[i].Value = formVal
			}
			switch self.Inputs[i].Type {
			case "email":
				if !strings.Contains(formVal, "@") {
					err := "*Email field requires an email address"
					self.Inputs[i].Error = err
					self.Errors[self.Inputs[i].Name] = err
					isValid = false
					//break
				} else {
					self.Inputs[i].Value = formVal
				}
			case "password":
				if self.Inputs[i].Min != 0 && len(formVal) < self.Inputs[i].Min {
					err := fmt.Sprintf("*Password minimum is %d", self.Inputs[i].Min)
					self.Inputs[i].Error = err
					self.Errors[self.Inputs[i].Name] = err
					isValid = false
					//break
				} else {
					self.Inputs[i].Value = formVal
				}
				if self.Inputs[i].Max != 0 && self.Inputs[i].Max > self.Inputs[i].Min && len(formVal) > self.Inputs[i].Max {
					err := fmt.Sprintf("*Password maximum is %d", self.Inputs[i].Max)
					self.Inputs[i].Error = err
					self.Errors[self.Inputs[i].Name] = err
					isValid = false
					//break
				} else {
					self.Inputs[i].Value = formVal
				}
			case "number":
				if !isNumber(formVal) {
					err := fmt.Sprintf("*%s requires a number", strings.Title(self.Inputs[i].Name))
					self.Inputs[i].Error = err
					self.Errors[self.Inputs[i].Name] = err
					isValid = false
					//break
				} else {
					self.Inputs[i].Value = formVal
				}
			}
		}
	}
	return isValid
}

func isNumber(s string) bool {
	isNum := true
	for _, b := range []byte(s) {
		if !('0' <= b && b <= '9') {
			isNum = false
		}
	}
	return isNum
}

func init() {
	var FUNCS = template.FuncMap{
		"safe": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
	}
	DEFAULT = template.Must(template.New("form").Funcs(FUNCS).Parse(DEFAULT_FORM))
	//INLINE = template.Must(template.New("form").Funcs(FUNCS).Parse(INLINE_FORM))
	//PARTIAL = template.Must(template.New("form").Funcs(FUNCS).Parse(PARTIAL_FORM))
}

var DEFAULT *template.Template
var DEFAULT_FORM = `{{range .Inputs}}
    	<div class="form-group line">
    		<div class="col-sm-12">
        		{{ if .Error }}<label class="text-danger text-left">{{ .Error }}</label>{{ end }}
        		<input class="form-control{{if .Class}} {{.Class}}{{end}}" type="{{.Type}}" name="{{.Name}}" {{ if .Value }}value="{{ .Value }}"{{ end }} placeholder="{{ .Holder }}" {{if .Min}}min="{{.Min}}"{{end}} {{if .Max}}max="{{.Max}}"{{end}} {{if .Required}}required{{end}}>
    		</div>
    	</div>
    {{end}}`

/*
var old_DEFAULT *template.Template
var old_DEFAULT_FORM = `<form {{ if .Id }}id="{{ .Id }}"{{ end }} method="post" action="{{ .Action }}" class="form form-horizontal" role="form">
    {{range .Inputs}}
    	<div class="form-group line">
    		<div class="col-sm-12">
        		{{ if .Error }}<label class="text-danger text-left">{{ .Error }}</label>{{ end }}
        		<input class="form-control{{if .Class}} {{.Class}}{{end}}" type="{{.Type}}" name="{{.Name}}" {{ if .Value }}value="{{ .Value }}"{{ end }} placeholder="{{ .Holder }}" {{if .Min}}min="{{.Min}}"{{end}} {{if .Max}}max="{{.Max}}"{{end}} {{if .Required}}required{{end}}>
    		</div>
    	</div>
    {{end}}
    	<div class="form-group">
    		<div class="col-sm-12">
    			<button {{ if .ButtonName }}name="{{ .ButtonName }}"{{ end }} type="submit" class="btn btn-zoom">{{ .Button }}</button>
    		</div>
    	</div>
    </form>`

var INLINE *template.Template
var INLINE_FORM = `<form role="form" method="post" action="{{ .Action }}" class="form-inline">
<p class="text-danger">{{ if .Errors }}{{ range .Inputs }}{{ .Error }}<br/>{{ end }}{{ end }}</p>
{{range .Inputs}}
    <input class="input{{if .Class}} {{.Class}}{{end}}" type="{{.Type}}" name="{{.Name}}" value="{{.Value}}" placeholder="{{ .Holder }}" {{if .Min}}min="{{.Min}}"{{end}} {{if .Max}}max="{{.Max}}"{{end}} {{if .Required}}required{{end}}>
{{end}}
    <button {{ if .ButtonName }}name="{{ .ButtonName }}"{{ end }} type="submit" class="btn">{{ .Button }}</button>
</form>`

var PARTIAL *template.Template
var PARTIAL_FORM = `{{range .Inputs}}
    	<div class="form-group">
        	<label class="text-danger text-left">{{ .Error }}</label>
        	<input class="input{{if .Class}} {{.Class}}{{end}}" type="{{.Type}}" name="{{.Name}}" {{ if .Value }}value="{{ .Value }}"{{ end }} placeholder="{{ .Holder }}" {{if .Min}}min="{{.Min}}"{{end}} {{if .Max}}max="{{.Max}}"{{end}} {{if .Required}}required{{end}}>
    	</div>
    {{end}}`
*/
