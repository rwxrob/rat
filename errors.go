package rat

import "fmt"

// ---------------------------- ErrExpected ---------------------------

type ErrExpected struct{ any }

func (e ErrExpected) Error() string { return fmt.Sprintf(_ErrExpected, e.any) }

// ---------------------------- ErrNotExist ---------------------------

type ErrNotExist struct{ any }

func (e ErrNotExist) Error() string {
	return fmt.Sprintf(_ErrNotExist, e.any)
}

// ---------------------------- ErrBadType ----------------------------

type ErrBadType struct{ any }

func (e ErrBadType) Error() string {
	return fmt.Sprintf(_ErrBadType, e.any)
}

// ------------------------------ ErrArgs -----------------------------

type ErrArgs struct{ any }

func (e ErrArgs) Error() string {
	return fmt.Sprintf(_ErrArgs, e.any)
}
