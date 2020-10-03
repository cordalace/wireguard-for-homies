package badgerdb

import "encoding/json"

func (t *badgerTx) GetOrCreateDeviceName(defaultDeviceName string) (string, error) {
	defaultValue, err := json.Marshal(defaultDeviceName)
	if err != nil {
		return "", err
	}

	value, err := t.getOrCreate("deviceName", defaultValue)
	if err != nil {
		return "", err
	}

	var deviceName string
	err = json.Unmarshal(value, &deviceName)
	if err != nil {
		return "", err
	}

	return deviceName, err
}
