package main

import "fmt"

var myInt int = 10
var myFloat float64 = 3.14
var myString string = "Hello Go"
var myBool bool = true

var anotherInt = 20
var anotherString = "Another example"

func hello_world() {
	fmt.Println("Hello, World!")
}

func variables() {
	anotherFloat := 5.67
	anotherBool := false

	fmt.Println(anotherFloat, anotherBool)

	var myInt int = 10
	myFloat := 3.14
	myString := "Hello Go"
	var myBool bool = true

	fmt.Printf("myInt: %v, Type: %T\n", myInt, myInt)
	fmt.Printf("myFloat: %v, Type: %T\n", myFloat, myFloat)
	fmt.Printf("myString: %v, Type: %T\n", myString, myString)
	fmt.Printf("myBool: %v, Type: %T\n", myBool, myBool)
}

func condition() {
	var number int
	fmt.Print("Enter an integer: ")
	fmt.Scanln(&number)

	if number > 0 {
		fmt.Println("The number is positive.")
	} else if number < 0 {
		fmt.Println("The number is negative.")
	} else {
		fmt.Println("The number is zero.")
	}
}

func cycle() {
	sum := 0
	for i := 1; i <= 10; i++ {
		sum += i
	}
	fmt.Println("The sum of the first 10 natural numbers is:", sum)
}

func switch_case() {
	var day int
	fmt.Print("Enter a number (1-7): ")
	fmt.Scanln(&day)

	switch day {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3:
		fmt.Println("Wednesday")
	case 4:
		fmt.Println("Thursday")
	case 5:
		fmt.Println("Friday")
	case 6:
		fmt.Println("Saturday")
	case 7:
		fmt.Println("Sunday")
	default:
		fmt.Println("Invalid day number")
	}
}

func add(x int, y int) int {
	return x + y
}

func swap(a, b string) (string, string) {
	return b, a
}

func divMod(x, y int) (int, int) {
	return x / y, x % y
}
