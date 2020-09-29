package badgerdb

func (t *badgerTX) GetOrCreateDeviceName(defaultDeviceName string) (string, error) {
	value, err := getOrCreate(t.txn, "deviceName", []byte(defaultDeviceName))
	if err != nil {
		return "", err
	}

	return string(value), err
}
