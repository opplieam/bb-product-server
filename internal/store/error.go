package store

import (
	"errors"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/lib/pq"
)

var ErrRecordNotFound = errors.New("record not found")
var ErrRecordAlreadyExists = errors.New("record already exists")

func DBTransformError(err error) error {
	if err == nil {
		return nil
	}

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		assertPq, ok := err.(*pq.Error)
		if ok {
			switch assertPq.Code {
			case "23505":
				return ErrRecordAlreadyExists
			}
		}

	}

	if errors.Is(err, qrm.ErrNoRows) {
		return ErrRecordNotFound
	}

	return err
}
