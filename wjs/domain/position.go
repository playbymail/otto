// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package domain

// Pos represents a position in the source file.
type Pos struct {
	Line   int    // 1-based
	Column int    // 1-based
	Offset int    // byte offset from start of file
	Script string // set only when processing a script file
}
