package palindromesort

import (
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	testCases := []struct {
		desc   string
		wanted bool
		given  []string
	}{
		{
			desc:   "correct palindromes",
			wanted: true,
			given: []string{
				"ala",
				"civic",
				"123321",
				"!2,..,2!",
				"╪",
				"",
			},
		},
		{
			desc:   "not a palindromes",
			wanted: false,
			given: []string{
				"alal",
				"╚Ä",
				"12",
				"1 1 ",
				"11 ",
				"aA",
				",.",
			},
		},
	}

	for _, tC := range testCases {
		for _, given := range tC.given {
			t.Run(tC.desc, func(t *testing.T) {
				got := IsPalindrome(seq(given))
				if tC.wanted != got {
					got := IsPalindrome(seq(given))

					t.Errorf("input %v, wanted %v, got %v", given, tC.wanted, got)
				}
			})
		}
	}
}
