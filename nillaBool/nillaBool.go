package nillaBool

import "encoding/json"

type NillaBool struct {
	value bool
	valid bool
}

func (n *NillaBool) IsBoolFlag() bool {
	return true
}

func (n *NillaBool) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.valid = false
		return nil
	}

	n.valid = true
	return json.Unmarshal(data, &n.value)
}

func (n NillaBool) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.value)
}

func (n NillaBool) String() string {
	if !n.valid {
		return "null"
	}

	if n.value {
		return "true"
	}

	return "false"
}

func (n *NillaBool) Set(value string) error {
	n.valid = true
	n.value = value == "true"
	return nil
}

func (n NillaBool) Get() (bool, bool) {
	return n.value, !n.valid
}

func (n *NillaBool) SetNull() {
	n.valid = false
}

func (n NillaBool) IsNull() bool {
	return !n.valid
}
