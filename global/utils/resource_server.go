package utils

import (
	jsoniter "github.com/json-iterator/go"
)

//UnmarshalMap is used to convert map[string]interface{} to a struct
func UnmarshalMap(data map[string]interface{}, toValue interface{}) {
	json := jsoniter.ConfigFastest
	jsonString, _ := json.Marshal(data)
	json.Unmarshal(jsonString, toValue)
}
