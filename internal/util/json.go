package util

import "encoding/json"

func MustMarshal(data map[string]interface{}) string {
	if data == nil {
		data = map[string]interface{}{}
	}
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

func MustUnmarshal(data string) map[string]interface{} {
	if data == "" {
		return map[string]interface{}{}
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return map[string]interface{}{}
	}
	return result
}
