package btree

import (
	"sort"
)

// node represents a node in a btree.
type node[K any] struct {
	// leaf indicates whether this node is a leaf node. It implies len(children) == 0.
	leaf bool

	// keys is an ordered list of K stored in this node. keys is always allocated such
	// that cap(keys) = 2*t - 1.
	keys []K

	// children is a list of children. For any i and j, children[i].keys[j] <= keys[i].
	// children is always allocated such that cap(children) == 2*t.
	children []node[K]

	// lessFn defines an ordering over K.
	lessFn func(K, K) bool
}

// Get returns a reference to the key which equals the
// given key.
func (n *node[K]) Get(key K) (*K, bool) {
	node, i := n.search(key)
	if node == nil {
		return nil, false
	}
	return &node.keys[i], true
}

// search returns
func (n *node[K]) search(key K) (*node[K], int) {
	// smallest index for which keys[index] >= key
	index := sort.Search(len(n.keys), func(i int) bool {
		// n.keys[i] >= key
		return !n.lessFn(n.keys[i], key)
	})
	if index < len(n.keys) && !n.lessFn(key, n.keys[index]) {
		// return if the keys are equal
		return n, index
	}
	if n.leaf {
		return nil, 0
	}
	return n.children[index].search(key)
}

func (n *node[K]) min() (*K, bool) {
	if n.leaf {
		if len(n.keys) > 0 {
			return &n.keys[0], true
		}
		return nil, false
	}
	return n.children[0].min()
}

func (n *node[K]) max() (*K, bool) {
	if n.leaf {
		if l := len(n.keys); l > 0 {
			return &n.keys[l-1], true
		}
		return nil, false
	}
	return n.children[len(n.children)-1].max()
}

// n is non-full, i.e. len(n.children) < 2*t -1. n.children[i] is full, i.e.
// len(n.children[i]) == 2*t - 1.
// Example for t = 2:

// Before:
//
//	c1.children = [0, 1, 2, 3]
//	c1.keys = [0, 1, 2]
//	m = 1
//
// After:
//
//	c1.children = [0, 1], c1.children = [2, 3]
//	c1.keys = [0], c2.keys = [2]
//	moved up keys: 1
func (n *node[K]) splitChild(i int) {
	t := cap(n.children) / 2
	c1 := &n.children[i]
	c2 := node[K]{
		leaf:     c1.leaf,
		lessFn:   c1.lessFn,
		keys:     make([]K, 0, 2*t-1),
		children: make([]node[K], 0, 2*t),
	}
	if !c2.leaf {
		c2.children, c1.children = move(c2.children, c1.children, t)
	}
	c2.keys, c1.keys = move(c2.keys, c1.keys, t)
	n.keys = insert(n.keys, c1.keys[t-1], i)
	c1.keys = c1.keys[:t-1]
	n.children = insert(n.children, c2, i+1)
}

// move moves the elements from src[i:] to dst, which is expected
// to be empty and have sufficient capacity.
func move[T any](dst []T, src []T, i int) ([]T, []T) {
	dst = dst[:(len(src) - i)]
	copy(dst, src[i:])
	src = src[:i]
	return dst, src
}

// insert inserts a new element e at position i, moving the elements
// after i accordingly.
func insert[T any](s []T, e T, i int) []T {
	s = append(s, e)
	if i < len(s)-1 {
		copy(s[i+1:], s[i:])
		s[i] = e
	}
	return s
}

func (n *node[K]) insertTree(key K) *node[K] {
	if len(n.keys) == cap(n.keys) {
		// optimization to avoid splitting the root
		// node if the key already exists
		if nd, i := n.search(key); nd != nil {
			nd.keys[i] = key
			return n
		}
		// split root
		r := &node[K]{
			leaf:     false,
			lessFn:   n.lessFn,
			keys:     make([]K, 0, cap(n.keys)),
			children: make([]node[K], 0, cap(n.children)),
		}
		r.children = append(r.children, *n)
		r.splitChild(0)
		n = r
	}
	n.insertTreeNotFull(key)
	return n
}

func (n *node[K]) insertTreeNotFull(key K) {
	// smallest i for which keys[i] >= key
	index := sort.Search(len(n.keys), func(i int) bool {
		return !n.lessFn(n.keys[i], key)
	})
	if index < len(n.keys) && !n.lessFn(key, n.keys[index]) {
		// replace equal key
		n.keys[index] = key
		return
	}
	if n.leaf {
		if index == len(n.keys) {
			n.keys = append(n.keys, key)
			return
		}
		n.keys = append(n.keys, key)
		if index < len(n.keys)-1 {
			copy(n.keys[index+1:], n.keys[index:])
			n.keys[index] = key
		}
	} else {
		// keys[index] > keys
		ch := n.children[index]
		if len(ch.keys) == cap(ch.keys) {
			// optimization to avoid splitting the
			// node if the key already exists
			if nd, i := ch.search(key); nd != nil {
				nd.keys[i] = key
				return
			}
			n.splitChild(index)
			if n.lessFn(n.keys[index], key) {
				index++
			} else if !n.lessFn(key, n.keys[index]) {
				// replace equal key
				n.keys[index] = key
				return
			}
		}
		n.children[index].insertTreeNotFull(key)
	}
}

func CreateBTree[K any](t int, lessFn func(K, K) bool) *node[K] {
	return &node[K]{
		leaf:     true,
		lessFn:   lessFn,
		keys:     make([]K, 0, 2*t-1),
		children: make([]node[K], 0, 2*t),
	}
}
