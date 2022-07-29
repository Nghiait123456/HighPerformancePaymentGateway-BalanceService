package main

import (
	"fmt"
	"github.com/gammazero/deque"
)

func main() {
	var q deque.Deque[string]
	q.PushBack("foo")
	q.PushBack("bar")
	q.PushBack("baz")

	fmt.Println(q.Len())   // Prints: 3
	fmt.Println(q.Front()) // Prints: foo
	fmt.Println(q.Back())  // Prints: baz

	q.PopFront() // remove "foo"
	q.PopBack()  // remove "baz"

	q.PushFront("hello")
	q.PushBack("world")

	// Consume deque and print elements.
	for q.Len() != 0 {
		fmt.Println(q.PopFront())
	}
}
