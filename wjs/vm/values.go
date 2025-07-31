// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package vm

import (
	"fmt"
	"github.com/playbymail/otto/wjs/domain"
	"reflect"
	"strconv"
	"strings"
)

// Value is the unified interface for all runtime values.
type Value interface{}

// Object represents a dynamic map of string keys to values (like a JS object).
type Object map[string]Value

// Array represents a list of values.
type Array []Value

// Map is the user-defined Worldographer map type returned by `load`.
type Map = any // replace with your actual *Map type

// RuntimeError is returned on any execution failure.
type RuntimeError struct {
	Pos     domain.Pos
	Message string
}

func (e *RuntimeError) Error() string {
	if e.Pos.Script == "" {
		return fmt.Sprintf("Runtime error at %d:%d: %s", e.Pos.Line, e.Pos.Column, e.Message)
	}
	return fmt.Sprintf("Runtime error at %s:%d:%d: %s", e.Pos.Script, e.Pos.Line, e.Pos.Column, e.Message)
}

func NewRuntimeError(pos domain.Pos, format string, args ...any) *RuntimeError {
	return &RuntimeError{
		Pos:     pos,
		Message: fmt.Sprintf(format, args...),
	}
}

// Type checking helpers

func IsNumber(v Value) bool {
	switch v.(type) {
	case int64, float64:
		return true
	default:
		return false
	}
}

func IsString(v Value) bool {
	_, ok := v.(string)
	return ok
}

func IsBool(v Value) bool {
	_, ok := v.(bool)
	return ok
}

func IsNull(v Value) bool {
	return v == nil
}

func IsArray(v Value) bool {
	_, ok := v.([]Value)
	return ok
}

func IsObject(v Value) bool {
	_, ok := v.(Object)
	return ok
}

// Stringify converts any value into its string form (used in template strings).
func Stringify(v Value) string {
	switch val := v.(type) {
	case nil:
		return "null"
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case string:
		return val
	case []Value:
		var sb strings.Builder
		sb.WriteString("[")
		for i, el := range val {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(Stringify(el))
		}
		sb.WriteString("]")
		return sb.String()
	case Object:
		var sb strings.Builder
		sb.WriteString("{")
		i := 0
		for k, v := range val {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(k)
			sb.WriteString(": ")
			sb.WriteString(Stringify(v))
			i++
		}
		sb.WriteString("}")
		return sb.String()
	default:
		// use reflection for struct types (e.g., Map)
		return fmt.Sprintf("%v", val)
	}
}

// Equal returns true if a and b are deeply equal.
func Equal(a, b Value) bool {
	return reflect.DeepEqual(a, b)
}

// ToFloat64 converts any number (int64 or float64) to float64
func ToFloat64(v Value) (float64, bool) {
	switch val := v.(type) {
	case int64:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

// PromoteNumbers converts two numeric values and returns them as the appropriate type.
// If both are int64, returns int64. If either is float64, both are promoted to float64.
func PromoteNumbers(a, b Value) (Value, Value, bool) {
	aInt, aIsInt := a.(int64)
	bInt, bIsInt := b.(int64)
	aFloat, aIsFloat := a.(float64)
	bFloat, bIsFloat := b.(float64)
	
	// Both must be numbers
	if !(aIsInt || aIsFloat) || !(bIsInt || bIsFloat) {
		return nil, nil, false
	}
	
	// If both are integers, keep them as integers
	if aIsInt && bIsInt {
		return aInt, bInt, true
	}
	
	// Otherwise, promote both to float64
	if aIsInt {
		aFloat = float64(aInt)
	}
	if bIsInt {
		bFloat = float64(bInt)
	}
	
	return aFloat, bFloat, true
}
