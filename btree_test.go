package btree

import (
	"math/rand"
	"testing"
)

func less(i, j int) bool {
	return i < j
}

func TestSmoke(t *testing.T) {
	n := 4000
	tree := createBTree(4, less)
	for _, i := range rand.Perm(n) {
		tree = tree.insertTree(i)
	}
	for _, i := range rand.Perm(n) {
		tree = tree.insertTree(i)
	}
	for i := 0; i < n; i++ {
		n, index := tree.search(i)
		if n == nil {
			t.Fatalf("tree.search(%d) = nil, %d, want not nil", i, index)
		}
		if n.keys[index] != i {
			t.Fatalf("%v.keys[%d] = %d, want %d", n, index, n.keys[index], i)
		}
	}
}
