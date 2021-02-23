package main

import (
	"bytes"
	"fmt"
)

// IntSet is a set of small non-negative integers.
// Its zero value represents the empty state.
type IntSet struct {
	words []uint64
}

// Has reports the set contains the non-negative value x
func (s *IntSet) Has(x int) bool {
	word, bit := x/64, uint(x%64)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set
func (s *IntSet) Add(x int) {
	word, bit := x/64, uint(x%64)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
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
		for j := 0; j < 64; j++ {
			if word&(1<<(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
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

	//doesnt work, it not counts correctly words like 12,52
	for _, v := range s.words {
		if v != 0 {
			n++
		}
	}
	return n
}

// Remove removes the word from the set
func (s *IntSet) Remove(x int) {
	word, bit := x/64, uint(x%64)

	if len(s.words) < word {
		return
	}
	s.words[word] &^= (1 << bit)
}

// Clear removes all words from the set
func (s *IntSet) Clear() {
	s.words = []uint64{} // or use for loop to set all words to 0
}

func (s *IntSet) Copy() *IntSet {
	wcount := len(s.words)
	words := make([]uint64, wcount, wcount)
	scopy := &IntSet{words}
	// for i := 0; i < wcount; i++ {
	// 	scopy.words[i] = s.words[i]
	// }
	copy(scopy.words, s.words)
	return scopy

}

func main() {
	intset := IntSet{}
	intset2 := IntSet{}
	intset2.Add(54643643)
	intset.Add(63)
	intset.Add(221)
	intset.Add(511)
	intset.UnionWith(&intset2)

	// fmt.Printf("%t\n", intset.Has(11))
	// fmt.Println(intset.String(), intset.Len())
	intset.Remove(221)
	// fmt.Println(intset.String(), intset.Len())
	intset2.Clear()
	// fmt.Println(intset2.String(), intset2.Len())
	intset3 := intset.Copy()
	intset.Add(22)
	fmt.Printf("%t\n", intset.Has(22))

	fmt.Println(intset.String(), intset.Len())
	fmt.Println(intset3.String(), intset3.Len())
}
