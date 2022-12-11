package encoding

import "encoding/json"

func JsonMarshal(payload interface{}) []byte {
	data, _ := json.Marshal(payload)
	return data
}

func JsonMarshalString(payload interface{}) string {
	return string(JsonMarshal(payload))
}

func JsonUnMarshalString(payload string, v interface{}) error {
	return json.Unmarshal([]byte(payload), v)
}
