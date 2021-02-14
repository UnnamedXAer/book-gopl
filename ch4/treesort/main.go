package main

import "fmt"

type tree struct {
	value       int
	left, right *tree
}

func main() {
	treeManipulation()
}

func treeManipulation() {
	nums := []int{17, 4, 32, 2, 56, 67, 7, 654, 4, 56, 2, 3, 4, 43}
	fmt.Println(nums)
	Sort(nums)
	fmt.Println(nums)
}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v)
	}

	appendValues(values[:0], root)
}

// appendValues appends the elements of t to values in order
// and return the resulting slice
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
		// equivalent to return &tree{value: value}.
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
