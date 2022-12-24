package rat

import (
	"fmt"
)

type ErrLit struct {
	Lit string
}

func (e ErrLit) Error() string {
	return fmt.Sprintf("expected literal %q", e.Lit)
}
