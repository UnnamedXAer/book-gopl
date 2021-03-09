package memo_test

import (
	"testing"

	memo "github.com/unnamedxaer/book-gopl/ch9/memo4"
	"github.com/unnamedxaer/book-gopl/ch9/memotest"
)

var m = memo.New(memotest.HTTPGetBody)

func TestMemoGetSequential(t *testing.T) {
	var m = memo.New(memotest.HTTPGetBody)
	memotest.GetSequential(t, m)
}

func TestMemoConcurrent(t *testing.T) {
	var m = memo.New(memotest.HTTPGetBody)
	memotest.GetConcurrent(t, m)
}
