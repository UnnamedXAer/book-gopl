package main

import (
	"bytes"
	"fmt"
)

const uintSize int = 32 << (^uint(0) >> 63)

// IntSet is a set of small non-negative integers.
// Its zero value represents the empty state.
type IntSet struct {
	words []uint
}

// Has reports the set contains the non-negative value x
func (s *IntSet) Has(x int) bool {
	word, bit := x/uintSize, uint(x%uintSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set
func (s *IntSet) Add(x int) {
	word, bit := x/uintSize, uint(x%uintSize)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// AddAll adds all passed non-negative parameters values to the set.
func (s *IntSet) AddAll(x ...int) {
	for _, v := range x {
		s.Add(v)
	}
}

// UnionWith sets s to the union of s and t
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String returns the set as a string of the form "{1 2 3}"
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')

	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < uintSize; j++ {
			if word&(1<<(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", uintSize*i+j)
			}
		}
	}
	buf.WriteByte('}')
	str := buf.String()
	return str
}

// Len returns a number of words in the set
func (s *IntSet) Len() int {
	var n int

	wcount := len(s.words)

	for i := 0; i < wcount; i++ {
		if s.words[i] == 0 {
			continue
		}
		for j := 0; j < uintSize; j++ {
			if (s.words[i] & (1 << j)) != 0 {
				n++
			}
		}
	}
	return n
}

// Remove removes the word from the set
func (s *IntSet) Remove(x int) {
	word, bit := x/uintSize, uint(x%uintSize)

	if len(s.words) < word {
		return
	}
	s.words[word] &^= (1 << bit)
}

// Clear removes all words from the set
func (s *IntSet) Clear() {
	s.words = []uint{} // or use for loop to set all words to 0
}

// Copy returns pointer to the copy of the set
func (s *IntSet) Copy() *IntSet {
	wcount := len(s.words)
	words := make([]uint, wcount, wcount)
	scopy := &IntSet{words}
	// for i := 0; i < wcount; i++ {
	// 	scopy.words[i] = s.words[i]
	// }
	copy(scopy.words, s.words)
	return scopy

}

// IntersectWith modifies the set to contain only elements that are present in s and t sets.
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, w := range s.words {
		if w == 0 {
			continue
		}

		for j := 0; j < uintSize; j++ {
			if (w & (1 << j)) != 0 {
				x := uintSize*i + j
				if t.Has(x) {
					s.Remove(x)
				}
			}
		}
	}
}

// DifferenceWith returns a set of values present in the s set and not in the t set.
func (s *IntSet) DifferenceWith(t *IntSet) *IntSet {
	u := &IntSet{}

	for i, w := range s.words {
		if w == 0 {
			continue
		}

		for j := 0; j < uintSize; j++ {
			if w&(1<<j) != 0 {
				x := uintSize*i + j
				if t.Has(x) == false {
					u.Add(x)
				}
			}

		}
	}

	return u
}

// SymmetricDifference returns a set of values that are present in t or s sets but not in both
func (s *IntSet) SymmetricDifference(t *IntSet) *IntSet {
	u := t.Copy()

	for i, w := range s.words {
		if w == 0 {
			continue
		}

		for j := 0; j < uintSize; j++ {
			x := uintSize*i + j

			if w&(1<<j) != 0 {
				if u.Has(x) == false {
					u.Add(x)
					continue
				}
				u.Remove(x)
			}

		}
	}

	return u
}

// Elems returns all elements in the set as slice of int
func (s *IntSet) Elems() []int {
	x := []int{}

	for i, w := range s.words {
		if w == 0 {
			continue
		}

		for j := 0; j < uintSize; j++ {
			if w&(1<<j) != 0 {
				x = append(x, (uintSize*i + j))
			}
		}
	}

	return x
}

func main() {
	intset := IntSet{}
	intset2 := IntSet{}
	intset2.Add(54643643)
	intset.Add(63)
	intset.Add(221)
	intset.AddAll(22, 11, 11, 511)
	intset.Add(511)
	intset.UnionWith(&intset2)

	fmt.Printf("%t\n", intset.Has(11))
	fmt.Println(intset.String(), intset.Len())
	intset.Remove(22)
	fmt.Println(intset.String(), intset.Len())
	intset2.Clear()
	fmt.Println(intset2.String(), intset2.Len())
	intset3 := intset.Copy()
	intset.Add(22)
	fmt.Printf("%t\n", intset.Has(22))

	fmt.Println(intset.String(), intset.Len())
	fmt.Println(intset3.String(), intset3.Len())
	fmt.Println()

	intset.AddAll(10, 100, 1000, 10000)
	intset3 = intset.Copy()
	intset3.Remove(100)
	intset3.Remove(1000)
	fmt.Println("inset ", intset.String(), intset.Len())
	fmt.Println("inset3", intset3.String(), intset3.Len())
	fmt.Println("difference between intset and intset3 is:")
	intsetDifference := intset.DifferenceWith(intset3)

	fmt.Println(intsetDifference)
	fmt.Println()
	fmt.Println("inset ", &intset)
	fmt.Println("inset3", intset3)
	fmt.Println()
	intset.IntersectWith(intset3)
	fmt.Println("intersected intset with intset3 is:")
	fmt.Println(&intset)

	fmt.Println()
	intset.Clear()
	intset3.Clear()
	intset.AddAll(1, 2, 3)
	intset3.AddAll(3, 4, 5)
	fmt.Println("inset ", &intset)
	fmt.Println("inset3", intset3)
	intset4 := intset.SymmetricDifference(intset3)
	fmt.Println("symmetric difference between intset and intset3 is:")
	fmt.Println(intset4)

	fmt.Println()
	fmt.Println("Elements of inset4 are:", intset4.Elems())
}
