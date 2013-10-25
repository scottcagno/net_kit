// *
// * -----------
// * forms.go ::: html form generator and validator
// * -----------
// * Copyright (c) 2013-Present, Scott Cagno & Greg Pechiro. All rights reserved.
// * BSD-style license can be found at https://sites.google.com/site/bsdc3license
// *

package html

import (
	"html/template"
	"net/http"
	"strings"
	"bytes"
	"fmt"
)

// form structure
type Form struct {
	Header		string
	Action	 	string
	Inputs 		[]*Input
	Submit		*Button
	Errors		map[string]string 
	t 			*template.Template
	values		[]string
}

// initialize form struct
func InitForm(action, header string, submit *Button) *Form {
	return &Form{
		Action	: action,
		Header  : header,
		Submit 	: submit,
		t 		: template.Must(template.New("form").Parse(DEFAULT_FORM)),
		Errors  : make(map[string]string),
	}
}

// add input to form struct
func (self *Form) Add(i *Input) {
	self.Inputs = append(self.Inputs, i)
}

// populate form values
func (self *Form) SetVals(ss []string) {
	if len(ss) != len(self.Inputs) {
		return
	}
	for i, _ := range self.Inputs {
		self.Inputs[i].Value = ss[i]
	}
}

//
func (self *Form) GetVals() []string {
	for i, _ := range self.Inputs {
		self.values = append(self.values, self.Inputs[i].Value)
	}
	return self.values
}

/*
	if !form.IsValid(r) {
		ts.Render(w, "example.html", html.M{"form":form.Render()})
	} else {
		// form is valid, db call for user
		...
		if ok {
			// continue
		} else {
			// session is dead
			http.Redirect(w, r, "/example", 303)
		}
	}
*/

// form validator
func (self *Form) IsValid(r *http.Request) bool {
	validState := true
	var match []string
	for i := 0; i < len(self.Inputs); i++ {
		formVal := r.FormValue(self.Inputs[i].Name)
		if self.Inputs[i].Type != "password" {
			formVal = clean(formVal)
		} else {
			match = append(match, formVal)
		}
		if self.Inputs[i].Required {
			if len(formVal) < 1 {
				err := fmt.Sprintf("*%s is a required field!", self.Inputs[i].Holder)
				self.Inputs[i].Err = err
				self.Errors[self.Inputs[i].Name] = err
				validState = false
			} else {
				self.Inputs[i].Value = formVal
			}
		}
		switch self.Inputs[i].Type {
		case "email":
			if strings.Count(formVal, "@") != 1 {
				err := fmt.Sprintf("*%s requires an email address!", self.Inputs[i].Holder)
				self.Inputs[i].Err = err
				self.Errors[self.Inputs[i].Name] = err
				validState = false
			} else {
				self.Inputs[i].Value = formVal
			}
		case "password":
			if strings.ContainsAny(formVal, ";, \n\r\t") {
				err := fmt.Sprintf("*%s contains invalid characters!", self.Inputs[i].Holder)
				self.Inputs[i].Err = err
				self.Errors[self.Inputs[i].Name] = err
				validState = false
			} else {
				if self.Inputs[i].Min != 0 && len(formVal) < self.Inputs[i].Min {
					err := fmt.Sprintf("*%s minimum is %d!", self.Inputs[i].Holder, self.Inputs[i].Min)
					self.Inputs[i].Err = err
					self.Errors[self.Inputs[i].Name] = err
					validState = false
				} else {
					self.Inputs[i].Value = formVal
				}
				if self.Inputs[i].Max != 0 && self.Inputs[i].Max > self.Inputs[i].Min && len(formVal) > self.Inputs[i].Max {
					err := fmt.Sprintf("*%s maximum is %d!", self.Inputs[i].Holder, self.Inputs[i].Max)
					self.Inputs[i].Err = err
					self.Errors[self.Inputs[i].Name] = err
					validState = false
				} else {
					self.Inputs[i].Value = formVal
				}
			}
		case "number":
			if !isNumber(formVal) {
				err := fmt.Sprintf("*%s requires a number!", self.Inputs[i].Holder)
				self.Inputs[i].Err = err
				self.Errors[self.Inputs[i].Name] = err
				validState = false
			} else {
				self.Inputs[i].Value = formVal
			}
		}
	}
	if len(match) == 2 && match[0] != match[1] {
		validState = false
		for i, _ := range self.Inputs {
			if self.Inputs[i].Type == "password" {
				self.Inputs[i].Err = "*Passwords do not match"
				self.Errors[self.Inputs[i].Name] = "*Passwords do not match"
			} 
		}
	}
	return validState
}

