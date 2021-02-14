package main

import (
	"fmt"
	"time"
)

// Employee represents a person which is employed in our company
type Employee struct {
	ID        int
	Name      string
	Address   string
	DoB       time.Time
	Position  string
	Salary    int
	ManagerID int
}

func main() {
	// employeeManipulation()
}

func employeeManipulation() {
	emp := Employee{
		ID:        99,
		Name:      "Jonhy",
		Address:   "15 North Avenue, Clearwater FL, 53533, USA",
		DoB:       time.Date(1750, 1, 13, 0, 0, 0, 0, time.UTC),
		Position:  "manager",
		Salary:    12500,
		ManagerID: 1,
	}
	fmt.Printf("%T\n", emp)

	var empPointer *Employee = &Employee{
		ID:        99,
		Name:      "Jonhy",
		Address:   "15 North Avenue, Clearwater FL, 53533, USA",
		DoB:       time.Date(1750, 1, 13, 0, 0, 0, 0, time.UTC),
		Position:  "manager",
		Salary:    12500,
		ManagerID: 1,
	}
	fmt.Printf("%T\n", empPointer)

	(empPointer).Position += " pro team player!"

	// position := &empPointer.Position
	// *position = "Senior " + *position
	fmt.Println(empPointer.Position)
	fmt.Println(empPointer.Salary)
	empPointer.Salary -= 2000
	fmt.Println(empPointer.Salary)
}
