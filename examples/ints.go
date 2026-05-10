package main

import "fmt"

func fi() int { return 1 }
func fi8() int8 { return 2 }
func fi16() int16 { return 3 }
func fi32() int32 { return 4 }
func fi64() int64 { return 5 }
func fu() uint { return 6 }
func fu8() uint8 { return 100 }
func fu16() uint16 { return 1000 }
func fu32() uint32 { return 7 }
func fu64() uint64 { return 8 }
func fup() uintptr { return 9 }
func fb() byte { return 10 }
func fr() rune { return 11 }

func ops32(a uint32, b uint32) {
	fmt.Println(a / b)
	fmt.Println(a % b)
	fmt.Println(a < b)
	fmt.Println(a <= b)
	fmt.Println(a > b)
	fmt.Println(a >= b)
	fmt.Println(a == b)
	fmt.Println(a != b)
	fmt.Println(a + b)
	fmt.Println(a - b)
	fmt.Println(a * b)
}

func ops64(x uint64, y uint64) {
	fmt.Println(x / y)
	fmt.Println(x % y)
	fmt.Println(x < y)
	fmt.Println(x <= y)
	fmt.Println(x > y)
	fmt.Println(x >= y)
	fmt.Println(x == y)
	fmt.Println(x != y)
	fmt.Println(x + y)
	fmt.Println(x - y)
	fmt.Println(x * y)
}

func main() {
	fmt.Println(fi())
	fmt.Println(fi8())
	fmt.Println(fi16())
	fmt.Println(fi32())
	fmt.Println(fi64())
	fmt.Println(fu())
	fmt.Println(fu8())
	fmt.Println(fu16())
	fmt.Println(fu32())
	fmt.Println(fu64())
	fmt.Println(fup())
	fmt.Println(fb())
	fmt.Println(fr())
	ops32(7, 2)
	ops64(7, 3)
}
