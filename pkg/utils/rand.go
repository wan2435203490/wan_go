package utils

import "math/rand"

func RandIntRange(from, to int) int {
	if to-from <= 0 {
		return from
	}
	return rand.Intn(to-from) + from
}
func RandFloat64Range(from, to float64) float64 {
	return rand.Float64()*(to-from) + from
}
