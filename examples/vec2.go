package main

import "fmt"

type Vec2 struct {
	x int
	y int
}

func (v Vec2) plus(other Vec2) Vec2 {
	return Vec2{v.x + other.x, v.y + other.y}
}

func main() {
	a := Vec2{1, 2}
	b := Vec2{3, 4}
	c := a.plus(b)
	fmt.Println(c.x)
	fmt.Println(c.y)
}
