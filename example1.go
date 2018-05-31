//example1.go
package main

import (
	"fmt"
	"math"
	"runtime"
	"strconv"
)

type Triangle struct {
	a float64
	b float64
	c float64
}

func main() {
	fmt.Printf("Hello, world !\n")
	dataTypes()
	var total = sumOfInts(5, 2, 1)
	fmt.Printf("Sum is : %d", total)
	//get 8th element of fibonacci sequence using recursion 1 1 2 3 5 8 13 21
	fmt.Printf("8th member of Fibonacci sequence is: %d", fibRecursive(8))
	fmt.Printf("Factorial of 5 is : %d\n", fact(5))
	fmt.Printf("Number 11 is prime ? %t", isPrime(3))
	fmt.Printf("There are %d  primes which are less than 100\n", numberOfPrimes(100))
	pointersExample()
	fmt.Printf("Area of triangle is : %f \n", calculateTriangleArea(Triangle{3, 4, 5}))
	arrayExample()
	mapExample()
	errorHandlingExample()
	switchExample()
	concurencyExample()
}

func dataTypes() {
	var x bool = true
	var counter int
	var text string
	version := 5.0
	//can be initialized and declared implicitly like this

	//imports which are not user, and variables that are defined but not user are not allowed

	//var notUsed int
	i := 42
	f := float64(i)

	counter = 3
	text = "Some text"
	fmt.Printf("X is : %t , version is %f, counter is : %d, text is : %s\n", x, version, counter, text)

	fmt.Println(i)
	fmt.Println(f)
}

func sumOfInts(x, y, z int) int {
	var total int
	total = x + y + z
	return total
}

func fibRecursive(n int) int {
	//} in the same line as else
	if n < 2 {
		return n
	} else {
		return fibRecursive(n-1) + fibRecursive(n-2)
	}
}

func fact(n int) int {
	//} in the same line as else
	if n < 1 {
		return 1
	} else {
		return n * fact(n-1)
	}

}

func isPrime(n int) bool {
	//6k + 1 || 6k -1 except of 2 and 3
	if n < 2 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func numberOfPrimes(n int) int {
	i := 0
	total := 0
	//forever loop
	for {
		if i == n {
			break
		}
		if isPrime(i) {
			total++
		}
		i++
	}
	return total
}

func pointersExample() {
	i, j := 42, 32.4
	p := &i         // point to i
	fmt.Println(*p) // read i through the pointer
	*p = 21         // set value of the i
	fmt.Println(*p) // write value of the variable on which the p is pointing
	fmt.Println(i)  // write value of the variable i - > same as above
	fmt.Println(j)  // is this line necessary ?
}

func calculateTriangleArea(t Triangle) float64 {
	//Heron's formula
	//import math
	s := (t.a + t.b + t.c) / 2
	area := math.Sqrt(s * (s - t.a) * (s - t.b) * (s - t.c))
	return area

}

func arrayExample() {
	var names [2]string
	names[0] = "Marc"
	names[1] = "John"
	fmt.Println(names[0], names[1])
	fmt.Println(names)

	primes := [8]int{2, 3, 5, 7, 11, 13, 17, 19}
	fmt.Println(primes)
	primesLessThan10 := primes[0:4]
	fmt.Println(primesLessThan10)
	//slice is dynamically-sized
	primesLessThan10 = primes[0:8]
	fmt.Println(primesLessThan10)
	primesLessThan10[0] = 13
	//primes is also changed, because slice is acutally reference to the array
	fmt.Println(primes)
	fmt.Printf("Number of elements in primes array is %d\n", len(primes))
	newSlice := make([]int, 5) // dinamically create slice
	fmt.Println(newSlice)
	newSlice = append(newSlice, 3)
	fmt.Println(newSlice)

}

func mapExample() {

	var m map[string]int
	m = make(map[string]int)
	m["S1"] = 11
	m["S2"] = 12
	m["M3"] = 15

	for key, value := range m {
		fmt.Printf("Key is : %s, Value is : %d\n", key, value)
	}
}

func errorHandlingExample() {

	/*
		Trying to convert non integer type to integer should invoke error which
		should be handled properly, with appropriate message returned to the user
	*/
	i, err := strconv.Atoi("44s")
	if err != nil {
		fmt.Printf("Not a valid number\n")
		//fmt.Printf("Not a valid number", err)
		return
	}
	fmt.Println("Converted integer:", i)
}

func switchExample() {
	fmt.Print("Go runs on ")
	//os:=runtime.GOOS; is done just before switch
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.\n")
	case "linux":
		fmt.Println("Linux.\n")
	default:
		fmt.Printf("%s.\n", os)
	}
}

func sumArrayIntoChannel(s []int, c chan int, ordNum int) {
	fmt.Printf("Started routine %d\n", ordNum)
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
	fmt.Printf("Finished routine %d\n", ordNum)
}
func concurencyExample() {

	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	c := make(chan int) // need to create channel
	go sumArrayIntoChannel(s[:len(s)/4], c, 1)
	go sumArrayIntoChannel(s[5:10], c, 2)
	go sumArrayIntoChannel(s[10:15], c, 3)
	go sumArrayIntoChannel(s[15:20], c, 4)
	x, y, z, v := <-c, <-c, <-c, <-c // receive from c

	fmt.Println(x, y, z, v, x+y+z+v)

}
