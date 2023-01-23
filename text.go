package rat

// KEEP APP TEXT HERE
// (This should be the only file to need translation, if needed.)

const (
	ErrIsZeroT      = `zero value: %T`
	ErrNotExistT    = `does not exist: %v`
	ErrExpectedT    = `expected: %v`
	ErrBadTypeT     = `unknown type: %v (%[1]T)`
	ErrArgsT        = `missing or incorrect arguments: %v (%[1]T)`
	ErrPackTypeT    = `invalid type`
	ErrNoCheckFuncT = `no check function assigned: %v`
)
