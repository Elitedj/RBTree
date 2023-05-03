package main

import "golang.org/x/exp/constraints"

type Color byte

const (
	Red Color = iota
	Black
)

type Direction byte

const (
	DirRoot Direction = iota
	DirLeft
	DirRight
)

type RBNode[T constraints.Ordered] struct {
	Val    T
	parent *RBNode[T]
	left   *RBNode[T]
	right  *RBNode[T]
	color  Color
}

func newRBNode[T constraints.Ordered](v T) *RBNode[T] {
	return &RBNode[T]{
		Val:    v,
		parent: nil,
		left:   nil,
		right:  nil,
		color:  Red,
	}
}

func (n *RBNode[T]) isLeaf() bool {
	return (n.left == nil && n.right == nil)
}

func (n *RBNode[T]) isRoot() bool {
	return n.parent == nil
}

func (n *RBNode[T]) isRed() bool {
	return n.color == Red
}

func (n *RBNode[T]) isBlack() bool {
	return n.color == Black
}

func (n *RBNode[T]) direction() Direction {
	if n.parent != nil {
		if n == n.parent.left {
			return DirLeft
		} else {
			return DirRight
		}
	} else {
		return DirRoot
	}
}

func (n *RBNode[T]) sibling() *RBNode[T] {
	if n.direction() == DirLeft {
		return n.parent.right
	} else {
		return n.parent.left
	}
}

func (n *RBNode[T]) hasSibling() bool {
	return (!n.isRoot() && n.sibling() != nil)
}

func (n *RBNode[T]) uncle() *RBNode[T] {
	return n.parent.sibling()
}

func (n *RBNode[T]) hasUncle() bool {
	return (!n.isRoot() && n.parent.hasSibling())
}

func (n *RBNode[T]) grandParent() *RBNode[T] {
	return n.parent.parent
}

func (n *RBNode[T]) hasGrandParent() bool {
	return (!n.isRoot() && !n.parent.isRoot())
}

func (n *RBNode[T]) release() {
	n.parent = nil
	if n.left != nil {
		n.left.release()
	}
	if n.right != nil {
		n.right.release()
	}
}

type RBTree[T constraints.Ordered] struct {
	root *RBNode[T]
	cnt  uint
}

func NewRBTree[T constraints.Ordered]() *RBTree[T] {
	return &RBTree[T]{
		root: nil,
		cnt:  0,
	}
}

func (t *RBTree[T]) Size() uint {
	return t.cnt
}

func (t *RBTree[T]) Empty() bool {
	return t.cnt == 0
}

func (t *RBTree[T]) Clear() {
	if t.root != nil {
		t.root.release()
		t.root = nil
	}
	t.cnt = 0
}

func search[T constraints.Ordered](n *RBNode[T], v T) (*RBNode[T], bool) {
	if n == nil {
		return nil, false
	}

	if n.Val == v {
		return n, true
	}
	if n.Val < v {
		return search(n.right, v)
	}
	return search(n.left, v)
}

func (t *RBTree[T]) Get(v T) (*RBNode[T], bool) {
	return search(t.root, v)
}

func (t *RBTree[T]) Has(v T) bool {
	_, ok := search(t.root, v)
	return ok
}

func (t *RBTree[T]) Insert(v T) {
	if t.root == nil {
		t.root = newRBNode(v)
		t.root.color = Black
		t.cnt++
	} else {
		t.insert(t.root, v)
	}
}

func (n *RBNode[T]) maintainRelationship() {
	if n.left != nil {
		n.left.parent = n
	}
	if n.right != nil {
		n.right.parent = n
	}
}

func (t *RBTree[T]) leftRotate(n *RBNode[T]) {
	if n == nil || n.right == nil {
		return
	}

	parent := n.parent
	dir := n.direction()

	successor := n.right
	n.right = successor.left
	successor.left = n

	n.maintainRelationship()
	successor.maintainRelationship()

	switch dir {
	case DirRoot:
		t.root = successor
		break
	case DirLeft:
		parent.left = successor
		break
	case DirRight:
		parent.right = successor
		break
	}

	successor.parent = parent
}

func (t *RBTree[T]) rightRotate(n *RBNode[T]) {
	if n == nil || n.left == nil {
		return
	}

	parent := n.parent
	dir := n.direction()

	successor := n.left
	n.left = successor.right
	successor.right = n

	n.maintainRelationship()
	successor.maintainRelationship()

	switch dir {
	case DirRoot:
		t.root = successor
		break
	case DirLeft:
		parent.left = successor
		break
	case DirRight:
		parent.right = successor
		break
	}

	successor.parent = parent
}

func (t *RBTree[T]) insert(n *RBNode[T], v T) {
	if n.Val == v {
		return
	}

	if v < n.Val {
		if n.left == nil {
			n.left = newRBNode(v)
			n.left.parent = n
			t.insertFixUp(n.left)
			t.cnt++
		} else {
			t.insert(n.left, v)
		}
	} else {
		if n.right == nil {
			n.right = newRBNode(v)
			n.right.parent = n
			t.insertFixUp(n.right)
			t.cnt++
		} else {
			t.insert(n.right, v)
		}
	}
}

