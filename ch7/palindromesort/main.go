package palindromesort

import "sort"

type seq []rune

func (r seq) Len() int           { return len(r) }
func (r seq) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r seq) Less(i, j int) bool { return r[i] < r[j] }

// func main() {

// }

// IsPalindrome reports whether the sequence s is a palindrome
func IsPalindrome(s sort.Interface) bool {
	count := s.Len()
	for i := 0; i < count/2; i++ {
		if (s.Less(i, count-1-i) == false && s.Less(count-1-i, i) == false) == false {
			return false
		}
	}

	return true
}
