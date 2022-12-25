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

type ErrOneOf struct {
	Rules []Rule
}

func (e ErrOneOf) Error() string {
	names := make([]string, len(e.Rules))
	for i, rule := range e.Rules {
		names[i] = FuncName(rule)
	}
	return fmt.Sprintf("expected one of %v", names)
}
