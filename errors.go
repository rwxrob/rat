package rat

import "fmt"

// ----------------------------- ErrIsZero ----------------------------

type ErrIsZero struct{ V any }

func (e ErrIsZero) Error() string { return fmt.Sprintf(_ErrIsZero, e.V) }

// ---------------------------- ErrExpected ---------------------------

type ErrExpected struct{ V any }

func (e ErrExpected) Error() string { return fmt.Sprintf(_ErrExpected, e.V) }

// ---------------------------- ErrNotExist ---------------------------
