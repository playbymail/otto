// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package vm

import (
	"github.com/playbymail/otto/wjs/domain"
)

// Callable is the interface for all VM functions (built-in or user-defined).
type Callable interface {
	Call(pos domain.Pos, args []Value) (Value, *RuntimeError)
	Name() string
	Arity() int // -1 for variadic
}
