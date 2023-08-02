package model

import (
	"fmt"
	"testing"
)

func TestRemoveSlice(t *testing.T) {

	s1 := []int{1, 2, 3, 4, 5, 6, 7}
	s2 := []int{2, 4, 6, 8}

	for _, card := range s2 {
		for j, old := range s1 {
			if old == card {
				s1 = append(s1[:j], s1[j+1:]...)
				break
			}
		}
	}

	fmt.Printf("%#v", s1)
}
