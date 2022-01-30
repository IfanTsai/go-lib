//go:build jsoniter

package json

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	_json = jsoniter.ConfigCompatibleWithStandardLibrary

	Marshal       = _json.Marshal
	MarshalIndent = _json.MarshalIndent
	Unmarshal     = _json.Unmarshal
)
