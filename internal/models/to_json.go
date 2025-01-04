// internal/models/base.go
package models

import (
	"encoding/json"
	"fmt"
)

// type JSONConverter[T any] interface {
//     ToJSON() map[string]interface{}
// }

func ToJSON[T any](v T) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	jsonData, err := json.Marshal(v)
	if err != nil {
		return map[string]interface{}{}
	}
	json.Unmarshal(jsonData, &jsonMap)

	fmt.Println(string(jsonData))

	return jsonMap
}
