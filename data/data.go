// -------
// data.go ::: data wrapper
// -------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package data

// data wrapper interface
type DataWrapper interface {
	Insert(v ...interface{}) interface{}
	Update(v ...interface{}) interface{}
	Return(v ...interface{}) interface{}
	Delete(v ...interface{}) interface{}
}
