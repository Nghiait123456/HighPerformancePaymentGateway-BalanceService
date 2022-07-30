package main

import "fmt"

var a int

func init() {
	fmt.Println(a)
	a = 10
	fmt.Println(a)
}

func main() {
	fmt.Println(a)
}
