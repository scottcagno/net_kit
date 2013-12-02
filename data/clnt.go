package data

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
	host string
	r    *bufio.Reader
}

func InitClient(host string) *Client {
	return &Client{
		host: host,
	}
}

func (self *Client) Open() bool {
	if conn, err := net.Dial("tcp", self.host); err == nil {
		self.conn = conn
		self.r = bufio.NewReader(self.conn)
		return true
	}
	return false
}

func (self *Client) Close() bool {
	if err := self.conn.Close(); err != nil {
		return false
	}
	return true
}

func (self *Client) Has(k string) bool {
	fmt.Fprintf(self.conn, "has %v\n", k)
	r := self.response()
	if r == nil || len(r) == 0 {
		return false
	}
	return btobool(r)
}

func (self *Client) Set(k string, v interface{}) bool {
	b, err := json.Marshal(v)
	if err != nil {
		return false
	}
	fmt.Fprintf(self.conn, "set %s %s\n", k, b)
	r := self.response()
	if r == nil || len(r) == 0 {
		return false
	}
	return btobool(r)
}

func (self *Client) Get(k ...string) (interface{}, bool) {
	if len(k) == 0 {
		return nil, false
	} else if len(k) == 1 {
		fmt.Fprintf(self.conn, "get %v\n", k[0])
	} else {
		b, err := json.Marshal(k)
		if err != nil {
			return nil, false
		}
		fmt.Fprintf(self.conn, "get %s\n", b)
	}
	r := self.response()
	if r == nil || len(r) == 0 {
		return nil, false
	}
	if len(r) > 1 {
		var v interface{}
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, false
		}
		return v, true
	}
	return nil, false
}

func (self *Client) Del(k string) bool {
	fmt.Fprintf(self.conn, "del %v\n", k)
	r := self.response()
	if r == nil || len(r) == 0 {
		return false
	}
	return btobool(r)
}

func (self *Client) Fnd(q map[string]interface{}) (interface{}, bool) {
	b, err := json.Marshal(q)
	if err != nil {
		return nil, false
	}
	fmt.Fprintf(self.conn, "fnd %s\n", b)
	r := self.response()
	if r == nil || len(r) == 0 {
		return nil, false
	}
	if len(r) > 1 {
		var v interface{}
		if err := json.Unmarshal(r, &v); err != nil {
			return nil, false
		}
		return v, true
	}
	return nil, false
}

func (self *Client) response() []byte {
	b, err := self.r.ReadBytes('\n')
	if err != nil {
		return nil
	}
	return b[:len(b)-1]
}

func btobool(b []byte) bool {
	if b[0] == '1' {
		return true
	} else if b[0] == '0' {
		return false
	}
	return false
}
