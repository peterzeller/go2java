package main

import "fmt"

func takeMid(in []int) []int {
	return in[1:3]
}

func main() {
	a := []int{1, 2, 3, 4}
	a = append(a, 5)
	b := a[1:4]
	c := takeMid(a)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}
