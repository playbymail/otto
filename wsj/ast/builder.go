// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package ast

func ApplySuffixes(base Expr, suffixes []Suffix) Expr {
	for _, s := range suffixes {
		base = s.Apply(base)
	}
	return base
}
