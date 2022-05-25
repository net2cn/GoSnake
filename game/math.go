package game

import (
	"golang.org/x/exp/constraints"
)

type number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Abs[T number](x T) T {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
