package data

import ()

type TestClient struct {
	m map[string]interface{}
}

type M map[string]interface{}

func InitTestClient(m *map[string]interface{}) *TestClient {
	return &TestClient{
		m: *m,
	}
}

func (self *TestClient) Open() bool {
	return true
}

func (self *TestClient) Close() bool {
	return true
}

func (self *TestClient) Has(k string) bool {
	_, ok := self.m[k]
	return ok
}

func (self *TestClient) Set(k string, v interface{}) bool {
	self.m[k] = v
	_, ok := self.m[k]
	return ok
}

func (self *TestClient) Get(k string) (interface{}, bool) {
	v, ok := self.m[k]
	if ok {
		return v, ok
	}
	return nil, ok
}

func (self *TestClient) Del(k string) bool {
	delete(self.m, k)
	_, ok := self.m[k]
	return ok
}

func (self *TestClient) Fnd(q map[string]interface{}) (interface{}, bool) {
	return mapsubset(self.m, q)
}

func mapsubset(m, subm map[string]interface{}) (interface{}, bool) {
	var value interface{}
	for k, v := range subm {
		if val, ok := m[k]; !ok || !subset(val, v) {
			return nil, false
		} else {
			value = val
		}
	}
	return value, true
}

func subset(v, subv interface{}) bool {
	return v == subv
}
