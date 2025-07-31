// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package vm provides the runtime for the WJS scripting language.
package vm

import (
	"github.com/playbymail/otto/wjs/ast"
	"github.com/playbymail/otto/wjs/domain"
)

func New(script string) *VM {
	return &VM{
		vars:   map[string]Value{},
		script: script,
	}
}

type VM struct {
	vars map[string]Value // for loop vars, etc.
	// funcs    map[string]CallableFunction // built-in function handlers
	script string // current script filename
}

func (vm *VM) Execute(program *ast.Program) error {
	return domain.ErrNotImplemented
}
