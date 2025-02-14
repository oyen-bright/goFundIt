package database

import (
	"errors"

	"gorm.io/gorm"
)

type DBError interface {
	error
	Is(err error) bool
	IsNotfound() bool
	IsForeignKeyViolated() bool
	IsPrimaryKeyViolated() bool
}

type WrapError struct {
	err error
}

func (e *WrapError) Error() string {
	return e.err.Error()
}

func (e *WrapError) Is(err error) bool {
	return errors.Is(e.err, err)
}

func (e *WrapError) IsNotfound() bool {
	return errors.Is(e.err, gorm.ErrRecordNotFound)
}

func (e *WrapError) IsForeignKeyViolated() bool {
	return errors.Is(e.err, gorm.ErrForeignKeyViolated)
}

func (e *WrapError) IsPrimaryKeyViolated() bool {
	return errors.Is(e.err, gorm.ErrDuplicatedKey)
}

func Error(err error) DBError {
	return &WrapError{err: err}
}
