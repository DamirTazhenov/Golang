package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Name string
	ID   string
}

type Manager struct {
	Employee
	Department string
}

func (p Person) Greet() {
	fmt.Printf("Hello, my name is %s and I am %d years old.\n", p.Name, p.Age)
}

func (e Employee) Work() {
	fmt.Printf("%s with ID %s is working.\n", e.Name, e.ID)
}
