package rat

import "encoding/json"

// Match is the fundamental result of a parse function and contains only
// reference information to the underlying, in-memory buffer. Buf is but
// a slice abstraction pointing to the same underlying array the same as
// all other matches.
type Match struct {
	Buf []rune  `json:"-"`           // back reference to data (copy of slice only)
	Beg int     `json:"B"`           // inclusive
	End int     `json:"E"`           // non-inclusive
	Sub []Match `json:"S,omitempty"` // equivalent to parens in regexp
}

// String fulfills the fmt.Stringer interface as JSON with "null" for
// any errors.
func (m Match) String() string {
	buf, err := json.Marshal(m)
	if err != nil {
		return `null`
	}
	return string(buf)
}
