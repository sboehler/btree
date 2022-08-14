package btree

import (
	"math/rand"
	"testing"
)

func less(i, j int) bool {
	return i < j
}

func TestInsertAndSearch(t *testing.T) {
	n := 4000
	tree := CreateBTree(4, less)
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

func TestMin(t *testing.T) {
	n := 4000
	tree := CreateBTree(40, less)
	k, ok := tree.min()
	if ok {
		t.Fatalf("tree.min() = %v, %t, want nil, false", k, ok)
	}
	var min = n
	for _, i := range rand.Perm(n) {
		if i < min {
			min = i
		}
		tree = tree.insertTree(i)
		k, ok := tree.min()
		if !ok || k == nil {
			t.Fatalf("tree.min() = %v, %t, want <value>, true", k, ok)
		}
		if !ok || *k != min {
			t.Fatalf("tree.min() = %v, %t, want %d, true", *k, ok, min)
		}
	}
}

func TestMax(t *testing.T) {
	n := 4000
	tree := CreateBTree(40, less)
	k, ok := tree.max()
	if ok {
		t.Fatalf("tree.max() = %v, %t, want nil, false", k, ok)
	}
	var max = 0
	for _, i := range rand.Perm(n) {
		if i > max {
			max = i
		}
		tree = tree.insertTree(i)
		k, ok := tree.max()
		if !ok || k == nil {
			t.Fatalf("tree.max() = %v, %t, want <value>, true", k, ok)
		}
		if *k != max {
			t.Fatalf("tree.max() = %v, %t, want %d, true", *k, ok, max)
		}
	}
}
