package rat

import (
	"fmt"

	"github.com/rwxrob/rat/pegn"
)

type ErrExpected struct {
	This any
}

func (e ErrExpected) Error() string {
	switch v := e.This.(type) {
	case rune:
		e.This = pegn.FromString(string(v))
	}
	return fmt.Sprintf(_ErrExpected, e.This)
}

type ErrNotExist struct{ This any }

func (e ErrNotExist) Error() string {
	return fmt.Sprintf(_ErrNotExist, e.This)
}
