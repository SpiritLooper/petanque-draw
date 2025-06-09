package intset

import (
	"errors"
	"math/rand"
)

type IntSet map[int]struct{}

func CreateIntSet(n int) IntSet {
	set := make(IntSet)
	for i := range n {
		set[i] = struct{}{}
	}
	return set
}

func (set IntSet) IsEmpty() bool {
	return len(set) == 0
}

func (set IntSet) Count() int {
	return len(set)
}

func (set *IntSet) RandomPop() (int, error) {
	if set.IsEmpty() {
		return 0, errors.New("Calling RandomPop on empty Set")
	}

	keys := make([]int, 0, len(*set))
	for k, _ := range *set {
		keys = append(keys, k)
	}

	randIndex := rand.Intn(len(keys))
	val := keys[randIndex]
	delete(*set, val)

	return val, nil
}