func (t *RBTree[T]) insertFixUp(n *RBNode[T]) {
	if n.isRoot() {
		n.color = Black
		return
	}

	if n.parent.isBlack() {
		// case 1: parent is black
		return
	}

	if n.parent.isRoot() {
		n.parent.color = Black
		return
	}

	if !n.hasUncle() || n.uncle().isBlack() {
		// case 2a: uncle is nil or black
		if n.direction() != n.parent.direction() {
			parent := n.parent
			if n.direction() == DirLeft {
				t.rightRotate(n.parent)
			} else {
				t.leftRotate(n.parent)
			}
			n = parent
		}

		if n.parent.direction() == DirLeft {
			t.rightRotate(n.grandParent())
		} else {
			t.leftRotate(n.grandParent())
		}

		n.parent.color = Black
		n.sibling().color = Red
		return
	}

	if n.hasUncle() && n.uncle().isRed() {
		// case 2b: uncle is red
		n.parent.color = Black
		n.uncle().color = Black
		n.grandParent().color = Red
		t.insertFixUp(n.grandParent())
		return
	}
}

func (t *RBTree[T]) Delete(v T) {
	if t.root == nil {
		return
	} else {
		t.delete(t.root, v)
	}
}

func (t *RBTree[T]) delete(n *RBNode[T], v T) {
	if n == nil {
		return
	}

	if n.Val != v {
		if v < n.Val {
			left := n.left
			t.delete(left, v)
		} else {
			right := n.right
			t.delete(right, v)
		}
		n.maintainRelationship()
		return
	}

	// only root
	if t.Size() == 1 {
		t.Clear()
		return
	}

	if n.left != nil && n.right != nil {
		successor := n.right
		parent := n
		// find successor of current
		for successor.left != nil {
			parent = successor
			successor = parent.left
		}
		// swap current and successor
		n.Val, successor.Val = successor.Val, n.Val
		n, successor = successor, n
		parent.maintainRelationship()
	}

	if n.isLeaf() {
		// black leaf node
		if n.isBlack() {
			t.deleteFixUp(n)
		}
		// delete
		if n.direction() == DirLeft {
			n.parent.left = nil
		} else {
			n.parent.right = nil
		}
	} else {
		// current node has a signle child
		// replace current with its child
		parent := n.parent
		var replacement *RBNode[T]
		if n.left != nil {
			replacement = n.left
		} else {
			replacement = n.right
		}
		switch n.direction() {
		case DirRoot:
			t.root = replacement
			break
		case DirLeft:
			parent.left = replacement
			break
		case DirRight:
			parent.right = replacement
			break
		}

		if !n.isRoot() {
			replacement.parent = parent
		}

		if n.isBlack() {
			if replacement.isRed() {
				replacement.color = Black
			} else {
				t.deleteFixUp(replacement)
			}
		}

		n = nil
	}

	t.cnt--
}

func (t *RBTree[T]) deleteFixUp(n *RBNode[T]) {
	if n.isRoot() {
		n.color = Black
		return
	}

	dir := n.direction()
	sibling := n.sibling()
	if sibling.isRed() {
		// case 1: current is black and sibling is red
		// step 1: recolor parent and sibling
		// step 2: rotate parent
		parent := n.parent
		parent.color = Red
		sibling.color = Black
		if dir == DirLeft {
			t.leftRotate(parent)
		} else {
			t.rightRotate(parent)
		}
		sibling = n.sibling()
	}

	var closeNephew, distantNephew *RBNode[T]
	var isCloseNephewBlack, isDistantNephewBlack bool
	if dir == DirLeft {
		closeNephew = sibling.left
		distantNephew = sibling.right
	} else {
		closeNephew = sibling.right
		distantNephew = sibling.left
	}
	if closeNephew == nil || closeNephew.isBlack() {
		isCloseNephewBlack = true
	}
	if distantNephew == nil || distantNephew.isBlack() {
		isDistantNephewBlack = true
	}

	if isCloseNephewBlack && isDistantNephewBlack {
		// case 2: sibling is black and has two black childs
		if n.parent.isRed() {
			n.parent.color = Black
			sibling.color = Red
		} else {
			sibling.color = Red
			t.deleteFixUp(n.parent)
		}
		return
	} else {
		if closeNephew != nil && closeNephew.isRed() {
			// case 3: sibling is black and close nephew is red
			closeNephew.color = Black
			sibling.color = Red
			if dir == DirLeft {
				t.rightRotate(sibling)
			} else {
				t.leftRotate(sibling)
			}
			sibling = n.sibling()
			if dir == DirLeft {
				closeNephew = sibling.left
				distantNephew = sibling.right
			} else {
				closeNephew = sibling.right
				distantNephew = sibling.left
			}
		}

		// case 4: sibling is black and distant nephew is red
		sibling.color = n.parent.color
		n.parent.color = Black
		if distantNephew != nil {
			distantNephew.color = Black
		}
		if dir == DirLeft {
			t.leftRotate(n.parent)
		} else {
			t.rightRotate(n.parent)
		}
		t.root.color = Black
		return
	}
}

func (t *RBTree[T]) InOrder() []T {
	res := make([]T, 0)
	t.inOrder(t.root, &res)
	return res
}

func (t *RBTree[T]) inOrder(n *RBNode[T], res *[]T) {
	if n == nil {
		return
	}

	if n.left != nil {
		t.inOrder(n.left, res)
	}
	*res = append(*res, n.Val)
	if n.right != nil {
		t.inOrder(n.right, res)
	}
}
