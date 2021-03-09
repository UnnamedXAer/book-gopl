package memo_test

import (
	"testing"

	memo "github.com/unnamedxaer/book-gopl/ch9/memo1"
	"github.com/unnamedxaer/book-gopl/ch9/memotest"
)

func MemoGetSequentialTest(t *testing.T) {
	m := memo.New(memotest.HTTPGetBody)
	memotest.GetSequential(t, m)
}
