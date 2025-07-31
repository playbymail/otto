// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package vm

import (
	"fmt"
	"github.com/playbymail/otto/wjs/domain"
	"strings"
)

// builtinFunc is the concrete type for built-in functions.
type builtinFunc struct {
	name  string
	arity int // use -1 for variadic
	fn    func(pos domain.Pos, args []Value) (Value, *RuntimeError)
}

func (b *builtinFunc) Call(pos domain.Pos, args []Value) (Value, *RuntimeError) {
	if b.arity >= 0 && len(args) != b.arity {
		return nil, NewRuntimeError(pos, "%s expects %d arguments, got %d", b.name, b.arity, len(args))
	}
	return b.fn(pos, args)
}

func (b *builtinFunc) Name() string {
	return b.name
}

func (b *builtinFunc) Arity() int {
	return b.arity
}

// RegisterBuiltins returns a map of standard built-in functions.
func RegisterBuiltins(loadFn func(path string) (*Map, error), saveFn func(*Map, string) error) map[string]Value {
	return map[string]Value{
		"print": &builtinFunc{
			name:  "print",
			arity: -1,
			fn: func(pos domain.Pos, args []Value) (Value, *RuntimeError) {
				out := make([]string, len(args))
				for i, arg := range args {
					out[i] = Stringify(arg)
				}
				fmt.Println(strings.Join(out, " "))
				return nil, nil
			},
		},

		"load": &builtinFunc{
			name:  "load",
			arity: 1,
			fn: func(pos domain.Pos, args []Value) (Value, *RuntimeError) {
				path, ok := args[0].(string)
				if !ok {
					return nil, NewRuntimeError(pos, "load expects a string path")
				}
				m, err := loadFn(path)
				if err != nil {
					return nil, NewRuntimeError(pos, "load error: %v", err)
				}
				return m, nil
			},
		},

		"save": &builtinFunc{
			name:  "save",
			arity: 2,
			fn: func(pos domain.Pos, args []Value) (Value, *RuntimeError) {
				mapPtr := args[0]
				path, ok := args[1].(string)
				if !ok {
					return nil, NewRuntimeError(pos, "save expects a string as the second argument")
				}
				m, ok := mapPtr.(*Map) // type check mapPtr is *Map
				if !ok {
					return nil, NewRuntimeError(pos, "save expects a Map as the first argument")
				}
				if err := saveFn(m, path); err != nil {
					return nil, NewRuntimeError(pos, "save error: %v", err)
				}
				return nil, nil
			},
		},
	}
}
