//go:build !jsoniter

package json

import "encoding/json"

var (
	MarshalIndent = json.MarshalIndent
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
)
