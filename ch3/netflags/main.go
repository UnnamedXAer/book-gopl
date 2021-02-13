package main

import (
	"fmt"
)

// Flag ...
type Flag uint

// package flags
const (
	FlagUp Flag = 1 << iota
	FlagBroadcast
	FlagLoopback
	FlagPointToPoint
	FlagMulticast
)

var status Flag

func IsUp(v Flag) bool {
	return v&FlagUp == FlagUp
}

func TurnUp(v *Flag) {
	*v |= FlagUp
}

func TurnDown(v *Flag) {
	*v &^= FlagUp
}

func SetFlag(v *Flag, f Flag) {
	*v |= f
}

func main() {

	fmt.Printf("% 5d -> %08[1]b\n", FlagUp)
	fmt.Printf("% 5d -> %08[1]b\n", FlagBroadcast)
	fmt.Printf("% 5d -> %08[1]b\n", FlagLoopback)
	fmt.Printf("% 5d -> %08[1]b\n", FlagPointToPoint)
	fmt.Printf("% 5d -> %08[1]b\n", FlagMulticast)
	fmt.Printf("% 5d -> %08[1]b\n", -5<<3)
	fmt.Printf("% 5d -> %08[1]b\n", 5<<3)
	fmt.Printf("% 5d -> %08[1]b\n", 5>>3)
	fmt.Printf("% 5d -> %08[1]b\n", 40>>3)
	fmt.Printf("% 5d -> %08[1]b\n", 40/2/2/2)
	fmt.Printf("% 5d -> %08[1]b\n", 5)
	fmt.Printf("% 5d -> %08[1]b\n", 5<<1)
	fmt.Printf("% 5d -> %08[1]b\n", 5<<2)
	fmt.Printf("% 5d -> %08[1]b\n", 5<<3<<57)
	fmt.Printf("% 5d -> %08[1]b\n", FlagUp)
	fmt.Printf("% 5d -> %08[1]b, %t\n", status, IsUp(status))
	TurnUp(&status)
	fmt.Printf("% 5d -> %08[1]b, %t\n", status, IsUp(status))
	TurnDown(&status)
	fmt.Printf("% 5d -> %08[1]b, %t\n", status, IsUp(status))
	SetFlag(&status, FlagBroadcast)
	SetFlag(&status, FlagPointToPoint)
	SetFlag(&status, FlagUp)
	fmt.Printf("% 5d -> %08[1]b, %t\n", status, IsUp(status))

	for i := 0; i < 10; i++ {
		fmt.Printf("%d. % 20v\n", i, 1<<(10*i))
	}

}
