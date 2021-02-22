package main

import (
	"fmt"
	"sort"
)

// preReqs maps computer science courses to their prerequisites.
var preReqs = map[string][]string{
	"algorithms": {"data structures"},
	"calculus":   {"linear algebra"},

	//added
	"linear algebra": {"calculus"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to promramming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	// for i, course := range topoSort(preReqs) {
	// 	fmt.Printf("%d:\t%s\n", i, course)
	// }
	// fmt.Println()
	fmt.Println()
	for name, pos := range topoSortM(preReqs) {
		fmt.Printf("%d:\t%s\n", pos, name)
	}
}

func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)

	// visitAll declated in two steps to be able to call it recursively
	// if combined into one step it wouldn't see itself
	var visitAll func(items []string)
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				order = append(order, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return order
}

func topoSortM(m map[string][]string) map[string]int {
	// var order []string
	order := make(map[string]int, len(m))
	seen := make(map[string]bool)

	// visitAll declated in two steps to be able to call it recursively
	// if combined into one step visitAll wouldn't see itself
	var visitAll func(items []string)
	var pos int
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				for _, v := range m[item] {

					for _, mv := range m[v] {
						if mv == item {
							fmt.Printf("warning, cycle detected, item: %q req %q\n", item, v)
						}
					}

					if v == item {
						fmt.Printf("warning, item: %q requires itself\n", item)
					}
				}

				seen[item] = true
				visitAll(m[item])
				pos++
				order[item] = pos
				// order = append(order, item)
			}
		}
	}

	// fmt.Println(len(m))
	// x := 0
	// // var keys []string
	// var keys []string = make([]string, len(m))
	// for key := range m {
	// 	fmt.Printf("ptr: %p len: % 2d,\t cap: % 2d\n", keys, len(keys), cap(keys))
	// 	// keys = append(keys, key)
	// 	keys[x] = key
	// 	x++
	// }
	// fmt.Printf("ptr: %p len: % 2d,\t cap: % 2d\n", keys, len(keys), cap(keys))

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return order
}
