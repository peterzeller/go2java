package main

import "fmt"

func ints() []int         { return []int{-1, 1, 2147483647} }
func int8s() []int8       { return []int8{-1, 1, 127} }
func int16s() []int16     { return []int16{-1, 1, 32767} }
func int32s() []int32     { return []int32{-1, 1, 2147483647} }
func int64s() []int64     { return []int64{-1, 1, 9223372036854775807} }
func uints() []uint       { return []uint{1, 2, 2147483647} }
func uint8s() []uint8     { return []uint8{1, 128, 255} }
func uint16s() []uint16   { return []uint16{1, 32768, 65535} }
func uint32s() []uint32   { return []uint32{1, 2, 3} }
func uint64s() []uint64   { return []uint64{1, 2, 3} }
func uintptrs() []uintptr { return []uintptr{1, 2, 9} }
func bytesVals() []byte   { return []byte{1, 128, 255} }
func runesVals() []rune   { return []rune{-1, 1, 1114111} }

func loopInt(vals []int) {
	i := 0
	for i < 3 {
		j := 0
		for j < 3 {
			a := vals[i]
			b := vals[j]
			fmt.Println(a < b)
			fmt.Println(a <= b)
			fmt.Println(a > b)
			fmt.Println(a >= b)
			fmt.Println(a == b)
			fmt.Println(a != b)
			if b != 0 {
				fmt.Println(a / b)
				fmt.Println(a % b)
			}
			j = j + 1
		}
		i = i + 1
	}
}

func loopUint32(vals []uint32) {
	i := 0
	for i < 3 {
		j := 0
		for j < 3 {
			a := vals[i]
			b := vals[j]
			fmt.Println(a < b)
			fmt.Println(a <= b)
			fmt.Println(a > b)
			fmt.Println(a >= b)
			fmt.Println(a == b)
			fmt.Println(a != b)
			if b != 0 {
				fmt.Println(a / b)
				fmt.Println(a % b)
			}
			j = j + 1
		}
		i = i + 1
	}
}

func loopUint8(vals []uint8) {
	i := 0
	for i < 3 {
		j := 0
		for j < 3 {
			a := vals[i]
			b := vals[j]
			fmt.Println(a < b)
			fmt.Println(a <= b)
			fmt.Println(a > b)
			fmt.Println(a >= b)
			fmt.Println(a == b)
			fmt.Println(a != b)
			if b != 0 {
				fmt.Println(a / b)
				fmt.Println(a % b)
			}
			j = j + 1
		}
		i = i + 1
	}
}

func main() {
	fmt.Println(ints()[0])
	fmt.Println(int8s()[1])
	fmt.Println(int16s()[2])
	fmt.Println(int32s()[0])
	fmt.Println(int64s()[1])
	fmt.Println(uints()[2])
	fmt.Println(uint8s()[2])
	fmt.Println(uint16s()[1])
	fmt.Println(uint32s()[1])
	fmt.Println(uint64s()[1])
	fmt.Println(uintptrs()[2])
	fmt.Println(bytesVals()[2])
	fmt.Println(runesVals()[2])

	loopInt(ints())
	loopUint8(uint8s())
	loopUint32(uint32s())
}
