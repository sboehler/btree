package btree

import (
	"sort"
)

type node[K any] struct {
	leaf     bool
	keys     []K
	children []node[K]
	lessFn   func(K, K) bool
}

func (n *node[K]) search(key K) (*node[K], int) {
	// smallest index for which keys[index] > key
	index := sort.Search(len(n.keys), func(i int) bool {
		return !n.lessFn(n.keys[i], key)
	})
	if index < len(n.keys) && !n.lessFn(key, n.keys[index]) {
		return n, index
	}
	if n.leaf {
		return nil, 0
	}
	return n.children[index].search(key)
}

func (n *node[K]) splitChild(i int) {
	t := cap(n.children) / 2
	c1 := &n.children[i]
	c2 := &node[K]{
		leaf:   c1.leaf,
		lessFn: c1.lessFn,
		keys:   make([]K, t-1, 2*t-1),
	}
	if c2.leaf {
		c2.children = make([]node[K], 0, 2*t)
	} else {
		c2.children = make([]node[K], t, 2*t)
		copy(c2.children, c1.children[t:])
		c1.children = c1.children[:t]
	}
	m := c1.keys[t-1]
	copy(c2.keys, c1.keys[t:])
	c1.keys = c1.keys[:t-1]
	// lift the key up
	n.keys = append(n.keys, m)
	if i < len(n.keys)-1 {
		copy(n.keys[i+1:], n.keys[i:])
	}
	n.keys[i] = m
	// insert new child
	n.children = append(n.children, *c2)
	if i < len(n.children)-2 {
		copy(n.children[i+2:], n.children[i+1:])
		n.children[i+1] = *c2
	}
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

func createBTree[K any](t int, lessFn func(K, K) bool) *node[K] {
	return &node[K]{
		leaf:     true,
		lessFn:   lessFn,
		keys:     make([]K, 0, 2*t-1),
		children: make([]node[K], 0, 2*t),
	}
}
