// ----------
// mgowrap.go ::: mongo data wrapper implementation
// ----------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package data

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

// mgo data wrapper
type MgoWrapper struct {
	Session  *mgo.Session
	Database *mgo.Database
	C        *mgo.Collection
}

// return a new data wrapper instance
func NewMgoWrapper(host string) *MgoWrapper {
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	return &MgoWrapper{
		Session: session,
	}
}

// set database
func (self *MgoWrapper) SetDb(db string) *MgoWrapper {
	self.Database = self.Session.DB(db)
	return self
}

func (self *MgoWrapper) Login(user, pass string) *MgoWrapper {
	err := self.Database.Login(user, pass)
	if err != nil {
		log.Panicf("MgoWrapper - login error: %v\n", err)
	}
	return self
}

// set collection
func (self *MgoWrapper) SetC(c string) *MgoWrapper {
	self.C = self.Database.C(c)
	return self
}

// insert
func (self *MgoWrapper) Insert(v ...interface{}) interface{} {
	err := self.C.Insert(v...)
	if err != nil {
		return err
	}
	return len(v)
}

// update
func (self *MgoWrapper) Update(v ...interface{}) interface{} {
	info, err := self.C.UpdateAll(v[0], v[1])
	if err != nil {
		return err
	}
	return info.Updated
}

// return
func (self *MgoWrapper) Return(v ...interface{}) interface{} {
	var lmt int
	var sel, set, ret interface{}
	for k, val := range v {
		switch val.(type) {
		case int:
			lmt = val.(int)
			v = append(v[:k], v[k+1:]...)
		}
	}
	switch len(v) {
	case 1:
		sel, set = bson.M{}, v[0]
	case 2:
		sel, set = v[0], v[1]
	default:
		ret = nil
	}
	switch lmt {
	case 0:
		ret = self.C.Find(sel).All(set)
	case 1:
		ret = self.C.Find(sel).Sort("-_id").One(set)
	default:
		ret = self.C.Find(sel).Limit(lmt).All(set)
	}
	return ret
}

// return
func (self *MgoWrapper) ReturnSort(sid string, v ...interface{}) interface{} {
	var lmt int
	var sel, set, ret interface{}
	for k, val := range v {
		switch val.(type) {
		case int:
			lmt = val.(int)
			v = append(v[:k], v[k+1:]...)
		}
	}
	switch len(v) {
	case 1:
		sel, set = bson.M{}, v[0]
	case 2:
		sel, set = v[0], v[1]
	default:
		ret = nil
	}
	switch lmt {
	case 0:
		ret = self.C.Find(sel).Sort(sid).All(set)
	case 1:
		ret = self.C.Find(sel).Sort("-_id").One(set)
	default:
		ret = self.C.Find(sel).Sort(sid).Limit(lmt).All(set)
	}
	return ret
}

// delete
func (self *MgoWrapper) Delete(v ...interface{}) interface{} {
	info, err := self.C.RemoveAll(v[0])
	if err != nil {
		return err
	}
	return info.Removed
}
