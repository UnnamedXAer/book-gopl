// Package cake1 prepares a cakes
// it demonstraits a serially confined approach as a way of avoiding data race
package cake1

import "math/rand"

type Cake struct {
	state string
	id    int
}

// Baker bakes cakes
func Baker(cooked chan<- *Cake) {
	for {
		cake := new(Cake)
		cake.state = "cooked"
		cake.id = rand.Int()
		cooked <- cake // baker never touches this cake again
	}
}

func Icer(iced chan<- *Cake, cooked <-chan *Cake) {
	for cake := range cooked {
		cake.state = "iced"
		iced <- cake // icer never touches this cake again
	}
}

func Deliver(delivered chan<- *Cake, iced <-chan *Cake) {
	for cake := range iced {
		cake.state = "delivered"
		delivered <- cake
	}
}
