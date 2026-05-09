package main

import "fmt"

type Vec2 struct {
	X int
	Y int
}

func (v Vec2) plus(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

func (v *Vec2) moveX(delta int) {
	v.X = v.X + delta
}

func main() {
	a := Vec2{1, 2}
	b := Vec2{3, 4}
	c := a.plus(b)
	c.moveX(10)
	fmt.Println(c.X)
	fmt.Println(c.Y)
}
