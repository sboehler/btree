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
		tree.InsertTree(i)
	}
	for _, i := range rand.Perm(n) {
		tree.InsertTree(i)
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

func TestIterate(t *testing.T) {
	n := 4000
	tree := CreateBTree(4, less)
	for _, i := range rand.Perm(n) {
		tree.InsertTree(i)
	}
	for _, i := range rand.Perm(n) {
		tree.InsertTree(i)
	}
	var want int
	tree.Iterate(func(got *int) bool {
		if *got != want {
			t.Fatalf("tree.Iterate(...) = %d, want %d", *got, want)
		}
		want++
		return true
	})
	var counter int
	ok := tree.Iterate(func(got *int) bool {
		if *got == 500 {
			return false
		}
		counter++
		return true
	})
	if ok || counter != 500 {
		t.Fatalf("tree.Iterate(...) = %t, %d, want false, 500", ok, counter)
	}
}

func TestMin(t *testing.T) {
	n := 4000
	tree := CreateBTree(40, less)
	k, ok := tree.Min()
	if ok {
		t.Fatalf("tree.min() = %v, %t, want nil, false", k, ok)
	}
	var min = n
	for _, i := range rand.Perm(n) {
		if i < min {
			min = i
		}
		tree.InsertTree(i)
		k, ok := tree.Min()
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
	k, ok := tree.Max()
	if ok {
		t.Fatalf("tree.max() = %v, %t, want nil, false", k, ok)
	}
	var max = 0
	for _, i := range rand.Perm(n) {
		if i > max {
			max = i
		}
		tree.InsertTree(i)
		k, ok := tree.Max()
		if !ok || k == nil {
			t.Fatalf("tree.max() = %v, %t, want <value>, true", k, ok)
		}
		if *k != max {
			t.Fatalf("tree.max() = %v, %t, want %d, true", *k, ok, max)
		}
	}
}

func TestSearchGE(t *testing.T) {
	n := 4000
	tree := CreateBTree(40, less)
	seq := rand.Perm(n)
	key := seq[len(seq)-1]
	node, index := tree.SearchGE(key)
	if node != nil {
		t.Fatalf("tree.SearchGE() = %v, %d, want nil, 0", node, index)
	}
	var min = n
	tree.InsertTree(n)
	node, index = tree.SearchGE(key)
	if node == nil {
		t.Fatalf("tree.SearchGE() = %v, %d, want <node>, <index>", node, index)
	}
	if node.keys[index] != min {
		t.Fatalf("tree.SearchGE() = %d, want %d", node.keys[index], n)
	}

	for _, i := range rand.Perm(n) {
		if i < min && i >= key {
			min = i
		}
		tree.InsertTree(i)
		node, index := tree.SearchGE(key)
		if node == nil {
			t.Fatalf("tree.SearchGE() = %v, %d, want <node>, <index>", node, index)
		}
		if node.keys[index] != min {
			t.Fatalf("tree.SearchGE() = %d, want %d", node.keys[index], min)
		}
	}
}
