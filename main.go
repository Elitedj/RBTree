package main

func main() {
	rbt := NewRBTree[int]()
	rbt.Insert(1)
	rbt.Has(1)
	rbt.Delete(1)
	rbt.Clear()
}
