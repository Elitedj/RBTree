package main_test

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	rbt "github.com/elitedj/rbtree"
)

func TestRBTree(t *testing.T) {
	tree := rbt.NewRBTree[int]()

	n := 100000
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	nums := rnd.Perm(n)
	for _, v := range nums {
		tree.Insert(v)
	}

	res := tree.InOrder()
	if !sort.IntsAreSorted(res) || len(res) != n {
		t.Error("error with inorder")
	}

	sort.Ints(nums)

	for i := 0; i < n; i++ {
		if nums[i] != res[i] {
			t.Error("error with inorder")
		}
	}

	ord := rnd.Perm(n)
	for i := 0; i < n; i++ {
		node, ok := tree.Get(ord[i])
		if !ok {
			t.Errorf("error with Get: %d not found", ord[i])
		}
		if node.Val != ord[i] {
			t.Errorf("error with Get: node.Val is %d, but %d expected", node.Val, ord[i])
		}

		tree.Delete(ord[i])
		ok = tree.Has(ord[i])
		if ok {
			t.Errorf("error with Delete: %d was deleted but tree had %d", ord[i], ord[i])
		}
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := rbt.NewRBTree[int]()
	for i := 0; i < b.N; i++ {
		tree.Insert(i)
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()

	tree := rbt.NewRBTree[int]()
	for i := 0; i < b.N; i++ {
		tree.Insert(i)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(i)
	}
}
