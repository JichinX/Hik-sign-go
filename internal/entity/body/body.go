package body

import "encoding/json"

type Values map[string]interface{}

func (val Values) Add(k string, v interface{}) {
	val[k] = v
}

func (val Values) String() string {
	b, _ := json.Marshal(val)
	return string(b)
}