// render form
func (self *Form) Render() string {
	var html bytes.Buffer
	self.t.Execute(&html, self)
	return html.String()
}

// sanitize
func clean(s string) string {
	return strings.Replace(strings.Replace(strings.Join(strings.Fields(s), " "), ",","", -1), ";","", -1)
}

// check is string is a number
func isNumber(s string) bool {
	ok := true
	for _, b := range []byte(s) {
		if !('0' <= b && b <= '9') {
			ok = false
		}
	}
	return ok
}

// input structure
type Input struct {
	Name, Type, Value, Holder, Err, Class 	string
	Required								bool
	Max, Min								int
}

// text input
func TextInput(name, holder string) *Input {
	return InitInput("text", name, holder)
}

// number input
func NumberInput(name, holder string) *Input {
	return InitInput("number", name, holder)
}
 
// email input
func EmailInput(name, holder string) *Input {
	return InitInput("email", name, holder)
}
// password input
func PassInput(name, holder string) *Input {
	return InitInput("password", name, holder)
}
// date input
func DateInput(name, holder string) *Input {
	return InitInput("date", name, holder)
}
// hidden input
func HiddenInput(name, holder string) *Input {
	return InitInput("hidden", name, holder)
}

// initialize input struct
func InitInput(typ, name, holder string) *Input {
	return &Input{
		Type 		: typ,
		Name 		: name,
		Holder 		: holder,
	}
}

// class
func (self *Input) SetClass(c string) *Input {
	self.Class = c
	return self
}

// max
func (self *Input) SetMax(n int) *Input {
	self.Max = n
	return self
}

// min
func (self *Input) SetMin(n int) *Input {
	self.Min = n
	return self
}

// required
func (self *Input) SetRequired() *Input {
	self.Required = true
	return self
}

// button structure
type Button struct {
	Icon, Color, Name 	string
}

// initialize button struct
func InitButton(icon, color, name string) *Button {
	return &Button{
		Icon 	: icon,
		Color 	: color,
		Name	: name,
	}
}

const (
	WHITE		= "btn-default"
	BLUE 		= "btn-primary"
	GREEN 		= "btn-success"
	LIGHT_BLUE 	= "btn-info"
	ORANGE 		= "btn-warning"
	RED 		= "btn-danger"
	LINK 		= "btn-link"
)

// default template
var DEFAULT *template.Template
var DEFAULT_FORM=`<div class="col-sm-12">
	<form class="form form-horizontal form-tight" role="form" method="post" action="{{.Action}}">
		{{if .Header}}<legend class="col-sm-12">{{.Header}}</legend>{{end}}
		{{range .Inputs}}
		<div class="form-group line">
			<div class="col-sm-12">
				{{if .Err}}<label class="text-danger text-left">{{.Err}}</label>{{end}}
				<input name="{{.Name}}" type="{{.Type}}" placeholder="{{.Holder}}" {{if .Value}}value="{{.Value}}"{{end}} class="form-control input-sm{{if .Class}}{{.Class}}{{end}}" {{if .Min}}min="{{.Min}}"{{end}} {{if .Max}}max="{{.Max}}"{{end}} {{if .Required}}required{{end}}>
			</div>
		</div>
		{{end}}
		<div class="form-group">
			<div class="col-sm-12">
				<button type="submit" class="btn {{.Submit.Color}}">
					<i class="{{.Submit.Icon}}"></i> {{.Submit.Name}}
				</button>
			</div>
		</div>
	</form>
</div>`