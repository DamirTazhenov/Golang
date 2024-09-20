package main

import "fmt"

func main() {
	hello_world()

	variables()

	condition()

	cycle()

	switch_case()

	fmt.Println(add(1, 5))

	first, second := swap("hello", "world")
	fmt.Println(first, second)

	quotient, remainder := divMod(10, 3)
	fmt.Println("Quotient:", quotient, "Remainder:", remainder)

	person := Person{Name: "John", Age: 30}
	person.Greet()

	manager := Manager{
		Employee: Employee{
			Name: "John Doe",
			ID:   "1234",
		},
		Department: "Engineering",
	}

	manager.Work()

	circle := Circle{Radius: 5}
	rectangle := Rectangle{Width: 10, Height: 5}

	PrintArea(circle)
	PrintArea(rectangle)

	testJson()
}
