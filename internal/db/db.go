package db

import (
	"errors"
)

var (
	ErrAlreadyExists = errors.New("object with the same ID already exists")
	ErrNotFound      = errors.New("object not found")
)

type TxMode int

const (
	TxModeReadOnly TxMode = iota
	TxModeReadWrite
)

type DB interface {
	Init() error
	Close() error
	// Begin(mode TxMode) (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback()
	DumpData() ([]byte, error)
	LoadData(data []byte) error
}
