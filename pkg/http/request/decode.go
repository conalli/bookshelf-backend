package request

import (
	"encoding/json"
	"io"
)

// DecodeJSONRequest takes in a JSON APIRequest and attempts to decode it.
func DecodeJSONRequest[T APIRequest](body io.ReadCloser) (T, error) {
	var data T
	err := json.NewDecoder(body).Decode(&data)
	return data, err
}
