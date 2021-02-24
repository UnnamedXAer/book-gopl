package main

import "fmt"

type tree struct {
	value       int
	left, right *tree
}

func (t tree) String() string {

	// ss := fmt.Sprintf("  %sL%dR%s  \n", t.left, t.value, t.right)
	ss := fmt.Sprintln(t.left, t.value, t.right)

	// s += t.visit(0, "")

	return ss
}

// func (t *tree) visit(depth int, side string) string {

// 	if t == nil {
// 		return ""
// 	}

// 	var s string = fmt.Sprintf("% 3d%s% *d\n", depth, side, depth, t.value)
// 	depth++
// 	s += t.left.visit(depth, "L")
// 	s += t.right.visit(depth, "R")
// 	return s
// }

func main() {
	treeManipulation()
}

func treeManipulation() {
	nums := []int{10, 100, 50, 500}
	// fmt.Println(nums)
	// Sort(nums)
	// fmt.Println(nums)

	t := add(nil, 15)
	t = add(t, 2)
	t = add(t, 5)
	t = add(t, 308)
	t = add(t, 200)
	t = add(t, 101)
	t = add(t, 11)
	t = add(t, -1)
	t = add(t, -1011)
	x := appendValues(nums, t)
	fmt.Println(t, x)
}

// func Sort(values []int) {
// 	var root *Tree
// 	for _, v := range values {
// 		root = add(root, v)
// 	}

// 	appendValues(values[:0], root)
// }

// // appendValues appends the elements of t to values in order
// // and return the resulting slice
// func appendValues(values []int, t *Tree) []int {
// 	if t != nil {
// 		values = appendValues(values, t.left)
// 		values = append(values, t.value)
// 		values = appendValues(values, t.right)
// 	}
// 	return values
// }

// // add adds a value to the tree and returns updated tree
// func add(t *Tree, value int) *Tree {
// 	if t == nil {
// 		// equivalent to return &tree{value: value}.
// 		t = new(Tree)
// 		t.value = value
// 		return t
// 	}

// 	if value < t.value {
// 		t.left = add(t.left, value)
// 	} else {
// 		t.right = add(t.right, value)
// 	}
// 	return t
// }

// Sort sorts values in place.
func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}
	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order
// and returns the resulting slice.
func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		// Equivalent to return &tree{value: value}.
		t = new(tree)
		t.value = value
		return t
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}
