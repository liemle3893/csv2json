package json

import (
	JSON "encoding/json"
)

// JsonObject is a wrapper around map[string]interface{} to provide some OOP-liked method for JSON.
type JsonObject map[string]interface{}

func NewJsonObject() JsonObject {
	m := make(map[string]interface{})
	return JsonObject(m)
}

func (json JsonObject) GetObject(key string) (JsonObject, bool) {
	val, ok := json[key].(JsonObject)
	return val, ok
}

func (json JsonObject) Get(key string) (interface{}, bool) {
	value, ok := json[key]
	return value, ok
}

func (json JsonObject) Put(key string, value interface{}) {
	json[key] = value
}

func (json JsonObject) String() string {
	s, err := JSON.Marshal(json)
	if err != nil {
		return "{}"
	} else {
		return string(s)
	}
}
