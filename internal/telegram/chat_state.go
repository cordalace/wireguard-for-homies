package telegram

import (
	"encoding/json"
)

type ChatState int

const (
	ChatStateInitial = iota
	ChatStateSubnetExpectCIDR
)

func (s ChatState) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

func ChatStateFromJSON(data []byte) (ChatState, error) {
	var ret ChatState
	if err := json.Unmarshal(data, &ret); err != nil {
		return ChatStateInitial, err
	}
	return ret, nil
}
