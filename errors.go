package rat

import "fmt"

// ----------------------------- ErrIsZero ----------------------------

type ErrIsZero struct{ V any }

func (e ErrIsZero) Error() string { return fmt.Sprintf(ErrIsZeroT, e.V) }

// ---------------------------- ErrExpected ---------------------------

type ErrExpected struct{ V any }

func (e ErrExpected) Error() string { return fmt.Sprintf(ErrExpectedT, e.V) }

// ------------------------------ ErrArgs -----------------------------

type ErrArgs struct{ any }

func (e ErrArgs) Error() string {
	return fmt.Sprintf(ErrArgsT, e.any)
}

// -------------------------- ErrNoCheckFunc --------------------------

type ErrNoCheckFunc struct{ V any }

func (e ErrNoCheckFunc) Error() string { return fmt.Sprintf(ErrNoCheckFuncT, e.V) }

// ---------------------------- ErrNotFound ---------------------------

type ErrNotFound struct{ any }
