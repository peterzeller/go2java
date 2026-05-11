package main

import "fmt"

func fi() int          { return 1 }
func fi8() int8        { return 2 }
func fi16() int16      { return 3 }
func fi32() int32      { return 4 }
func fi64() int64      { return 5 }
func fu() uint         { return 6 }
func fu8() uint8       { return 255 }
func fu8one() uint8    { return 1 }
func fu16() uint16     { return 65535 }
func fu32one() uint32  { return 1 }
func fu32two() uint32  { return 2 }
func fu32high() uint32 { return 0x80000000 }
func fup() uintptr     { return 9 }
func fb() byte         { return 255 }
func fr() rune         { return 11 }

func ops8(a uint8, b uint8) {
	fmt.Println(a <= b)
	fmt.Println(a < b)
	fmt.Println(a > b)
	fmt.Println(a >= b)
}

func loopOps32() {
	i := 0
	for i < 3 {
		a := fu32one()
		if i == 1 {
			a = fu32two()
		}
		if i == 2 {
			a = fu32high()
		}

		j := 0
		for j < 3 {
			b := fu32one()
			if j == 1 {
				b = fu32two()
			}
			if j == 2 {
				b = fu32high()
			}
			fmt.Println(a < b)
			fmt.Println(a <= b)
			fmt.Println(a > b)
			fmt.Println(a >= b)
			fmt.Println(a == b)
			fmt.Println(a != b)
			if i < 2 {
				if j < 2 {
					fmt.Println(a / b)
					fmt.Println(a % b)
				}
			}
			j = j + 1
		}
		i = i + 1
	}
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
	fmt.Println(fup())
	fmt.Println(fb())
	fmt.Println(fr())
	ops8(fu8one(), fu8())
	ops8(fu8(), fu8one())
	loopOps32()
}
