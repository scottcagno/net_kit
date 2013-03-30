// ------------
// utilities.go ::: misc helpers and utilities
// ------------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// hex encoder
func EncodeHex(s string) string {
	return hex.EncodeToString([]byte(s))
}

// hex decoder
func DecodeHex(s string) string {
	dat, _ := hex.DecodeString(s)
	return string(dat)
}

// json encoder
func EncodeJSON(v interface{}) string {
	dat, _ := json.Marshal(v)
	return string(dat)
}

// json decoder
func DecodeJSON(s string, v interface{}) {
	json.Unmarshal([]byte(s), &v)
}

// json marshaler
func MarshalJSON(v, m interface{}) {
	dat, _ := json.Marshal(v)
	json.Unmarshal(dat, &m)
}

// url encode a string
func EncodeURL(s string) string {
	encoder := base64.URLEncoding
	encoded := make([]byte, encoder.EncodedLen(len([]byte(s))))
	encoder.Encode(encoded, []byte(s))
	return string(encoded)
}

// url decode a string
func DecodeURL(s string) string {
	encoder := base64.URLEncoding
	decoded := make([]byte, encoder.EncodedLen(len([]byte(s))))
	_, err := encoder.Decode(decoded, []byte(s))
	if err != nil {
		return fmt.Sprintln(err)
	}
	return string(decoded)
}

// random string generator, n number of bytes
func Random(n int) string {
	e := make([]byte, n)
	rand.Read(e)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(e)))
	base64.URLEncoding.Encode(b, e)
	return strings.Replace(string(b), "=", "", -1)
}

// snake a string
func Snake(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "_", -1)
}
